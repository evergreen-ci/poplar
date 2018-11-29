package internal

import (
	"context"

	"github.com/evergreen-ci/poplar"
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

func (s *poplarService) CreateRecorder(context.Context, *RecorderID) (*RecorderInfo, error) {
	return nil, nil
}

func (s *poplarService) CloseRecorder(context.Context, *RecorderID) (*RecorderInfo, error) {
	return nil, nil
}

func (s *poplarService) BeginEvent(context.Context, *RecorderID) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) ResetRecorder(context.Context, *RecorderID) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) End(context.Context, *EventSendDuration) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) SetTime(context.Context, *EventSendTime) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) SetDuration(context.Context, *EventSendDuration) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) SetTotalDuration(context.Context, *EventSendDuration) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) SetState(context.Context, *EventSendInt) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) SetWorkers(context.Context, *EventSendInt) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) SetBool(context.Context, *EventSendBool) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) IncOps(context.Context, *EventSendInt) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) IncSize(context.Context, *EventSendInt) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) IncError(context.Context, *EventSendInt) (*EventResponse, error) {
	return nil, nil
}

func (s *poplarService) IncIterations(context.Context, *EventSendInt) (*EventResponse, error) {
	return nil, nil
}
