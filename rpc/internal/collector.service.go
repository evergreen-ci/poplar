package internal

import (
	"container/heap"
	"context"
	"io"
	"sync"
	"time"

	"github.com/evergreen-ci/poplar"
	"github.com/evergreen-ci/utility"
	"github.com/mongodb/ftdc/events"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var gloablStreamsCoorindator = &streamsCoordinator{}

type collectorService struct {
	registry *poplar.RecorderRegistry
}

func (s *collectorService) CreateCollector(ctx context.Context, opts *CreateOptions) (*PoplarResponse, error) {
	if _, ok := s.registry.GetCollector(opts.Name); !ok {
		_, err := s.registry.Create(opts.Name, opts.Export())
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	return &PoplarResponse{Name: opts.Name, Status: true}, nil
}

func (s *collectorService) CloseCollector(ctx context.Context, id *PoplarID) (*PoplarResponse, error) {
	err := s.registry.Close(id.Name)

	grip.Error(message.WrapError(err, message.Fields{
		"message":  "problem closing recorder",
		"recorder": id.Name,
	}))

	return &PoplarResponse{Name: id.Name, Status: err == nil}, nil

}

func (s *collectorService) SendEvent(ctx context.Context, event *EventMetrics) (*PoplarResponse, error) {
	collector, ok := s.registry.GetEventsCollector(event.Name)

	if !ok {
		return nil, status.Errorf(codes.NotFound, "no registry named %s", event.Name)
	}

	err := collector.AddEvent(event.Export())

	return &PoplarResponse{Name: event.Name, Status: err == nil}, nil

}

func (s *collectorService) StreamEvents(srv PoplarEventCollector_StreamEventsServer) error {
	ctx := srv.Context()

	var (
		collector events.Collector
		eventName string
	)

	for {
		event, err := srv.Recv()
		if err == io.EOF {
			return srv.SendAndClose(&PoplarResponse{
				Name:   eventName,
				Status: true,
			})
		} else if err != nil {
			return srv.SendAndClose(&PoplarResponse{
				Name:   eventName,
				Status: false,
			})
		}
		if collector == nil {
			if event.Name == "" {
				return status.Error(codes.InvalidArgument, "registries must be named")
			}

			eventName = event.Name
			var ok bool
			collector, ok = s.registry.GetEventsCollector(eventName)
			if !ok {
				return status.Errorf(codes.NotFound, "no registry named %s", eventName)
			}
		}

		if err := collector.AddEvent(event.Export()); err != nil {
			return status.Errorf(codes.Internal, "problem persisting argument %s", err.Error())
		}

		if ctx.Err() != nil {
			return status.Errorf(codes.Canceled, "operation canceled for %s", eventName)
		}
	}
}

type streamsCoordinator struct {
	ctx    context.Context
	groups map[string]*streamGroup
	mu     sync.Mutex
}

type streamGroup struct {
	collector events.Collector
	streams   map[string]grip.Catcher
	eventHeap *PerformanceHeap
	cancel    context.CancelFunc
	closed    bool
	mu        sync.Mutex
}

func (sc *streamsCoordinator) addStream(ctx context.Context, name string, collector events.Collector) (*streamGroup, string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	group, ok := sc.groups[name]
	if !ok || group.closed {
		ctx, cancel := context.WithCancel(ctx)
		group := &streamGroup{
			collector: collector,
			eventHeap: &PerformanceHeap{},
			cancel:    cancel,
		}
		heap.Init(group.eventHeap)
		go group.timedFlush(5 * time.Second)
		sc.groups[name] = group
	}

	group.mu.Lock()
	defer group.mu.Unlock()

	id := utility.RandomString()
	sc.streams[id] = grip.NewBasicCatcher()

	return group, id
}

func (sg *streamGroup) addEvent(id string, event *events.Performance) error {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	catcher, ok := sg.streams[id]
	if !ok {
		return errors.Errorf("stream %s does not exist in stream group for %s", id, name)
	}

	if catcher.HasErrors() {
		return catcher.Resolve()
	}

	sg.eventHeap.SafePush(event)
}

func (sg *streamGroup) getErrors(id string) error {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	catcher, ok := sg.streams[id]
	if !ok {
		return errors.Errorf("stream %s does not exist in stream group for %s", id, name)
	}
	sg.streams[id] = grip.NewBasicCatcher()

	return catcher.Resolve()
}

func (sg *streamGroup) closeGroup() error {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	sg.cancel()
	sg.flush(time.Hour)

	catcher := grip.NewBasicCatcher()
	for _, c := range sg.streams {
		catcher.Add(c.Resolve())
	}

	sg.closed = true

	return catcher.Resolve()
}

func (sg *streamGroup) timedFlush(ctx context.Context, flushInterval time.Duration) {
	timer := time.NewTimer(b.opts.FlushInterval)
	defer timer.Stop()

	for {
		select {
		case <-ctx.Done():
			sg.mu.Lock()
			sg.flush(time.Hour)
			sg.mu.Unlock()
			return
		case <-timer.C:
			sg.mu.Lock()
			sg.flush(flushInterval)
			sg.mu.Unlock()
		}
	}
}

func (sg *streamGroup) flush(flushInterval time.Duration) {
	if sg.closed {
		return
	}

	for sg.eventHeap.Len() > 0 {
		ts := sg.eventHeap.SafePeek()
		if time.Since(ts) < flushInterval {
			item := sg.eventHeap.SafePop()
			err := sg.collector.AddEvent(item.event)
			sg.streams[id].Add(err)
		} else {
			break
		}
	}
}

// PerformanceHeap is a min heap of ftdc/events.Performance objects.
type PerformanceHeap struct {
	items []performanceHeapItem
}

type performanceHeapItem struct {
	id    string
	event *events.Performance
}

// Len returns the size of the heap.
func (h PerformanceHeap) Len() int { return len(h.items) }

// Less returns true if the object at index i is less than the object at index
// j in the heap, false otherwise.
func (h PerformanceHeap) Less(i, j int) bool {
	return h.items[i].event.Timestamp.Before(h.items[j].event.Timestamp)
}

// Swap swaps the objects at indexes i and j.
func (h PerformanceHeap) Swap(i, j int) { h.items[i], h.items[j] = h.items[j], h.items[i] }

// Push appends a new object of type Performance to the heap. Note that if x is
// not a performanceHeapItem object nothing happens.
func (h *PerformanceHeap) Push(x interface{}) {
	item, ok := x.(performanceHeapItem)
	if !ok {
		return
	}

	h.items = append(h.items, item)
}

// Pop returns the next object (as an empty interface) from the heap. Note that
// if the heap is empty this will panic.
func (h *PerformanceHeap) Pop() interface{} {
	old := h.items
	n := len(old)
	x := old[n-1]
	h.items = old[0 : n-1]
	return x
}

// SafePush is a wrapper function around heap.Push that ensures, during compile
// time, that the correct type of object is put in the heap.
func (h *PerformanceHeap) SafePush(item performanceHeapItem) {
	heap.Push(h, item)
}

// SafePop is a wrapper function around heap.Pop that converts the returned
// interface into a pointer to a  performanceHeapItem object before returning
// it.
func (h *PerformanceHeap) SafePop() *performanceHeapItem {
	if h.Len() == 0 {
		return nil
	}

	i := heap.Pop(h)
	item := i.(*performanceHeapItem)
	return it
}

// SafePeek returns the timestamp of the item at the top of the heap without
// popping it off the heap. If there are not items in the heap, the zero
// time.Time value is returned.
func (h *PerformanceHeap) SafePeek() time.Time {
	if h.Len() == 0 {
		return time.Time{}
	}

	return h.items[0].event.Timestamp
}
