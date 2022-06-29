package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws/credentials"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/evergreen-ci/juniper/gopb"
	"github.com/evergreen-ci/poplar"
	"github.com/evergreen-ci/poplar/rpc/internal"
	"github.com/google/uuid"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/grip/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

// UploadReportOptions contains all the required options for uploading a report
// to both Cedar and the Data Pipes service.
type UploadReportOptions struct {
	Report              *poplar.Report
	ClientConn          *grpc.ClientConn
	SerializeUpload     bool
	AWSAccessKey        string
	AWSSecretKey        string
	AWSToken            string
	DataPipesHost       string
	DataPipesRegion     string
	DataPipesHTTPClient *http.Client
	DryRun              bool
}

// SignedURL is a struct representing a signed url returned from the data pipes API.
type SignedURL struct {
	URL            string `json:"signed_url"`
	ExpirationSecs int    `json:"expiration_secs"`
}

// UploadReport does the following:
// 1. Ingest the report.json.
// 2. Send all artifact files to the user-specified buckets using user-specified keys.
// 3. Send the metrics metadata and pre-calculated summaries to Cedar over gRPC using the Cedar creds.
// 4. Send the metrics metadata and pre-calculated summaries to the Data Pipes service using AWS keys.
func UploadReport(ctx context.Context, opts UploadReportOptions) error {

	// Errors uploading to Data Pipes will be only be logged while Cedar is in use.
	defer func() {
		if err := errors.Wrap(opts.validate(), "validating options"); err != nil {
			grip.Warning(message.Fields{
				"op":    "validate options",
				"error": err,
			})
			return
		}

		err := errors.Wrap(uploadResultsToDataPipes(&opts), "uploading results to DataPipes")
		if err != nil {
			grip.Warning(message.Fields{
				"op":    "uploadResultsToDataPipes",
				"error": err,
			})
		}
	}()

	if err := opts.convertAndUploadArtifacts(ctx); err != nil {
		return errors.Wrap(err, "uploading tests for report")
	}
	return errors.Wrap(uploadTests(ctx, gopb.NewCedarPerformanceMetricsClient(opts.ClientConn), opts.Report, opts.Report.Tests, opts.DryRun), "uploading tests for report")
}

func (opts *UploadReportOptions) convertAndUploadArtifacts(ctx context.Context) error {
	catcher := grip.NewBasicCatcher()
	testChan := make(chan poplar.Test, len(opts.Report.Tests)*2)
	go opts.artifactProducer(ctx, testChan, catcher)

	var wg sync.WaitGroup
	workers := 1
	if !opts.SerializeUpload {
		workers = runtime.NumCPU()
	}
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go opts.artifactConsumer(ctx, testChan, catcher, &wg)
	}

	wg.Wait()
	return catcher.Resolve()
}

func (opts *UploadReportOptions) artifactProducer(ctx context.Context, testChan chan poplar.Test, catcher grip.Catcher) {
	defer func() {
		catcher.Add(recovery.HandlePanicWithError(recover(), nil, "artifact upload producer"))
		close(testChan)
	}()

	testQueue := make([]poplar.Test, len(opts.Report.Tests))
	for i, test := range opts.Report.Tests {
		testQueue[i] = test
	}
	for len(testQueue) != 0 {
		test := testQueue[0]
		testQueue = testQueue[1:]
		select {
		case testChan <- test:
		case <-ctx.Done():
			return
		}

		for _, subTest := range test.SubTests {
			testQueue = append(testQueue, subTest)
		}
	}
}

func (opts *UploadReportOptions) artifactConsumer(ctx context.Context, testChan chan poplar.Test, catcher grip.Catcher, wg *sync.WaitGroup) {
	defer func() {
		catcher.Add(recovery.HandlePanicWithError(recover(), nil, "artifact upload consumer"))
		wg.Done()
	}()

	for test := range testChan {
		for j := range test.Artifacts {
			if err := ctx.Err(); err != nil {
				catcher.Add(err)
				return
			}

			if err := test.Artifacts[j].Convert(ctx); err != nil {
				catcher.Wrap(err, "converting artifact")
				continue
			}

			if err := test.Artifacts[j].SetBucketInfo(opts.Report.BucketConf); err != nil {
				catcher.Wrap(err, "setting bucket info")
				continue
			}

			grip.Info(message.Fields{
				"op":     "uploading artifact",
				"path":   test.Artifacts[j].Path,
				"bucket": test.Artifacts[j].Bucket,
				"prefix": test.Artifacts[j].Prefix,
				"file":   test.Artifacts[j].LocalFile,
			})
			catcher.Wrapf(test.Artifacts[j].Upload(ctx, opts.Report.BucketConf, opts.DryRun), "uploading artifact")
		}
	}
}

func (opts *UploadReportOptions) validate() error {
	catcher := grip.NewBasicCatcher()

	catcher.NewWhen(opts.AWSAccessKey == "", "must provide an AWS access key")
	catcher.NewWhen(opts.AWSSecretKey == "", "must provide an AWS secret key")
	catcher.NewWhen(opts.DataPipesHost == "", "must provide the Data Pipes service hostname")
	catcher.NewWhen(opts.DataPipesRegion == "", "must provide the Data Pipes service region")

	return errors.Wrap(catcher.Resolve(), "invalid upload report options")
}

// getSignedURL calls the Data Pipes API to retrieve a signed URL, where it can PUT the report JSON.
// Data Pipes docs: https://github.com/10gen/data-pipes.
func getSignedURL(opts *UploadReportOptions) (string, error) {
	service := "execute-api"
	resultType := "cedar-report"
	if opts.AWSAccessKey == "" || opts.AWSSecretKey == "" || opts.DataPipesHost == "" || opts.DataPipesRegion == "" {
		return "", errors.New("Getting signed URL failed. AWS access key, AWS secret key, data pipes host and data pipes region required.")
	}

	// See the Data Pipes documentation for more information.
	name := uuid.New()
	url := fmt.Sprintf("%s/v1/results/evergreen/%s/%s/%s/%s", opts.DataPipesHost, opts.Report.TaskID, strconv.Itoa(opts.Report.Execution), resultType, name.String())
	grip.Debug(message.Fields{
		"request_url": url,
	})
	//bytes.NewBuffer(data)
	req, err := http.NewRequest(http.MethodPut, url, nil)
	if err != nil {
		return "", err
	}

	AWSCredentials := credentials.NewStaticCredentials(opts.AWSAccessKey, opts.AWSSecretKey, opts.AWSToken)
	signer := v4.NewSigner(AWSCredentials)
	_, err = signer.Sign(req, nil, service, opts.DataPipesRegion, time.Now())
	if err != nil {
		return "", err
	}

	response, err := opts.DataPipesHTTPClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "failed getting signed url")
	}

	defer response.Body.Close()
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return "", err
	}

	var responseBody SignedURL
	err = json.Unmarshal(body, &responseBody)
	if error != nil {
		return "", errors.Wrap(err, "failed parsing signed url response")
	}

	grip.Debug(message.Fields{
		"signed_url": responseBody.URL,
	})
	return responseBody.URL, nil
}

func uploadTestReport(signedURL string, data []byte, client *http.Client) error {
	req, err := http.NewRequest(http.MethodPut, signedURL, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	response, err := client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed report upload to given signed url")
	}

	grip.Debug(message.Fields{
		"message":  "upload to Data Pipes response",
		"function": "uploadTestReport",
		"response": response,
	})

	return nil
}

// uploadResultsToDataPipes uploads the report JSON to Data Pipes for further processing.
func uploadResultsToDataPipes(opts *UploadReportOptions) error {
	// The bucket configuration contains the user's sensitive information that the data pipes
	// and processors don't need, so we redact it.
	opts.Report.BucketConf = poplar.BucketConfiguration{}
	if opts.DryRun {
		grip.Info(message.Fields{
			"message":     "dry-run mode",
			"function":    "uploadResultsToDataPipes",
			"result_data": opts.Report,
		})
		return nil
	}
	jsonResp, err := json.Marshal(opts.Report)
	if err != nil {
		return err
	}
	signedURL, err := getSignedURL(opts)
	if err != nil {
		return err
	}
	err = uploadTestReport(signedURL, jsonResp, opts.DataPipesHTTPClient)
	if err != nil {
		return err
	}
	return nil
}

func uploadTests(ctx context.Context, client gopb.CedarPerformanceMetricsClient, report *poplar.Report, tests []poplar.Test, dryRun bool) error {
	for idx, test := range tests {
		grip.Info(message.Fields{
			"num":     idx,
			"total":   len(tests),
			"parent":  test.Info.Parent != "",
			"name":    test.Info.TestName,
			"task":    report.TaskID,
			"dry_run": dryRun,
		})

		artifacts, err := extractArtifacts(ctx, report, test)
		if err != nil {
			return err
		}
		metrics, err := extractMetrics(ctx, test)
		if err != nil {
			return err
		}
		resultData := &gopb.ResultData{
			Id: &gopb.ResultID{
				Project:   report.Project,
				Version:   report.Version,
				Order:     int32(report.Order),
				Variant:   report.Variant,
				TaskName:  report.TaskName,
				TaskId:    report.TaskID,
				Mainline:  report.Mainline,
				Execution: int32(report.Execution),
				TestName:  test.Info.TestName,
				Trial:     int32(test.Info.Trial),
				Tags:      test.Info.Tags,
				Arguments: test.Info.Arguments,
				Parent:    test.Info.Parent,
				CreatedAt: internal.ExportTimestamp(test.CreatedAt),
			},
			Artifacts: artifacts,
			Rollups:   metrics,
		}

		if dryRun {
			grip.Info(message.Fields{
				"message":     "dry-run mode",
				"function":    "CreateMetricSeries",
				"result_data": resultData,
			})
		} else {
			var resp *gopb.MetricsResponse
			resp, err = client.CreateMetricSeries(ctx, resultData)
			if err != nil {
				return errors.Wrapf(err, "submitting test %d of %d", idx+1, len(tests))
			} else if !resp.Success {
				return errors.New("operation return failed state")
			}

			test.ID = resp.Id
			for i := range test.SubTests {
				test.SubTests[i].Info.Parent = test.ID
			}
		}

		if err = uploadTests(ctx, client, report, test.SubTests, dryRun); err != nil {
			return errors.Wrapf(err, "submitting subtests of '%s'", test.ID)
		}

		end := &gopb.MetricsSeriesEnd{
			Id:          test.ID,
			IsComplete:  true,
			CompletedAt: internal.ExportTimestamp(test.CompletedAt),
		}

		if dryRun {
			grip.Info(message.Fields{
				"message":           "dry-run mode",
				"function":          "CloseMetrics",
				"metric_series_end": end,
			})
		} else {
			var resp *gopb.MetricsResponse
			resp, err = client.CloseMetrics(ctx, end)
			if err != nil {
				return errors.Wrapf(err, "closing metrics series for '%s'", test.ID)
			} else if !resp.Success {
				return errors.New("operation return failed state")
			}
		}
	}

	return nil
}

func extractArtifacts(ctx context.Context, report *poplar.Report, test poplar.Test) ([]*gopb.ArtifactInfo, error) {
	artifacts := make([]*gopb.ArtifactInfo, 0, len(test.Artifacts))
	for _, a := range test.Artifacts {
		if err := a.Validate(); err != nil {
			return nil, errors.Wrap(err, "invalid artifacts")
		}
		artifacts = append(artifacts, internal.ExportArtifactInfo(&a))
		artifacts[len(artifacts)-1].Location = gopb.StorageLocation_CEDAR_S3
	}

	return artifacts, nil
}

func extractMetrics(ctx context.Context, test poplar.Test) ([]*gopb.RollupValue, error) {
	rollups := make([]*gopb.RollupValue, 0, len(test.Metrics))
	names := map[string]bool{}
	for _, r := range test.Metrics {
		if ok := names[r.Name]; ok {
			return nil, errors.Errorf("duplicate metric name '%s'", r.Name)
		}
		names[r.Name] = true

		rollups = append(rollups, internal.ExportRollup(&r))
	}

	return rollups, nil
}
