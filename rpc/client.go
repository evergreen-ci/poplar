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
// to both Cedar and the S3 data-pipes service.
type UploadReportOptions struct {
	Report             *poplar.Report
	ClientConn         *grpc.ClientConn
	SerializeUpload    bool
	AWSAccessKey       string
	AWSSecretKey       string
	AWSToken           string
	ResultsHandlerHost string
	DryRun             bool
}

// SignedUrl is a struct representing a signed url returned from the data pipes API.
type SignedUrl struct {
	SignedUrl      string `json:"signed_url"`
	ExpirationSecs int    `json:"expiration_secs"`
}

// UploadReport does the following:
// 1. Ingest the report.json.
// 2. Send all artifact files to the user-specified buckets using user-specified keys.
// 3. Send the metrics metadata and pre-calculated summaries to Cedar over gRPC using the Cedar creds.
// 4. Send the metrics metadata and pre-calculated summaries to the data-pipes service using AWS keys.
func UploadReport(ctx context.Context, opts UploadReportOptions) error {
	var returnError error
	if err := opts.convertAndUploadArtifacts(ctx); err != nil {
		returnError = errors.Wrap(err, "uploading tests for report")
	} else {
		returnError = errors.Wrap(uploadTests(ctx, gopb.NewCedarPerformanceMetricsClient(opts.ClientConn), opts.Report, opts.Report.Tests, opts.DryRun), "uploading tests for report")
	}
	err := errors.Wrap(uploadResultsToDataPipes(opts.Report, opts.AWSAccessKey, opts.AWSSecretKey, opts.AWSToken, opts.ResultsHandlerHost, opts.DryRun), "uploading results to DataPipes")
	// Errors uploading to Data Pipes will be only be logged while Cedar is in use.
	if err != nil {
		grip.Warning(message.Fields{
			"op":    "uploadResultsToDataPipes",
			"error": err,
		})
	}

	return returnError
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

// getSignedUrl calls the Data Pipes API to retrieve a signed URL, where it can PUT the report json.
// Data Pipes docs: https://github.com/10gen/data-pipes
func getSignedURL(task string, execution int, AWSRegion string, data []byte, AWSAccessKey string, AWSSecretKey string, AWSToken string, resultsHandlerHost string) (string, error) {
	service := "execute-api"
	resultType := "cedar-report"
	if AWSAccessKey == "" || AWSSecretKey == "" || resultsHandlerHost == "" || resultType == "" {
		return "", errors.New("Getting signed URL failed. AWS access key, AWS secret key, results handler host and result type required.")
	}
	client := &http.Client{}
	name := uuid.New()
	url := fmt.Sprintf("%s/evergreen/%s/%s/%s/%s", resultsHandlerHost, task, strconv.Itoa(execution), resultType, name.String())
	grip.Debug(message.Fields{
		"request_url": url,
	})
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
		return "", err
	}
	AWSCredentials := credentials.NewStaticCredentials(AWSAccessKey, AWSSecretKey, AWSToken)
	signer := v4.NewSigner(AWSCredentials)
	_, err = signer.Sign(req, nil, service, AWSRegion, time.Now())
	if err != nil {
		return "", err
	}
	response, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, error := ioutil.ReadAll(response.Body)
	if error != nil {
		return "", err
	}
	// close response body
	response.Body.Close()
	var responseBody SignedUrl
	err = json.Unmarshal(body, &responseBody)
	if error != nil {
		return "", err
	}
	grip.Info(message.Fields{
		"signed_url": responseBody.SignedUrl,
	})
	return responseBody.SignedUrl, nil
}

func uploadTestReport(signedUrl string, data []byte) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, signedUrl, bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	response, err := client.Do(req)
	if err != nil {
		return err
	}
	grip.Info(message.Fields{
		"message":  "Upload to Data Pipes response",
		"function": "uploadTestReport",
		"response": response,
	})
	return nil
}

// uploadResultsToDataPipes uploads the report json to Data Pipes for further processing.
func uploadResultsToDataPipes(report *poplar.Report, AWSAccessKey string, AWSSecretKey string, AWSToken string, resultsHandlerHost string, dryRun bool) error {
	region := report.BucketConf.Region
	report.BucketConf = poplar.BucketConfiguration{}
	if dryRun {
		grip.Info(message.Fields{
			"message":     "dry-run mode",
			"function":    "uploadResultsToDataPipes",
			"result_data": report,
		})
	} else {
		jsonResp, err := json.Marshal(report)
		if err != nil {
			return err
		}
		signedUrl, err := getSignedURL(report.TaskName, report.Execution, region, jsonResp, AWSAccessKey, AWSSecretKey, AWSToken, resultsHandlerHost)
		if err != nil {
			return err
		}
		err = uploadTestReport(signedUrl, jsonResp)
		if err != nil {
			return err
		}
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
