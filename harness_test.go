package poplar

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/mongodb/ftdc"
	"github.com/mongodb/ftdc/events"
	"github.com/mongodb/grip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCaseType(t *testing.T) {
	t.Run("SatifiesMinimum", func(t *testing.T) {
		t.Run("ZeroValue", func(t *testing.T) {
			c := BenchmarkCase{}
			r := &BenchmarkResult{}
			// no minimum specified
			assert.True(t, c.satisfiedMinimumRuntime(r))
			assert.False(t, c.satisfiedMinimumIterations(r))
			assert.False(t, c.satisfiedMinimumns(r))
		})
		t.Run("RuntimeUnfulfilled", func(t *testing.T) {
			c := BenchmarkCase{MinIterations: 2, MinRuntime: time.Hour}
			r := &BenchmarkResult{Iterations: 2, Runtime: time.Minute}
			assert.False(t, c.satisfiedMinimumRuntime(r))
			assert.True(t, c.satisfiedMinimumIterations(r))
			assert.False(t, c.satisfiedMinimumns(r))
		})
		t.Run("PassingCaseIterations", func(t *testing.T) {
			c := BenchmarkCase{MinIterations: 2}
			r := &BenchmarkResult{Iterations: 3}
			assert.True(t, c.satisfiedMinimumRuntime(r))
			assert.True(t, c.satisfiedMinimumIterations(r))
			assert.True(t, c.satisfiedMinimumns(r))
		})
		t.Run("PassingCase", func(t *testing.T) {
			c := BenchmarkCase{MinIterations: 2, MinRuntime: time.Minute}
			r := &BenchmarkResult{Iterations: 3, Runtime: 5 * time.Minute}
			assert.True(t, c.satisfiedMinimumRuntime(r))
			assert.True(t, c.satisfiedMinimumIterations(r))
			assert.True(t, c.satisfiedMinimumns(r))
		})
		t.Run("RuntimeOnlySatisfied", func(t *testing.T) {
			c := BenchmarkCase{MinRuntime: time.Minute, MinIterations: 0}
			r := &BenchmarkResult{Runtime: time.Hour, Iterations: 1}
			assert.True(t, c.satisfiedMinimumRuntime(r))
			assert.False(t, c.satisfiedMinimumIterations(r))
			assert.False(t, c.satisfiedMinimumns(r))
		})
	})
	t.Run("ExceedsMinimums", func(t *testing.T) {
		t.Run("ZeroValue", func(t *testing.T) {
			c := BenchmarkCase{}
			r := &BenchmarkResult{}
			// no minimum specified
			assert.False(t, c.exceededMaximumRuntime(r))
			assert.False(t, c.exceededMaximumIterations(r))
			assert.False(t, c.exceededMaximums(r))
		})
		t.Run("UnderRuntime", func(t *testing.T) {
			c := BenchmarkCase{MaxRuntime: time.Minute}
			r := &BenchmarkResult{Runtime: time.Second}
			// no minimum specified
			assert.False(t, c.exceededMaximumRuntime(r))
			assert.False(t, c.exceededMaximumIterations(r))
			assert.False(t, c.exceededMaximums(r))
		})
		t.Run("OverRuntime", func(t *testing.T) {
			c := BenchmarkCase{MaxRuntime: time.Minute}
			r := &BenchmarkResult{Runtime: time.Hour}
			// no minimum specified
			assert.True(t, c.exceededMaximumRuntime(r))
			assert.False(t, c.exceededMaximumIterations(r))
			assert.True(t, c.exceededMaximums(r))
		})
		t.Run("UnderIterations", func(t *testing.T) {
			c := BenchmarkCase{MaxIterations: 60}
			r := &BenchmarkResult{Iterations: 1}
			// no minimum specified
			assert.False(t, c.exceededMaximumRuntime(r))
			assert.False(t, c.exceededMaximumIterations(r))
			assert.False(t, c.exceededMaximums(r))
		})
		t.Run("OverIterations", func(t *testing.T) {
			c := BenchmarkCase{MaxIterations: 60}
			r := &BenchmarkResult{Iterations: 360}
			// no minimum specified
			assert.False(t, c.exceededMaximumRuntime(r))
			assert.True(t, c.exceededMaximumIterations(r))
			assert.True(t, c.exceededMaximums(r))
		})
	})
	t.Run("Validate", func(t *testing.T) {
		t.Run("SetDefault", func(t *testing.T) {
			c := BenchmarkCase{}
			assert.Zero(t, c)
			assert.Error(t, c.Validate())
			assert.NotZero(t, c)
			assert.Equal(t, RecorderPerf, c.Recorder)
			assert.Equal(t, 1, c.Count)
			assert.Equal(t, 10*time.Minute, c.Timeout)
			assert.Equal(t, 5*time.Minute, c.IterationTimeout)
		})
		for _, test := range []struct {
			Name string
			Case func(*testing.T, BenchmarkCase)
		}{
			{
				Name: "ValidateFixture",
				Case: func(t *testing.T, c BenchmarkCase) {
					assert.NoError(t, c.Validate())
				},
			},
			{
				Name: "InvalidIterationTimeout",
				Case: func(t *testing.T, c BenchmarkCase) {
					c.IterationTimeout = time.Hour
					c.Timeout = time.Minute
					assert.Error(t, c.Validate())
				},
			},
			{
				Name: "InvalidRuntimes",
				Case: func(t *testing.T, c BenchmarkCase) {
					c.MinRuntime = time.Hour
					c.MaxRuntime = time.Minute
					assert.Error(t, c.Validate())
				},
			},
			{
				Name: "Invaliditerations",
				Case: func(t *testing.T, c BenchmarkCase) {
					c.MinIterations = 1000
					c.MaxIterations = 100
					assert.Error(t, c.Validate())
				},
			},
			{
				Name: "MinIterationsMustBeSet",
				Case: func(t *testing.T, c BenchmarkCase) {
					c.MinIterations = 0
					assert.Error(t, c.Validate())
				},
			},
			{
				Name: "MinRunttimeMustBeSet",
				Case: func(t *testing.T, c BenchmarkCase) {
					c.MinRuntime = 0
					assert.Error(t, c.Validate())
				},
			},
			{
				Name: "StringFormat",
				Case: func(t *testing.T, c BenchmarkCase) {
					assert.NotZero(t, c.String())
					assert.Contains(t, c.String(), "name=1,")
					c.CaseName = "foo"
					assert.Contains(t, c.String(), c.CaseName)
				},
			},
		} {
			t.Run(test.Name, func(t *testing.T) {
				c := BenchmarkCase{
					Bench:         func(_ context.Context, _ events.Recorder, _ int) error { return nil },
					MinIterations: 10,
					MinRuntime:    time.Minute,
					MaxRuntime:    time.Hour,
					MaxIterations: 100,
				}
				require.NoError(t, c.Validate())
				test.Case(t, c)
			})
		}
	})
	t.Run("Fluent", func(t *testing.T) {
		for _, test := range []struct {
			Name string
			Case func(*testing.T, *BenchmarkCase)
		}{
			{
				Name: "ValidateFixture",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.NotZero(t, c)
					assert.Zero(t, *c)
					assert.Error(t, c.Validate())
				},
			},
			{
				Name: "NameRoundTrip",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.Equal(t, "foo", c.SetName("foo").Name())
				},
			},
			{
				Name: "Recorder",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.Equal(t, RecorderType(RecorderPerfSingle), c.SetRecorder(RecorderPerfSingle).Recorder)
				},
			},
			{
				Name: "Benchmark",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.NotZero(t, c.SetBench(func(_ context.Context, _ events.Recorder, _ int) error { return nil }).Bench)
				},
			},
			{
				Name: "Count",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.Equal(t, 42, c.SetCount(42).Count)
				},
			},
			{
				Name: "MaxDuration",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.Equal(t, time.Minute, c.SetMaxDuration(time.Minute).MaxRuntime)
				},
			},
			{
				Name: "MaxIterations",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.Equal(t, 42, c.SetMaxIterations(42).MaxIterations)
				},
			},
			{
				Name: "IterationTimeout",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.Equal(t, time.Minute, c.SetIterationTimeout(time.Minute).IterationTimeout)
				},
			},
			{
				Name: "Timeout",
				Case: func(t *testing.T, c *BenchmarkCase) {
					assert.Equal(t, time.Minute, c.SetTimeout(time.Minute).Timeout)
				},
			},
			{
				Name: "Iterations",
				Case: func(t *testing.T, c *BenchmarkCase) {
					c = c.SetIterations(10)
					assert.Equal(t, 10, c.MinIterations)
					assert.Equal(t, 100, c.MaxIterations)
				},
			},
			{
				Name: "SetDurationHour",
				Case: func(t *testing.T, c *BenchmarkCase) {
					c = c.SetDuration(time.Hour)
					assert.Equal(t, time.Hour, c.MinRuntime)
					assert.Equal(t, 61*time.Minute, c.MaxRuntime)
				},
			},
			{
				Name: "SetDurationMinute",
				Case: func(t *testing.T, c *BenchmarkCase) {
					c = c.SetDuration(time.Minute)
					assert.Equal(t, time.Minute, c.MinRuntime)
					assert.Equal(t, 10*time.Minute, c.MaxRuntime)
				},
			},
			{
				Name: "SetDurationSecond",
				Case: func(t *testing.T, c *BenchmarkCase) {
					c = c.SetDuration(time.Second)
					assert.Equal(t, time.Second, c.MinRuntime)
					assert.Equal(t, 10*time.Second, c.MaxRuntime)
				},
			},
			{
				Name: "SetDurationTwoMinute",
				Case: func(t *testing.T, c *BenchmarkCase) {
					c = c.SetDuration(2 * time.Minute)
					assert.Equal(t, 2*time.Minute, c.MinRuntime)
					assert.Equal(t, 4*time.Minute, c.MaxRuntime)
				},
			},
			{
				Name: "Standard",
				Case: func(t *testing.T, c *BenchmarkCase) {
					c.SetBench(func(ctx context.Context, _ events.Recorder, count int) error { return errors.New("is nil") })
					assert.NotZero(t, c.Bench)
					assert.NotNil(t, c.Bench.Standard(nil))
					assert.NotNil(t, c.Standard(nil))
					assert.Panics(t, func() { c.Standard(nil)(nil) })
					assert.NotPanics(t, func() {
						res := testing.Benchmark(c.Standard(nil))
						assert.Zero(t, res)
					})

				},
			},
		} {
			t.Run(test.Name, func(t *testing.T) {
				test.Case(t, &BenchmarkCase{})
			})
		}
	})
	t.Run("Execution", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
		defer cancel()

		collector := ftdc.NewBaseCollector(10)
		recorder := events.NewSingleRecorder(collector)

		t.Run("NoopBenchmark", func(t *testing.T) {
			c := BenchmarkCase{
				Bench:            func(ctx context.Context, _ events.Recorder, count int) error { return nil },
				MinIterations:    2,
				MaxIterations:    4,
				MinRuntime:       time.Millisecond,
				MaxRuntime:       time.Second,
				Timeout:          time.Minute,
				IterationTimeout: time.Microsecond,
			}
			require.NoError(t, c.Validate())
			res := c.Run(ctx, recorder)
			assert.NoError(t, res.Error)
			assert.True(t, res.Runtime >= time.Millisecond)
			assert.True(t, res.Runtime <= time.Second)
			assert.True(t, res.Iterations > c.MinIterations)
			assert.False(t, res.Iterations < c.MaxIterations)

			assert.True(t, res.CompletedAt.Sub(res.StartAt) < time.Second)
			assert.True(t, res.CompletedAt.Sub(res.StartAt) > time.Millisecond)

		})
		t.Run("CanceledContext", func(t *testing.T) {
			c := BenchmarkCase{
				Bench: func(ctx context.Context, _ events.Recorder, count int) error {
					return nil
				},
				MinIterations:    2,
				MaxIterations:    4,
				MinRuntime:       time.Minute,
				MaxRuntime:       time.Hour,
				Timeout:          2 * time.Hour,
				IterationTimeout: 20 * time.Second,
			}
			require.NoError(t, c.Validate())
			ctx, cancel := context.WithCancel(ctx)
			cancel()
			res := c.Run(ctx, recorder)
			assert.NoError(t, res.Error)
			assert.Zero(t, res.Iterations)
			assert.True(t, res.CompletedAt.Sub(res.StartAt) < time.Second)
			grip.Info(res)
		})
		t.Run("Error", func(t *testing.T) {
			err := errors.New("foo")
			c := BenchmarkCase{
				Bench: func(ctx context.Context, _ events.Recorder, count int) error {
					return err
				},
				MinIterations:    2,
				MaxIterations:    4,
				MinRuntime:       time.Minute,
				MaxRuntime:       time.Hour,
				Timeout:          2 * time.Hour,
				IterationTimeout: 20 * time.Second,
			}
			require.NoError(t, c.Validate())

			res := c.Run(ctx, recorder)
			assert.Error(t, res.Error)
			assert.Equal(t, err, res.Error)
			assert.Zero(t, res.Iterations)
			assert.True(t, res.CompletedAt.Sub(res.StartAt) < time.Millisecond)
		})
	})
}
