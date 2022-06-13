package rpc

import (
	"context"
	"runtime"
	"sync"
	"time"
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strconv"
	"github.com/evergreen-ci/juniper/gopb"
	"github.com/evergreen-ci/poplar"
	"github.com/evergreen-ci/poplar/rpc/internal"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/grip/recovery"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"github.com/aws/aws-sdk-go/aws/signer/v4"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/google/uuid"
)

const (
	ResultsHandlerHost = "https://yv5aascdn3.execute-api.us-east-1.amazonaws.com/prod/v1/results"
)

type UploadReportOptions struct {
	Report          *poplar.Report
	ClientConn      *grpc.ClientConn
	SerializeUpload bool
	AWSSecretKey    string
	AWSAccessKey    string
	DryRun          bool
}

type SignedUrl struct {
	Signed_Url string
	Expiration_secs int
}

func UploadReport(ctx context.Context, opts UploadReportOptions) error {
	if err := opts.convertAndUploadArtifacts(ctx); err != nil {
		return errors.Wrap(err, "uploading tests for report")
	}
	return errors.Wrap(uploadTests(ctx, gopb.NewCedarPerformanceMetricsClient(opts.ClientConn), opts.Report, opts.Report.Tests, opts.AWSSecretKey, opts.AWSAccessKey, opts.DryRun),
		"uploading tests for report")
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

			// if err := test.Artifacts[j].Convert(ctx); err != nil {
			// 	catcher.Wrap(err, "converting artifact")
			// 	continue
			// }

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
			// catcher.Wrapf(test.Artifacts[j].Upload(ctx, opts.Report.BucketConf, opts.DryRun), "uploading artifact")
		}
	}
}


func getSignedURL(task string, execution int, awsRegion string, data []byte, awsSecretKey string, awsAccessKey string) (string, error) {
	resultType := "cedar-report"
	service := "execute-api"
	if awsSecretKey == "" || awsAccessKey == "" {
		return "", errors.New("Getting signed URL failed. AWS secret key or/and access key not specified")
	}
	client := &http.Client{}
	name := uuid.New()
	url := fmt.Sprintf("%s/evergreen/%s/%s/%s/%s", ResultsHandlerHost, task, strconv.Itoa(execution), resultType, name.String())
	req, err := http.NewRequest(http.MethodPut, url, bytes.NewBuffer(data))
	if err != nil {
			return "", err
	}
	aws_credentials := credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, "")
	signer := v4.NewSigner(aws_credentials)
	_, err = signer.Sign(req, nil, service, awsRegion, time.Now())
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
	var response_body SignedUrl
	err = json.Unmarshal(body, &response_body)
	if error != nil {
		 return "", err
	}
	grip.Info(message.Fields{
		"request_url": url,
		"signed_url": response_body.Signed_Url,
	})
	return response_body.Signed_Url, nil
}

func uploadTestReport(signedUrl string, data []byte) error {
	client := &http.Client{}
	req, err := http.NewRequest(http.MethodPut, signedUrl, bytes.NewBuffer(data))
	if err != nil {
			return err
	}
	_, err = client.Do(req)
	if err != nil {
			return err
	}
	return nil
}

// uploadResults uploads cedar results to S3 bucket which triggers rollups in perf-summarizer-service(PSS)
func uploadResultsToPSS(report *poplar.Report, test poplar.Test, awsSecretKey string, awsAccessKey string, dryRun bool,) error {
	if !dryRun {
		grip.Info(message.Fields{
			"message":     "dry-run mode",
			"function":    "uploadResultsToPSS",
			"result_data": test,
		})
	} else {
		jsonResp, err := json.Marshal(test)
		if err != nil {
				return err
		}
		signedUrl, err := getSignedURL(report.TaskName, report.Execution, report.BucketConf.Region, jsonResp, awsSecretKey, awsAccessKey)
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

func uploadTests(ctx context.Context, client gopb.CedarPerformanceMetricsClient, report *poplar.Report, tests []poplar.Test, awsSecretKey string, awsAccessKey string, dryRun bool) error {
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

		err = uploadResultsToPSS(report, test, awsSecretKey, awsAccessKey, dryRun)
		if err !=nil {
			return errors.Wrapf(err, "submit results to PSS %d of %d", idx+1, len(tests))
		}


		if err = uploadTests(ctx, client, report, test.SubTests, awsSecretKey, awsAccessKey, dryRun); err != nil {
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
