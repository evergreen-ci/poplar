package internal

import (
	"bytes"
	"container/list"
	"context"
	fmt "fmt"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/evergreen-ci/poplar"
	"github.com/mongodb/ftdc"
	"github.com/mongodb/grip"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestCreateCollector(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "create-collector-test")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpDir))
	}()
	svc := getTestCollectorService(t, tmpDir)
	defer func() {
		assert.NoError(t, closeCollectorService(svc))
	}()

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
			if test.hasErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, test.resp.Name, resp.Name)
				assert.Equal(t, test.resp.Status, resp.Status)
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
	svc := getTestCollectorService(t, tmpDir)
	svc.coordinator.groups["group"] = &streamGroup{
		streams: map[string]*stream{
			"id1": {
				buffer: &list.List{},
			},
		},
		eventHeap: &PerformanceHeap{},
	}
	defer func() {
		assert.NoError(t, closeCollectorService(svc))
	}()

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
			assert.Equal(t, test.resp.Name, resp.Name)
			assert.Equal(t, test.resp.Status, resp.Status)
			assert.NoError(t, err)
			_, ok := svc.registry.GetCollector(test.id.Name)
			assert.False(t, ok)
			_, ok = svc.coordinator.groups["group"].streams["id1"]
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
	svc := getTestCollectorService(t, tmpDir)
	defer func() {
		assert.NoError(t, closeCollectorService(svc))
	}()

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
				Time: &timestamppb.Timestamp{},
				Timers: &EventMetricsTimers{
					Total:    &durationpb.Duration{},
					Duration: &durationpb.Duration{},
				},
			},
			resp: &PoplarResponse{Name: "collector", Status: true},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp, err := svc.SendEvent(context.TODO(), test.event)
			if test.hasErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, err)
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, test.resp.Name, resp.Name)
				assert.Equal(t, test.resp.Status, resp.Status)
			}
		})
	}
}

func TestRegisterStream(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "register-stream-test")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpDir))
	}()
	svc := getTestCollectorService(t, tmpDir)
	defer func() {
		assert.NoError(t, closeCollectorService(svc))
	}()

	for _, test := range []struct {
		name          string
		collectorName *CollectorName
		resp          *PoplarResponse
		hasErr        bool
	}{
		{
			name:          "EmptyName",
			collectorName: &CollectorName{},
			resp:          nil,
			hasErr:        true,
		},
		{
			name:          "CollectorDNE",
			collectorName: &CollectorName{Name: "DNE"},
			resp:          nil,
			hasErr:        true,
		},
		{
			name:          "CollectorExists",
			collectorName: &CollectorName{Name: "collector"},
			resp:          &PoplarResponse{Name: "collector", Status: true},
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			resp, registerErr := svc.RegisterStream(context.TODO(), test.collectorName)
			id, group, getErr := svc.coordinator.getStream(test.collectorName.Name)
			if test.hasErr {
				assert.Error(t, registerErr)
				assert.Error(t, getErr)
				assert.Nil(t, resp)
			} else {
				require.NoError(t, registerErr)
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, test.resp.Name, resp.Name)
				assert.Equal(t, test.resp.Status, resp.Status)
				assert.NotEmpty(t, id)
				assert.NotNil(t, group)
				assert.NoError(t, getErr)
			}
		})
	}
}

func TestStreamEvent(t *testing.T) {
	tmpDir, err := ioutil.TempDir(".", "stream-event-test")
	require.NoError(t, err)
	defer func() {
		assert.NoError(t, os.RemoveAll(tmpDir))
	}()
	svc := getTestCollectorService(t, tmpDir)
	defer func() {
		assert.NoError(t, closeCollectorService(svc))
	}()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	address := fmt.Sprintf("localhost:%d", 7070)
	require.NoError(t, startCollectorService(ctx, svc, address))
	client, err := getCollectorGRPCClient(ctx, address, []grpc.DialOption{grpc.WithInsecure()})
	require.NoError(t, err)

	for _, test := range []struct {
		name     string
		events   []*EventMetrics
		resp     *PoplarResponse
		register bool
		hasErr   bool
	}{
		{
			name:   "Unregistered",
			events: []*EventMetrics{{Name: "DNE"}},
			resp:   nil,
			hasErr: true,
		},
		{
			name:   "NoName",
			events: []*EventMetrics{{}},
			resp:   nil,
			hasErr: true,
		},
		{
			name: "DifferentNames",
			events: []*EventMetrics{
				{
					Name: "collector",
					Time: &timestamppb.Timestamp{},
					Timers: &EventMetricsTimers{
						Total:    &durationpb.Duration{},
						Duration: &durationpb.Duration{},
					},
				},
				{
					Name: "anotherName",
					Time: &timestamppb.Timestamp{},
					Timers: &EventMetricsTimers{
						Total:    &durationpb.Duration{},
						Duration: &durationpb.Duration{},
					},
				},
			},
			resp:     nil,
			register: true,
			hasErr:   true,
		},
		{
			name: "AddEvents",
			events: []*EventMetrics{
				{
					Name: "collector",
					Time: &timestamppb.Timestamp{},
					Timers: &EventMetricsTimers{
						Total:    &durationpb.Duration{},
						Duration: &durationpb.Duration{},
					},
				},
				{
					Name: "collector",
					Time: &timestamppb.Timestamp{},
					Timers: &EventMetricsTimers{
						Total:    &durationpb.Duration{},
						Duration: &durationpb.Duration{},
					},
				},
			},
			resp:     &PoplarResponse{Name: "collector", Status: true},
			register: true,
		},
	} {
		t.Run(test.name, func(t *testing.T) {
			stream, err := client.StreamEvents(ctx)
			require.NoError(t, err)

			if test.register {
				_, err = client.RegisterStream(ctx, &CollectorName{Name: "collector"})
				require.NoError(t, err)
			}

			catcher := grip.NewBasicCatcher()
			for i := 0; i < len(test.events); i++ {
				catcher.Add(stream.Send(test.events[i]))
			}
			resp, err := stream.CloseAndRecv()
			catcher.Add(err)

			if test.hasErr {
				assert.Error(t, catcher.Resolve())
				assert.Nil(t, resp)
			} else {
				require.NoError(t, catcher.Resolve())
				require.NoError(t, err)
				require.NotNil(t, resp)
				assert.Equal(t, test.resp.Name, resp.Name)
				assert.Equal(t, test.resp.Status, resp.Status)
			}
		})
	}
	t.Run("MultipleStreams", func(t *testing.T) {
		event := EventMetrics{
			Name: "multiple",
			Time: &timestamppb.Timestamp{},
			Timers: &EventMetricsTimers{
				Total:    &durationpb.Duration{},
				Duration: &durationpb.Duration{},
			},
		}
		streams := make([]PoplarEventCollector_StreamEventsClient, 4)
		for i := range streams {
			var err error
			_, err = client.RegisterStream(ctx, &CollectorName{Name: "multiple"})
			require.NoError(t, err)
			streams[i], err = client.StreamEvents(ctx)
			require.NoError(t, err)

		}

		// Add one event to the first stream and close it. This will
		// check that streams closed early are not flushed until all
		// other streams have received at least one item.
		event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(24 * time.Hour).Unix()}
		require.NoError(t, streams[0].Send(&event))
		_, err := streams[0].CloseAndRecv()
		require.NoError(t, err)
		streams = streams[1:]

		for i := range streams {
			event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(time.Duration(i+3) * -time.Minute).Unix()}
			require.NoError(t, streams[i].Send(&event))
			event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(time.Duration(i+2) * -time.Minute).Unix()}
			require.NoError(t, streams[i].Send(&event))
			event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(time.Duration(i+1) * -time.Minute).Unix()}
			require.NoError(t, streams[i].Send(&event))
		}
		event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(-30 * time.Second).Unix()}
		require.NoError(t, streams[0].Send(&event))
		for i := range streams {
			event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(time.Duration(i+3) * -time.Second).Unix()}
			require.NoError(t, streams[i].Send(&event))
			event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(time.Duration(i+2) * -time.Second).Unix()}
			require.NoError(t, streams[i].Send(&event))
			event.Time = &timestamppb.Timestamp{Seconds: time.Now().Add(time.Duration(i+1) * -time.Second).Unix()}
			require.NoError(t, streams[i].Send(&event))

			_, err := streams[i].CloseAndRecv()
			require.NoError(t, err)
		}

		collector, ok := svc.registry.GetEventsCollector("multiple")
		require.True(t, ok)
		data, err := collector.Resolve()
		require.NoError(t, err)
		chunkIt := ftdc.ReadChunks(context.TODO(), bytes.NewReader(data))
		defer chunkIt.Close()

		count := 0
		var lastTS int64
		for i := 0; chunkIt.Next(); i++ {
			chunk := chunkIt.Chunk()
			for _, metric := range chunk.Metrics {
				if metric.Key() == "ts" {
					require.NotEmpty(t, metric.Values)
					firstTS := metric.Values[0]
					// The first TS of this chunk should be
					// greater than or equal to the last
					// TS of the previous chunk.
					require.True(t, firstTS > lastTS || firstTS == lastTS)

					var lastVal int64
					for _, val := range metric.Values {
						require.True(t, lastVal <= val)
						lastVal = val
						count++
					}
					lastTS = lastVal
				}
			}
		}
		assert.Equal(t, 20, count)
	})
}

func getTestCollectorService(t *testing.T, tmpDir string) *collectorService {
	registry := poplar.NewRegistry()
	_, err := registry.Create("collector", poplar.CreateOptions{
		Path:      filepath.Join(tmpDir, "exists"),
		ChunkSize: 5,
		Recorder:  poplar.RecorderPerf,
		Dynamic:   true,
		Events:    poplar.EventsCollectorBasic,
	})
	require.NoError(t, err)
	_, err = registry.Create("multiple", poplar.CreateOptions{
		Path:      filepath.Join(tmpDir, "multiple"),
		ChunkSize: 5,
		Recorder:  poplar.RecorderPerf,
		Dynamic:   true,
		Events:    poplar.EventsCollectorBasic,
	})
	require.NoError(t, err)

	return &collectorService{
		registry: registry,
		coordinator: &streamsCoordinator{
			groups: map[string]*streamGroup{},
		},
	}
}

func closeCollectorService(svc *collectorService) error {
	catcher := grip.NewBasicCatcher()
	catcher.Add(svc.registry.Close("collector"))
	catcher.Add(svc.registry.Close("multiple"))
	catcher.Add(svc.registry.Close("new"))

	return catcher.Resolve()
}

func startCollectorService(ctx context.Context, svc *collectorService, address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return errors.WithStack(err)
	}

	s := grpc.NewServer()
	RegisterPoplarEventCollectorServer(s, svc)

	go func() {
		_ = s.Serve(lis)
	}()
	go func() {
		<-ctx.Done()
		s.Stop()
	}()

	return nil
}

func getCollectorGRPCClient(ctx context.Context, address string, opts []grpc.DialOption) (PoplarEventCollectorClient, error) {
	conn, err := grpc.DialContext(ctx, address, opts...)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	go func() {
		<-ctx.Done()
		_ = conn.Close()
	}()

	return NewPoplarEventCollectorClient(conn), nil
}
