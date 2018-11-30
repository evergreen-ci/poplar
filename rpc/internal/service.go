package internal

import (
	"context"

	"github.com/evergreen-ci/poplar"
	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

func AttachService(registry *poplar.RecorderRegistry, s *grpc.Server) error {
	srv := &poplarService{
		registry: registry,
	}

	RegisterPoplarMetricsRecorderServer(s, srv)

	return nil
}

type poplarService struct {
	registry *poplar.RecorderRegistry
}

func (s *poplarService) CreateRecorder(ctx context.Context, info *CreateOptions) (*EventResponse, error) {
	_, err := s.registry.Create(info.Name, info.Export())
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return &EventResponse{Name: info.Name, Status: true}, nil
}

func (s *poplarService) CloseRecorder(ctx context.Context, id *RecorderID) (*EventResponse, error) {
	if err := s.registry.Close(id.Name); err != nil {
		return nil, errors.WithStack(err)
	}

	return &EventResponse{Name: id.Name, Status: true}, nil
}

func (s *poplarService) BeginEvent(ctx context.Context, id *RecorderID) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(id.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", id.Name)
	}

	rec.Begin()

	return &EventResponse{Name: id.Name, Status: true}, nil
}

func (s *poplarService) ResetEvent(ctx context.Context, id *RecorderID) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(id.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", id.Name)
	}

	rec.Reset()

	return &EventResponse{Name: id.Name, Status: true}, nil
}

func (s *poplarService) EndEvent(ctx context.Context, val *EventSendDuration) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	dur, err := ptypes.Duration(val.Duration)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert duration value")
	}

	rec.End(dur)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) SetTime(ctx context.Context, t *EventSendTime) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(t.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", t.Name)
	}

	ts, err := ptypes.Timestamp(t.Time)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert timestamp value")
	}

	rec.SetTime(ts)

	return &EventResponse{Name: t.Name, Status: true}, nil
}

func (s *poplarService) SetDuration(ctx context.Context, val *EventSendDuration) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	dur, err := ptypes.Duration(val.Duration)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert duration value")
	}

	rec.SetDuration(dur)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) SetTotalDuration(ctx context.Context, val *EventSendDuration) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	dur, err := ptypes.Duration(val.Duration)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert duration value")
	}

	rec.SetTotalDuration(dur)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) SetState(ctx context.Context, val *EventSendInt) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	rec.SetState(val.Value)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) SetWorkers(ctx context.Context, val *EventSendInt) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	rec.SetWorkers(val.Value)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) SetFailed(ctx context.Context, val *EventSendBool) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	rec.SetFailed(val.Value)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) IncOps(ctx context.Context, val *EventSendInt) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	rec.IncOps(val.Value)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) IncSize(ctx context.Context, val *EventSendInt) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	rec.IncSize(val.Value)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) IncError(ctx context.Context, val *EventSendInt) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	rec.IncError(val.Value)

	return &EventResponse{Name: val.Name, Status: true}, nil
}

func (s *poplarService) IncIterations(ctx context.Context, val *EventSendInt) (*EventResponse, error) {
	rec, ok := s.registry.GetRecorder(val.Name)
	if !ok {
		return nil, errors.Errorf("could not find recorder '%s'", val.Name)
	}

	rec.IncIterations(val.Value)

	return &EventResponse{Name: val.Name, Status: true}, nil
}
