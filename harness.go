package poplar

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/mongodb/ftdc/events"
	"github.com/mongodb/grip"
	"github.com/pkg/errors"
)

type BenchmarkCase struct {
	Bench         Benchmark
	MinRuntime    time.Duration
	MaxRuntime    time.Duration
	Count         int
	MinIterations int
	MaxIterations int
	Recorder      RecorderType
}

type BenchmarkResult struct {
	Name         string
	Runtime      time.Duration
	Count        int
	Iterations   int
	Error        error
	ArtifactPath string
	StartAt      time.Time
	CompletedAt  time.Time
}

type BenchmarkSuite []BenchmarkCase
type BenchmarkSuiteResults []BenchmarkResult

func (res BenchmarkSuiteResults) String() string {
	out := make([]string, len(res))

	for idx := range res {
		out[idx] = res[idx].Report()
	}

	return strings.Join(out, "\n")
}

func RunBenchmarks(ctx context.Context, prefix string, cases BenchmarkSuite) (BenchmarkSuiteResults, error) {
	registry := NewRegistry()
	catcher := grip.NewBasicCatcher()
	res := make(BenchmarkSuiteResults, 0, len(cases))

	for _, test := range cases {
		name := test.Name()
		path := filepath.Join(prefix, name+".ftdc")
		recorder, err := registry.Create(name, CreateOptions{
			Path:      path,
			ChunkSize: 1024,
			Streaming: true,
			Dynamic:   true,
			Recorder:  test.Recorder,
		})
		if err != nil {
			catcher.Add(err)
			break
		}

		startAt := time.Now()
		out := test.Run(ctx, recorder)
		out.CompletedAt = time.Now()
		out.StartAt = startAt
		out.ArtifactPath = path

		catcher.Add(out.Error)
		res = append(res, out)
		catcher.Add(registry.Close(name))
	}

	return res, catcher.Resolve()
}

type Benchmark func(context.Context, events.Recorder, int) error

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

func (c *BenchmarkCase) Name() string { return getName(c.Bench) }

func (c *BenchmarkCase) String() string {
	return fmt.Sprintf("name=%s, count=%d, min_dur=%s, max_dur=%s, min_iters=%s, max_iters=%d",
		c.Name(), c.Count, c.MinRuntime, c.MaxRuntime, c.MinIterations, c.MaxIterations)
}

func (c *BenchmarkCase) Validate() error {
	if c.Recorder == "" {
		c.Recorder = RecorderPerf
	}

	if c.MaxRuntime == 0 {
		c.MaxRuntime = 5 * time.Minute
	}

	if c.Count == 0 {
		c.Count = 1
	}

	return nil
}

func (c *BenchmarkCase) Run(ctx context.Context, recorder events.Recorder) BenchmarkResult {
	res := BenchmarkResult{
		Iterations: 1,
		Count:      c.Count,
	}

	var cancel context.CancelFunc
	ctx, cancel = context.WithTimeout(ctx, 2*c.MaxRuntime)
	defer cancel()

	for {
		if ctx.Err() != nil {
			break
		}

		if c.MaxRuntime > 0 && res.Runtime >= c.MaxRuntime {
			if c.MinIterations == 0 || res.Iterations > c.MinIterations {
				break
			}
		}

		if c.MaxIterations > 0 && res.Iterations >= c.MaxIterations {
			if c.MinRuntime == 0 || res.Runtime > c.MinRuntime {
				break
			}
		}

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

	return res
}

func (res *BenchmarkResult) Report() string {
	out := []string{
		"=== RUN", res.Name,
		"    --- REPORT: " + fmt.Sprintf("count=%d, iters=%s, runtime=%s", res.Count, res.Iterations, roundDurationMS(res.Runtime)),
	}

	if res.Error != nil {
		out = append(out,
			fmt.Sprintf("    --- ERRORS: %s", res.Error.Error()),
			fmt.Sprintf("--- FAIL: %s (%s)", res.Name, roundDurationMS(res.Runtime)))
	} else {
		out = append(out, fmt.Sprintf("--- PASS: %s", res.Name, roundDurationMS(res.Runtime)))
	}

	return strings.Join(out, "\n")

}
