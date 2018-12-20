package poplar

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/mongodb/ftdc/events"
	"github.com/mongodb/grip"
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
	CaseName         string
	Bench            Benchmark
	MinRuntime       time.Duration
	MaxRuntime       time.Duration
	Timeout          time.Duration
	IterationTimeout time.Duration
	Count            int
	MinIterations    int
	MaxIterations    int
	Recorder         RecorderType
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
// for the benchmark function. Use the CaseName field/SetName function
// when you define the case as a function literal, or to override the
// function name.
func (c *BenchmarkCase) Name() string {
	if c.CaseName != "" {
		return c.CaseName
	}
	return getName(c.Bench)
}

// SetName sets the case's name, overriding the symbol name if
// needed, and is part of the BenchmarkCase's fluent interface.
func (c *BenchmarkCase) SetName(n string) *BenchmarkCase { c.CaseName = n; return c }

// SetRecorder overrides, the default event recorder type, which
// allows you to change the way that intrarun data is collected and
// allows you to use histogram data if needed for longer runs, and is
// part of the BenchmarkCase's fluent interface.
func (c *BenchmarkCase) SetRecorder(r RecorderType) *BenchmarkCase { c.Recorder = r; return c }

// SetBench allows you set the benchmark cunftion, and is part of the
// BenchmarkCase's fluent interface.
func (c *BenchmarkCase) SetBench(b Benchmark) *BenchmarkCase { c.Bench = b; return c }

// SetCount allows you to set the count number passed to the benchmark
// function which should control the number of internal iterations,
// and is part of the BenchmarkCase's fluent interface.
//
// If running as a standard library test, this value is ignored.
func (c *BenchmarkCase) SetCount(v int) *BenchmarkCase { c.Count = v; return c }

// SetMaxDuration allows to specify a maximum duration for the
// test. If the test has been running for more than this period of
// time, then it will stop running. If you do not specify a timeout
func (c *BenchmarkCase) SetMaxDuration(dur time.Duration) *BenchmarkCase { c.MaxRuntime = dur; return c }
func (c *BenchmarkCase) SetMaxIterations(v int) *BenchmarkCase           { c.MaxIterations = v; return c }
func (c *BenchmarkCase) SetIterationTimeout(dur time.Duration) *BenchmarkCase {
	c.IterationTimeout = dur
	return c
}
func (c *BenchmarkCase) SetTimeout(dur time.Duration) *BenchmarkCase { c.Timeout = dur; return c }

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

	if c.Timeout == 0 {
		c.Timeout = 3 * c.MaxRuntime
	}

	if c.Timeout == 0 {
		c.Timeout = 10 * time.Minute
	}

	if c.IterationTimeout == 0 {
		c.IterationTimeout = 2 * c.MaxRuntime
	}

	if c.IterationTimeout == 0 {
		c.IterationTimeout = 5 * time.Minute
	}

	catcher := grip.NewBasicCatcher()

	if c.Bench == nil {
		catcher.Add(errors.New("must specify a valid benchmark"))
	}

	if c.MinRuntime >= c.MaxRuntime {
		catcher.Add(errors.New("min runtime must not be >= max runtime "))
	}

	if c.MinIterations >= c.MaxIterations {
		catcher.Add(errors.New("min iterations must not be >= max iterations"))
	}

	if c.IterationTimeout > c.Timeout {
		catcher.Add(errors.New("iteration timeout cannot be longer than case timeout"))
	}

	return catcher.Resolve()
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
			bctx, bcancel := context.WithTimeout(ctx)
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
