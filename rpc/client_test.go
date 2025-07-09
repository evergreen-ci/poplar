package rpc

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/evergreen-ci/pail"
	"github.com/evergreen-ci/poplar"
	"github.com/evergreen-ci/utility"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func mockUploadReport(ctx context.Context, report *poplar.Report, serialize bool, dryRun bool) error {
	opts := UploadReportOptions{
		Report:          report,
		SerializeUpload: serialize,
		DryRun:          dryRun,
	}
	if err := opts.convertAndUploadArtifacts(ctx); err != nil {
		return errors.Wrap(err, "converting and uploading artifacts for report")
	}
	return nil
}

func TestClient(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	testdataDir := filepath.Join("..", "testdata")
	s3Name := "build-test-curator"
	s3Prefix := "poplar-client-test"
	s3Opts := pail.S3Options{
		Name:   s3Name,
		Prefix: s3Prefix,
		Region: "us-east-1",
	}

	client := utility.GetHTTPClient()
	defer utility.PutHTTPClient(client)

	s3Bucket, err := pail.NewS3BucketWithHTTPClient(ctx, client, s3Opts)
	require.NoError(t, err)

	report := generateTestReport(testdataDir, s3Name, s3Prefix, false)
	expectedTests := []poplar.Test{
		report.Tests[0],
		report.Tests[0].SubTests[0],
		report.Tests[1],
		report.Tests[1].SubTests[0],
	}
	for i := range expectedTests {
		for j := range expectedTests[i].Artifacts {
			require.NoError(t, expectedTests[i].Artifacts[j].Convert(ctx))
			require.NoError(t, expectedTests[i].Artifacts[j].SetBucketInfo(report.BucketConf))
			require.NoError(t, os.RemoveAll(filepath.Join(testdataDir, expectedTests[i].Artifacts[j].Path)))
		}
	}

	defer func() {
		for _, test := range expectedTests {
			for _, artifact := range test.Artifacts {
				assert.NoError(t, s3Bucket.Remove(ctx, artifact.Path))
				assert.NoError(t, os.RemoveAll(filepath.Join(testdataDir, artifact.Path)))
			}
		}
	}()
	t.Run("WetRun", func(t *testing.T) {
		for _, serialize := range []bool{true, false} {
			testReport := generateTestReport(testdataDir, s3Name, s3Prefix, false)
			require.NoError(t, mockUploadReport(ctx, &testReport, serialize, false))
		}
	})

	for _, test := range expectedTests {
		for _, artifact := range test.Artifacts {
			require.NoError(t, s3Bucket.Remove(ctx, artifact.Path))
			require.NoError(t, os.RemoveAll(filepath.Join(testdataDir, artifact.Path)))
		}
	}

	t.Run("DryRun", func(t *testing.T) {
		for _, serialize := range []bool{true, false} {
			testReport := generateTestReport(testdataDir, s3Name, s3Prefix, false)
			require.NoError(t, mockUploadReport(ctx, &testReport, serialize, true))
		}
	})

	for _, test := range expectedTests {
		for _, artifact := range test.Artifacts {
			require.NoError(t, s3Bucket.Remove(ctx, artifact.Path))
			require.NoError(t, os.RemoveAll(filepath.Join(testdataDir, artifact.Path)))
		}
	}

	t.Run("DuplicateMetricName", func(t *testing.T) {
		for _, serialize := range []bool{true, false} {
			testReport := generateTestReport(testdataDir, s3Name, s3Prefix, true)
			assert.Error(t, mockUploadReport(ctx, &testReport, serialize, true))
		}
	})
}

func generateTestReport(testdataDir, s3Name, s3Prefix string, duplicateMetric bool) poplar.Report {
	report := poplar.Report{
		Project:   "project",
		Version:   "version",
		Order:     2,
		Variant:   "variant",
		TaskName:  "taskName",
		TaskID:    "taskID",
		Mainline:  true,
		Execution: 2,

		BucketConf: poplar.BucketConfiguration{
			Name:   s3Name,
			Region: "us-east-1",
		},

		Tests: []poplar.Test{
			{
				Info: poplar.TestInfo{
					TestName:  "test0",
					Trial:     2,
					Tags:      []string{"tag0", "tag1"},
					Arguments: map[string]int32{"thread_level": 1},
				},
				Artifacts: []poplar.TestArtifact{
					{
						Bucket:           s3Name,
						Prefix:           s3Prefix,
						Path:             "bson_example.ftdc",
						LocalFile:        filepath.Join(testdataDir, "bson_example.bson"),
						ConvertBSON2FTDC: true,
					},
					{
						Prefix:      s3Prefix,
						LocalFile:   filepath.Join(testdataDir, "bson_example.bson"),
						ConvertGzip: true,
					},
				},
				CreatedAt:   time.Date(2018, time.July, 4, 12, 0, 0, 0, time.UTC),
				CompletedAt: time.Date(2018, time.July, 4, 12, 1, 0, 0, time.UTC),
				SubTests: []poplar.Test{
					{
						Info: poplar.TestInfo{
							TestName: "test00",
						},
						Metrics: []poplar.TestMetrics{
							{
								Name:    "mean",
								Version: 1,
								Value:   1.5,
								Type:    "MEAN",
							},
							{
								Name:    "sum",
								Version: 1,
								Value:   10,
								Type:    "SUM",
							},
						},
					},
				},
			},
			{
				Info: poplar.TestInfo{
					TestName: "test1",
				},
				Artifacts: []poplar.TestArtifact{
					{
						Bucket:           s3Name,
						Prefix:           s3Prefix,
						Path:             "json_example.ftdc",
						LocalFile:        filepath.Join(testdataDir, "json_example.json"),
						CreatedAt:        time.Date(2018, time.July, 4, 11, 59, 0, 0, time.UTC),
						ConvertJSON2FTDC: true,
					},
				},
				SubTests: []poplar.Test{
					{
						Info: poplar.TestInfo{
							TestName: "test10",
						},
						Metrics: []poplar.TestMetrics{
							{
								Name:    "mean",
								Version: 1,
								Value:   1.5,
								Type:    "MEAN",
							},
							{
								Name:    "sum",
								Version: 1,
								Value:   10,
								Type:    "SUM",
							},
						},
					},
				},
			},
		},
	}

	if duplicateMetric {
		report.Tests[0].SubTests[0].Metrics = append(
			report.Tests[0].SubTests[0].Metrics,
			poplar.TestMetrics{
				Name:    "mean",
				Version: 1,
				Value:   2,
				Type:    "MEAN",
			},
		)
	}

	return report
}
