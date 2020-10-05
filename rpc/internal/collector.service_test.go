package internal

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/evergreen-ci/poplar"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getTestCollectorService(tmpDir string) *collectorService {
	registry := poplar.NewRegistry()
	registry.Create("collector", poplar.CreateOptions{
		Path:     filepath.Join(tmpDir, "exists"),
		Recorder: poplar.RecorderPerf,
		Dynamic:  true,
		Events:   poplar.EventsCollectorBasic,
	})

	return &collectorService{
		registry: registry,
	}
}

func TestCreateCollector(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "create-collector-test")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpDir))
	}()
	svc := getTestCollectorService(tmpDir)

	for _, test := range []struct {
		name   string
		opts   *CreateOptions
		resp   *PoplarResponse
		hasErr bool
	}{
		{
			name: "CollectorDNE",
			opts: &CreateOptions{
				Name:     "new",
				Path:     filepath.Join(tmpDir, "new"),
				Recorder: CreateOptions_PERF,
				Events:   CreateOptions_BASIC,
			},
			resp: &PoplarResponse{Name: "new", Status: true},
		},
		{
			name: "CollectorExists",
			opts: &CreateOptions{
				Name:     "collector",
				Path:     filepath.Join(tmpDir, "exists"),
				Recorder: CreateOptions_PERF,
				Events:   CreateOptions_BASIC,
			},
			resp: &PoplarResponse{Name: "collector", Status: true},
		},
		{
			name: "InvalidOpts",
			opts: &CreateOptions{
				Name:   "invalid",
				Path:   filepath.Join(tmpDir, "invalid"),
				Events: CreateOptions_BASIC,
			},
			resp:   nil,
			hasErr: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp, err := svc.CreateCollector(context.TODO(), test.opts)
			assert.Equal(t, test.resp, resp)
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestCloseCollector(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "close-collector-test")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpDir))
	}()
	svc := getTestCollectorService(tmpDir)

	for _, test := range []struct {
		name string
		id   *PoplarID
		resp *PoplarResponse
	}{
		{
			name: "Exists",
			id:   &PoplarID{Name: "collector"},
			resp: &PoplarResponse{Name: "collector", Status: true},
		},
		{
			name: "DNE",
			id:   &PoplarID{Name: "dne"},
			resp: &PoplarResponse{Name: "dne", Status: true},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp, err := svc.CloseCollector(context.TODO(), test.id)
			assert.Equal(t, test.resp, resp)
			assert.NoError(t, err)
			_, ok := svc.registry.GetCollector(test.id.Name)
			assert.False(t, ok)
		})
	}
}

func TestSendEvent(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "send-event-test")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpDir))
	}()
	svc := getTestCollectorService(tmpDir)

	for _, test := range []struct {
		name   string
		event  *EventMetrics
		resp   *PoplarResponse
		hasErr bool
	}{
		{
			name:   "CollectorDNE",
			event:  &EventMetrics{Name: "DNE"},
			resp:   nil,
			hasErr: true,
		},
		{
			name: "AddEvent",
			event: &EventMetrics{
				Name: "collector",
				Time: &timestamp.Timestamp{},
				Timers: &EventMetricsTimers{
					Total:    &duration.Duration{},
					Duration: &duration.Duration{},
				},
			},
			resp: &PoplarResponse{Name: "collector", Status: true},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp, err := svc.SendEvent(context.TODO(), test.event)
			assert.Equal(t, test.resp, resp)
			if test.hasErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
