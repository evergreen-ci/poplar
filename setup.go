package poplar

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const (
	ProjectEnv      = "PROJECT"
	VersionEnv      = "VERSION"
	OrderEnv        = "ORDER"
	VariantEnv      = "VARIANT"
	TaskNameEnv     = "TASK_NAME"
	ExecutionEnv    = "EXECUTION"
	MainlineEnv     = "MAINLINE"
	APIKeyEnv       = "API_KEY"
	APISecretEnv    = "API_SECRET"
	APITokenEnv     = "API_TOKEN"
	BucketNameEnv   = "BUCKET_NAME"
	BucketPrefixEnv = "BUCKET_PREFIX"
	BucketRegionEnv = "BUCKET_REGION"
)

type ReportType string

const (
	ReportTypeJSON ReportType = "JSON"
	ReportTypeBSON ReportType = "BSON"
	ReportTypeYAML ReportType = "YAML"
	ReportTypeEnv  ReportType = "ENV"
)

func ReportSetup(reportType ReportType, filename string) (*Report, error) {
	switch reportType {
	case ReportTypeJSON:
		return reportSetupUnmarshal(reportType, filename, json.Unmarshal)
	case ReportTypeBSON:
		return reportSetupUnmarshal(reportType, filename, bson.Unmarshal)
	case ReportTypeYAML:
		return reportSetupUnmarshal(reportType, filename, yaml.Unmarshal)
	case ReportTypeEnv:
		return reportSetupEnv()
	default:
		return nil, errors.Errorf("invalid report type %s", reportType)
	}
}

func reportSetupUnmarshal(reportType ReportType, filename string, unmarshal func([]byte, interface{}) error) (*Report, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "problem opening file %s", filename)
	}

	report := &Report{}
	if err := unmarshal(data, report); err != nil {
		return nil, errors.Wrapf(err, "problem unmarshalling %s from %s", reportType, filename)
	}

	return report, nil
}

func reportSetupEnv() (*Report, error) {
	var order int
	var execution int
	var mainline bool
	var err error
	if os.Getenv(OrderEnv) != "" {
		order, err = strconv.Atoi(os.Getenv(OrderEnv))
		if err != nil {
			return nil, errors.Wrapf(err, "env var %s should be an int", OrderEnv)
		}
	}
	if os.Getenv(ExecutionEnv) != "" {
		execution, err = strconv.Atoi(os.Getenv(ExecutionEnv))
		if err != nil {
			return nil, errors.Wrapf(err, "env var %s should be an int", ExecutionEnv)
		}
	}
	if os.Getenv(MainlineEnv) != "" {
		mainline, err = strconv.ParseBool(os.Getenv(MainlineEnv))
		if err != nil {
			return nil, errors.Wrapf(err, "env var %s should be a bool", MainlineEnv)
		}
	}

	return &Report{
		Project:   os.Getenv(ProjectEnv),
		Version:   os.Getenv(VersionEnv),
		Order:     order,
		Variant:   os.Getenv(VariantEnv),
		TaskName:  os.Getenv(TaskNameEnv),
		Execution: execution,
		Mainline:  mainline,
		BucketConf: BucketConfiguration{
			APIKey:    os.Getenv(APIKeyEnv),
			APISecret: os.Getenv(APISecretEnv),
			APIToken:  os.Getenv(APITokenEnv),
			Name:      os.Getenv(BucketNameEnv),
			Prefix:    os.Getenv(BucketPrefixEnv),
			Region:    os.Getenv(BucketRegionEnv),
		},
		Tests: []Test{},
	}, nil
}
