package poplar

import (
	"io"
	"os"
	"sync"
	"time"

	"github.com/mongodb/ftdc"
	"github.com/mongodb/ftdc/events"
	"github.com/pkg/errors"
)

type RecorderType string

const (
	RecorderPerf            RecorderType = "perf"
	RecorderPerfSingle                   = "perf-single"
	RecorderPerf100ms                    = "perf-grouped-100ms"
	RecorderPerf1s                       = "perf-grouped-1s"
	RecorderHistogramSingle              = "histogram-single"
	RecorderHistogram100ms               = "histogram-grouped-100ms"
	RecorderHistogram1s                  = "histogram-grouped-1s"
)

func (t RecorderType) Validate() error {
	switch t {
	case RecorderPerf, RecorderPerfSingle, RecorderPerf100ms, RecorderPerf1s,
		RecorderHistogramSingle, RecorderHistogram100ms, RecorderHistogram1s:

		return nil
	default:
		return errors.Errorf("%s is not a supported recorder type", t)
	}
}

type recorderInstance struct {
	file      io.WriteCloser
	collector ftdc.Collector
	recorder  events.Recorder
	isDynamic bool
}

type RecorderRegistry struct {
	cache map[string]*recorderInstance
	mu    sync.Mutex
}

func NewRegistry() *RecorderRegistry {
	return &RecorderRegistry{
		cache: map[string]*recorderInstance{},
	}
}

// Create builds a new collector, of the given name with the specified
// options controling the collector type and configuration.
//
// If the options specify a filename that already exists, then Create
// will return an error.
func (r *RecorderRegistry) Create(key string, collOpts CreateOptions) (events.Recorder, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	_, ok := r.cache[key]
	if ok {
		return nil, errors.Errorf("a recorder named '%s' already exists", key)
	}

	instance, err := collOpts.build()
	if err != nil {
		return nil, errors.Wrap(err, "could not construct recorder output")
	}

	r.cache[key] = instance

	return instance.recorder, nil
}

// GetRecorder returns the Recorder instance for this key. Returns
// false when the recorder does not exist.
func (r *RecorderRegistry) GetRecorder(key string) (events.Recorder, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	impl, ok := r.cache[key]
	if !ok {
		return nil, false
	}

	return impl.recorder, true
}

// GetCollector returns the collector instance for this key. Will
// return false, when the collector does not exist OR if the collector
// is dynamic.
func (r *RecorderRegistry) GetCollector(key string) (ftdc.Collector, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	impl, ok := r.cache[key]

	if !ok || !impl.isDynamic {
		return nil, false
	}

	return impl.collector, true
}

// Close flushes and closes the underlying recorder and collector and
// then removes it from the cache.
func (r *RecorderRegistry) Close(key string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if impl, ok := r.cache[key]; ok {
		if err := impl.recorder.Flush(); err != nil {
			return errors.Wrap(err, "problem flushing recorder")
		}
		if err := ftdc.FlushCollector(impl.collector, impl.file); err != nil {
			return errors.Wrap(err, "problem writing collector contents to file")
		}
		if err := impl.file.Close(); err != nil {
			return errors.Wrap(err, "problem closing open file")
		}
	}

	delete(r.cache, key)
	return nil
}

type CreateOptions struct {
	Path      string
	ChunkSize int
	Streaming bool
	Dynamic   bool
	Recorder  RecorderType
}

func (opts *CreateOptions) build() (*recorderInstance, error) {
	if err := opts.Recorder.Validate(); err != nil {
		return nil, errors.Wrap(err, "invalid recorder type")
	}

	if _, err := os.Stat(opts.Path); !os.IsNotExist(err) {
		return nil, errors.Errorf("could not create '%s' because it exists", opts.Path)
	}

	file, err := os.Create(opts.Path)
	if err != nil {
		return nil, errors.Wrapf(err, "problem opening file '%s'", opts.Path)
	}

	out := &recorderInstance{
		isDynamic: opts.Dynamic,
	}

	switch {
	case opts.Streaming && opts.Dynamic:
		out.collector = ftdc.NewStreamingDynamicCollector(opts.ChunkSize, file)
	case !opts.Streaming && opts.Dynamic:
		out.collector = ftdc.NewDynamicCollector(opts.ChunkSize)
	case opts.Streaming && !opts.Dynamic:
		out.collector = ftdc.NewStreamingCollector(opts.ChunkSize, file)
	case !opts.Streaming && !opts.Dynamic:
		out.collector = ftdc.NewBatchCollector(opts.ChunkSize)
	default:
		return nil, errors.New("invalid collector defined")
	}

	switch opts.Recorder {
	case RecorderPerf:
		out.recorder = events.NewRawRecorder(out.collector)
	case RecorderPerfSingle:
		out.recorder = events.NewSingleRecorder(out.collector)
	case RecorderPerf100ms:
		out.recorder = events.NewGroupedRecorder(out.collector, 100*time.Millisecond)
	case RecorderPerf1s:
		out.recorder = events.NewGroupedRecorder(out.collector, time.Second)
	case RecorderHistogramSingle:
		out.recorder = events.NewSingleHistogramRecorder(out.collector)
	case RecorderHistogram100ms:
		out.recorder = events.NewHistogramGroupedRecorder(out.collector, 100*time.Millisecond)
	case RecorderHistogram1s:
		out.recorder = events.NewHistogramGroupedRecorder(out.collector, time.Second)
	default:
		return nil, errors.New("invalid recorder defined")
	}

	return out, nil
}
