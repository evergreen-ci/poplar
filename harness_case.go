package poplar

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mongodb/ftdc/events"
	"github.com/pkg/errors"
)

// BenchmarkCase describes a single benchmark, and describes how to
// run a benchmark, including minimum and maximum runtimes and
// iterations.
//
// With poplar's exceution, via the Run method, cases will execute
// until both the minimum runtime and iteration count are reached, and
// will end as soon as either the maximum iteration or runtime counts
// are exceeded.
//
// You can also use the Standard() function to convert the
// BenchmarkCase into a more conventional go standard library
// Bencharmk function.
type BenchmarkCase struct {
	CaseName      string
	Bench         Benchmark
	MinRuntime    time.Duration
	MaxRuntime    time.Duration
	Count         int
	MinIterations int
	MaxIterations int
	Recorder      RecorderType
}

// Benchmark defines a function signature for running a benchmark
// test.
//
// These functions take a context to support course timeouts, a
// Recorder instance to capture intra-test data, and a count number to
// tell the test the number of times the function should run.
//
// In general, functions should resemble the following:
//
//    func(ctx context.Context, r events.Recorder, count int) error {
//         ticks := 0
//         for i :=0, i < count; i++ {
//             r.Begin()
//             ticks += 4
//             r.IncOps(4)
//             r.End()
//         }
//    }
//
type Benchmark func(context.Context, events.Recorder, int) error

// Standard converts Benchmark function into a standard library test
// function, using the recorder to capture events. The count value
// passed to the benchmark function is b.N.
func (bench Benchmark) Standard(recorder events.Recorder) func(*testing.B) {
	return func(b *testing.B) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		shim := events.NewShimRecorder(recorder, b)
		b.ResetTimer()
		err := bench(ctx, shim, b.N)
		b.StopTimer()
		if err != nil {
			b.Fatal(errors.Wrap(err, "benchmark failed"))
		}
	}
}

// Standard produces a standard library test function, as a
// passthrough for BenchmarkCase.Bench.Standard.
func (c *BenchmarkCase) Standard(r events.Recorder) func(*testing.B) { return c.Bench.Standard(r) }

// Name returns either the CaseName value OR the name of the symbol
// for the benchmark function. Use CaseName when you define the case
// as a function literal, defaulting
func (c *BenchmarkCase) Name() string {
	if c.CaseName != "" {
		return c.CaseName
	}
	return getName(c.Bench)
}

func (c *BenchmarkCase) SetName(n string) *BenchmarkCase                 { c.CaseName = n; return c }
func (c *BenchmarkCase) SetRecorder(r RecorderType) *BenchmarkCase       { c.Recorder = r; return c }
func (c *BenchmarkCase) SetBench(b Benchmark) *BenchmarkCase             { c.Bench = b; return c }
func (c *BenchmarkCase) SetCount(v int) *BenchmarkCase                   { c.Count = v; return c }
func (c *BenchmarkCase) SetMaxDuration(dur time.Duration) *BenchmarkCase { c.MaxRuntime = dur; return c }
func (c *BenchmarkCase) SetMaxIterations(v int) *BenchmarkCase           { c.MaxIterations = v; return c }

func (c *BenchmarkCase) SetDuration(dur time.Duration) *BenchmarkCase {
	c.MinRuntime = dur
	if dur <= time.Minute {
		c.MaxRuntime = 10 * dur
	} else if dur <= 10*time.Minute {
		c.MaxRuntime = 2 * dur
	} else {
		c.MaxRuntime = time.Minute + dur
	}
	return c
}

func (c *BenchmarkCase) SetIterations(v int) *BenchmarkCase {
	c.MinIterations = v
	c.MaxIterations = 10 * v
	return c
}

func (c *BenchmarkCase) String() string {
	return fmt.Sprintf("name=%s, count=%d, min_dur=%s, max_dur=%s, min_iters=%s, max_iters=%d",
		c.Name(), c.Count, c.MinRuntime, c.MaxRuntime, c.MinIterations, c.MaxIterations)
}

func (c *BenchmarkCase) Validate() error {
	if c.Recorder == "" {
		c.Recorder = RecorderPerf
	}

	if c.Count == 0 {
		c.Count = 1
	}

	if c.Bench == nil {
		return errors.New("must specify a valid benchmark")
	}

	if c.MinRuntime >= c.MaxRuntime {
		return errors.New("min runtime must not be >= max runtime ")
	}

	if c.MinIterations >= c.MaxIterations {
		return errors.New("min iterations must not be >= max iterations")
	}

	return nil
}

func (c *BenchmarkCase) Run(ctx context.Context, recorder events.Recorder) BenchmarkResult {
	res := &BenchmarkResult{
		Iterations: 1,
		Count:      c.Count,
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 2*c.MaxRuntime)
	defer cancel()

	for {
		switch {
		case ctx.Err() != nil:
			break
		case c.satisfiedMinimumRuntime(res) && c.satisfiedMinimumIterations(res):
			break
		case c.exceededMaximumRuntime(res) || c.exceededMaximumIterations(res):
			break
		default:
			startAt := time.Now()
			bctx, bcancel := context.WithCancel(ctx)
			res.Error = c.Bench(bctx, recorder, c.Count)
			res.Runtime += time.Since(startAt)
			res.Iterations++
			bcancel()

			if res.Error != nil {
				recorder.SetFailed(true)
				break
			}
		}
	}

	return *res
}

func (c *BenchmarkCase) satisfiedMinimumRuntime(res *BenchmarkResult) bool {
	return (c.MinRuntime > 0 && c.MaxRuntime <= 0) && res.Runtime >= c.MinRuntime
}

func (c *BenchmarkCase) exceededMaximumRuntime(res *BenchmarkResult) bool {
	return (c.MaxRuntime > 0 && res.Runtime >= c.MaxRuntime) && (c.MinIterations == 0 || res.Iterations > c.MinIterations)
}

func (c *BenchmarkCase) satisfiedMinimumIterations(res *BenchmarkResult) bool {
	return (c.MinIterations > 0 && c.MaxIterations <= 0) && res.Iterations >= c.MinIterations
}

func (c *BenchmarkCase) exceededMaximumIterations(res *BenchmarkResult) bool {
	return (c.MaxIterations > 0 && res.Iterations >= c.MaxIterations) && (c.MinRuntime == 0 || res.Runtime > c.MinRuntime)
}
