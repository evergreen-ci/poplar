package rpc

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"runtime"
	"sync"

	"github.com/evergreen-ci/juniper/gopb"
	"github.com/evergreen-ci/poplar"
	"github.com/evergreen-ci/poplar/rpc/internal"
	"github.com/evergreen-ci/utility"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/grip/recovery"
	"github.com/pkg/errors"
)

// UploadReportOptions contains all the required options for uploading a report
// to both Cedar and SPS.

type UploadReportOptions struct {
	Report          *poplar.Report
	HTTPClient      *http.Client
	SPSURL          string
	SerializeUpload bool
	DryRun          bool
}

// UploadReport does the following:
// 1. Ingest the report.json.
// 2. Upload the report to SPS over a REST endpoint.
func UploadReport(ctx context.Context, opts UploadReportOptions) error {
	if err := opts.convertAndUploadArtifacts(ctx); err != nil {
		return errors.Wrap(err, "uploading tests for report")
	}
	if !opts.DryRun {
		if opts.HTTPClient == nil {
			opts.HTTPClient = utility.GetHTTPClient()
			defer utility.PutHTTPClient(opts.HTTPClient)
		}
		if err := uploadTestsToSPS(ctx, opts.Report, opts.HTTPClient, opts.SPSURL); err != nil {
			return errors.Wrap(err, "uploading metrics for report to SPS")
		}
	}
	return nil
}

func uploadTestsToSPS(ctx context.Context, report *poplar.Report, client *http.Client, spsURL string) error {
	taskInformation := &internal.TaskInformation{
		Project:   report.Project,
		Version:   report.Version,
		Variant:   report.Variant,
		Order:     report.Order,
		TaskName:  report.TaskName,
		TaskId:    report.TaskID,
		Execution: report.Execution,
		Mainline:  report.Mainline,
	}
	perfResults := internal.GatherPerfResults(report)
	requestBody := &internal.SubmitPerfResultRequest{
		Id:      *taskInformation,
		Results: perfResults,
	}
	marshalledBody, err := json.Marshal(requestBody)
	if err != nil {
		return errors.Wrap(err, "marshalling request body")
	}
	bodyReader := bytes.NewReader(marshalledBody)
	requestURL := spsURL + "/raw_perf_results"
	request, err := http.NewRequestWithContext(ctx, "POST", requestURL, bodyReader)
	if err != nil {
		return errors.Wrap(err, "creating request")
	}
	response, err := client.Do(request)
	if err != nil {
		return errors.Wrap(err, "sending request")
	}
	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err == nil {
			return errors.Errorf("unexpected status code: %d: %s", response.StatusCode, body)
		} else {
			return errors.Wrapf(err, "unexpected status code: %d", response.StatusCode)
		}

	}
	return nil
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
