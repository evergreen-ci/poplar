package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/evergreen-ci/aviation"
	"github.com/evergreen-ci/poplar/collector"
	"github.com/golang/protobuf/ptypes"
	"github.com/mongodb/grip"
	"google.golang.org/grpc"
)

const numSeconds = 10

func main() {
	grip.Info("hello world!")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conn, err := grpc.DialContext(ctx, "localhost:2288",
		grpc.WithUnaryInterceptor(aviation.MakeRetryUnaryClientInterceptor(10)),
		grpc.WithStreamInterceptor(aviation.MakeRetryStreamClientInterceptor(10)),
		grpc.WithInsecure())

	grip.EmergencyFatal(err)

	name := fmt.Sprintf("poplar-%d", time.Now().Unix())

	client := collector.NewPoplarEventCollectorClient(conn)
	_, err = client.CreateCollector(ctx, &collector.CreateOptions{
		Name:      name,
		Path:      name + ".ftdc",
		Streaming: true,
		Dynamic:   false,
		ChunkSize: 100,
		Recorder:  collector.CreateOptions_PERF,
	})
	grip.EmergencyFatal(err)
	var events int64
	startAt := time.Now()

	if os.Getenv("POPLAR_SINGLE") != "" {
		for {
			_, err = client.SendEvent(ctx, &collector.EventMetrics{
				Name:     name,
				Id:       events,
				Time:     ptypes.TimestampNow(),
				Counters: &collector.EventMetricsCounters{},
				Timers:   &collector.EventMetricsTimers{},
				Gauges:   &collector.EventMetricsGauges{},
			})
			grip.EmergencyFatal(err)
			events++
			if time.Since(startAt) > numSeconds*time.Second {
				break
			}
		}
	} else {
		srv, err := client.StreamEvents(ctx)
		grip.EmergencyFatal(err)
		for {
			err = srv.Send(&collector.EventMetrics{
				Name:     name,
				Id:       events,
				Time:     ptypes.TimestampNow(),
				Counters: &collector.EventMetricsCounters{},
				Timers:   &collector.EventMetricsTimers{},
				Gauges:   &collector.EventMetricsGauges{},
			})
			events++

			if err != nil || time.Since(startAt) > numSeconds*time.Second {
				_, err = srv.CloseAndRecv()
				grip.EmergencyFatal(err)
				break
			}
		}

	}

	_, err = client.CloseCollector(ctx, &collector.PoplarID{Name: name})
	grip.EmergencyFatal(err)

	dur := time.Since(startAt)
	grip.Infof("%s.ftdc - %d events per %s [%d ops/sec]", name, events, dur, events/numSeconds)

}
