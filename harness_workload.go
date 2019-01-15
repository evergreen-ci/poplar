package poplar

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/mongodb/grip/recovery"
	"github.com/pkg/errors"
)

// BenchmarkWorkload provides a way to express a more complex
// performance test, that involves multiple instances of a benchmark
// test running at the same time.
//
// You can specify the workload as either a single benchmark case, or
// as a ordered list of benchmark operations, however it is not valid
// to do both in the same workload instance.
//
// If you specify a group of workload operations when executing,
// poplar will run each sub-workload (with however many instances of
// the workload are specified,) sequentially, with no inter-workload
// synchronization.
type BenchmarkWorkload struct {
	Name      string
	Timeout   *time.Duration
	Case      *BenchmarkCase
	Group     []BenchmarkWorkload
	Instances int
	Recorder  RecorderType
}

// Validate ensures that the workload is well formed. Additionally,
// ensures that the all cases and workload groups are valid, and have
// the same recorder type defined.
func (w *BenchmarkWorkload) Validate() error {
	catcher := grip.NewBasicCatcher()

	if w.Timeout != nil && *w.Timeout < time.Millisecond {
		catcher.Add(errors.New("cannot specify timeout less than a millisecond"))
	}

	if w.Case == nil && w.Group == nil {
		catcher.Add(errors.New("cannot define a workload with out work"))
	}

	if w.Case != nil && w.Group != nil {
		catcher.Add(errors.New("cannot define a workload with both a case and a group"))
	}

	if w.Instances <= 1 {
		catcher.Add(errors.New("must define more than a single instance in a workload"))
	}

	if w.Name == "" {
		catcher.Add(errors.New("must specify a name for a workload"))
	}

	if w.Case != nil {
		catcher.Add(w.Case.Validate())
		w.Case.Recorder = w.Recorder
	}

	for _, gwl := range w.Group {
		catcher.Add(gwl.Validate())
		gwl.Recorder = w.Recorder
	}

	return catcher.Resolve()
}

// Run executes the workload, and has similar semantics to the
// BenchmarkSuite implementation.
func (w *BenchmarkWorkload) Run(ctx context.Context, prefix string) (BenchmarkResultGroup, error) {
	var cancel context.CancelFunc
	ctx, cancel = context.WithCancel(ctx)
	defer cancel()

	if w.Timeout != nil {
		ctx, cancel = context.WithTimeout(ctx, *w.Timeout)
		defer cancel()
	}

	registry := NewRegistry()
	catcher := grip.NewBasicCatcher()
	wg := &sync.WaitGroup{}
	results := make(chan BenchmarkResult)
	res := BenchmarkResultGroup{}

	go func() {
		defer func() {
			catcher.Add(recovery.AnnotateMessageWithPanicError(recover(), nil, message.Fields{
				"name":     w.Name,
				"op":       "collecting results",
				"executor": "native",
			}))
		}()

		for {

			select {
			case r := <-results:
				res = append(res, r)
			case <-ctx.Done():
				return
			}
		}
	}()

	for i := 0; i < w.Instances; i++ {
		wg.Add(1)
		go func(instanceIdx int) {
			defer wg.Done()
			defer func() {
				catcher.Add(recovery.AnnotateMessageWithPanicError(recover(), nil, message.Fields{
					"idx":      instanceIdx,
					"name":     w.Name,
					"executor": "native",
					"op":       "running workload",
				}))
			}()

			name := fmt.Sprintf("%s.%d", w.Name, instanceIdx)
			path := filepath.Join(prefix, name+".ftdc")
			recorder, err := registry.Create(name, CreateOptions{
				Path:      path,
				ChunkSize: 1024,
				Streaming: true,
				Dynamic:   true,
				Recorder:  w.Recorder,
			})
			if err != nil {
				catcher.Add(err)
				return
			}

			if w.Case != nil {
				select {
				case results <- w.Case.Run(ctx, recorder):
					return
				case <-ctx.Done():
					return
				}
			}

			for idx, wlg := range w.Group {
				resultGroup, err := wlg.Run(ctx, fmt.Sprintf("%s.%s.%d.%d", prefix, name, idx))
				catcher.Add(err)

				for _, r := range resultGroup {
					select {
					case results <- r:
						continue
					case <-ctx.Done():
						return
					}
				}
			}
		}(i)
	}
	wg.Wait()
	close(results)

	for r := range results {
		r.Workload = true
		if r.Instances == 0 {
			r.Instances = w.Instances
		}
	}

	return res, catcher.Resolve()
}

// Standard produces a standard golang benchmarking function from a
// poplar workload.
//
// These invocations are not able to respect the top-level workload
// timeout, and *do* perform pre-flight workload validation.
func (w *BenchmarkWorkload) Standard(registry *RecorderRegistry) func(*testing.B) {
	return func(b *testing.B) {
		if err := w.Validate(); err != nil {
			b.Fatal(errors.Wrap(err, "benchmark workload failed"))
		}

		wg := &sync.WaitGroup{}
		for i := 0; i < w.Instances; i++ {
			wg.Add(1)
			go func(id int) {
				defer func() {
					err := recovery.AnnotateMessageWithPanicError(recover(), nil, message.Fields{
						"idx":      id,
						"name":     w.Name,
						"op":       "running workload",
						"executor": "standard",
					})

					if err != nil {
						b.Fatal(err)
					}
				}()

				if w.Case != nil {
					b.Run(fmt.Sprintf("WorkloadCase%s%s#%d", w.Name, w.Case.Name(), id), w.Case.Standard(registry))
					return
				}

				for idx, wlg := range w.Group {
					b.Run(fmt.Sprintf("WorkloadGroup%s%s%d#%d", w.Name, wlg.Name, idx, id), wlg.Standard(registry))
				}
			}(i)
		}
		wg.Done()
	}
}
