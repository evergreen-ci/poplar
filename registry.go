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
}

type RecorderRegistry struct {
	cache map[string]*recorderInstance
	mu    sync.Mutex
}

func (r *RecorderRegistry) Create(key string, flavor RecorderType, collOpts CollectorCreationOptions) (events.Recorder, error) {
	var out events.Recorder

	if err := flavor.Validate(); err != nil {
		return nil, errors.Wrap(err, "could not build recorder")
	}

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

func (r *RecorderRegistry) Get(key string) (events.Recorder, bool) {
	r.mu.Lock()
	defer r.mu.Unlock()

	impl, ok := r.mu[key]

	return impl.recorder, ok
}

func (r *RecorderRegistry) Drop(key string) error {
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

type CollectorCreationOptions struct {
	Path      string
	ChunkSize int
	Streaming bool
	Dynamic   bool
	Recorder  RecorderType
}

func (opts *CollectorCreationOptions) build() (*recorderInstance, error) {
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

	var collector ftdc.Collector
	switch {
	case opts.Streaming && opts.Dynamic:
		collector = ftdc.NewStreamingDynamicCollector(opts.ChunkSize, file)
	case !opts.Streaming && opts.Dynamic:
		collector = ftdc.NewDynamicCollector(opts.ChunkSize)
	case opts.Streaming && !opts.Dynamic:
		collector = ftdc.NewStreamingCollector(opts.ChunkSize, file)
	case !opts.Streaming && !opts.Dynamic:
		collector = ftdc.NewBatchCollector(opts.ChunkSize)
	default:
		return nil, errors.New("invalid collector defined")
	}

	var recorder events.Recorder
	switch {
	case RecorderPerf:
		recorder = events.NewRawRecorder(collector)
	case RecorderPerfSingle:
		recorder = events.NewSingleRecorder(collector)
	case RecorderPerf100ms:
		recorder = events.NewGroupedRecorder(collector, 100*time.Millisecond)
	case RecorderPerf1s:
		recorder = events.NewGroupedRecorder(collector, time.Second)
	case RecorderHistogramSingle:
		recorder = events.NewSingleHistogramRecorder(collector)
	case RecorderHistogram100ms:
		recorder = events.NewHistogramGroupedRecorder(collector, 100*time.Millisecond)
	case RecorderHistogram1s:
		recorder = events.NewHistogramGroupedRecorder(collector, time.Second)
	default:
		return nil, errors.New("invalid recorder defined")
	}

	return &recorderInstance{
		file:      file,
		collector: collector,
		recorder:  recorder,
	}, nil
}
