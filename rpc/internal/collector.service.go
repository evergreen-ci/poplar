package internal

import (
	"container/heap"
	"context"
	"io"
	"sync"

	"github.com/evergreen-ci/poplar"
	"github.com/evergreen-ci/utility"
	"github.com/mongodb/ftdc/events"
	"github.com/mongodb/grip"
	"github.com/mongodb/grip/message"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var globalStreamsCoordinator = &streamsCoordinator{}

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
		group     *streamGroup
		streamID  string
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
		if group == nil {
			if event.Name == "" {
				return status.Error(codes.InvalidArgument, "registries must be named")
			}

			eventName = event.Name
			collector, ok := s.registry.GetEventsCollector(eventName)
			if !ok {
				return status.Errorf(codes.NotFound, "no registry named %s", eventName)
			}

			group, streamID = globalStreamsCoordinator.addStream(eventName, collector)
			defer group.removeStream(streamID)
		}

		if event.Name != eventName {
			return status.Errorf(codes.InvalidArgument, "cannot request different registries in the same stream")
		}

		if err := group.addEvent(streamID, event.Export()); err != nil {
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
	streams   map[string]chan error
	eventHeap *PerformanceHeap
	mu        sync.Mutex
}

func (sc *streamsCoordinator) addStream(name string, collector events.Collector) (*streamGroup, string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	group, ok := sc.groups[name]
	if !ok {
		group := &streamGroup{
			collector: collector,
			eventHeap: &PerformanceHeap{},
		}
		heap.Init(group.eventHeap)
		sc.groups[name] = group
	}

	group.mu.Lock()
	defer group.mu.Unlock()

	id := utility.RandomString()
	group.streams[id] = make(chan error)

	return group, id
}

func (sg *streamGroup) addEvent(id string, event *events.Performance) error {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	errChan, ok := sg.streams[id]
	if !ok {
		return errors.Errorf("stream %s does not exist in this stream group", id)
	}

	if sg.eventHeap.Len() == len(sg.streams) {
		sg.pop()
	}
	sg.eventHeap.SafePush(&performanceHeapItem{errChan: errChan, event: event})

	sg.mu.Unlock()
	err := <-errChan
	sg.mu.Lock()

	return err
}

func (sg *streamGroup) pop() {
	item := sg.eventHeap.SafePop()
	if item != nil {
		item.errChan <- sg.collector.AddEvent(item.event)
	}
}

func (sg *streamGroup) removeStream(id string) {
	sg.mu.Lock()
	defer sg.mu.Unlock()

	delete(sg.streams, id)
}

// PerformanceHeap is a min heap of ftdc/events.Performance objects.
type PerformanceHeap struct {
	items []*performanceHeapItem
}

type performanceHeapItem struct {
	errChan chan error
	event   *events.Performance
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
	item, ok := x.(*performanceHeapItem)
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
func (h *PerformanceHeap) SafePush(item *performanceHeapItem) {
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
	return item
}
