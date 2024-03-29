// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.3
// source: collector.proto

package internal

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PoplarEventCollectorClient is the client API for PoplarEventCollector service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PoplarEventCollectorClient interface {
	CreateCollector(ctx context.Context, in *CreateOptions, opts ...grpc.CallOption) (*PoplarResponse, error)
	SendEvent(ctx context.Context, in *EventMetrics, opts ...grpc.CallOption) (*PoplarResponse, error)
	RegisterStream(ctx context.Context, in *CollectorName, opts ...grpc.CallOption) (*PoplarResponse, error)
	StreamEvents(ctx context.Context, opts ...grpc.CallOption) (PoplarEventCollector_StreamEventsClient, error)
	CloseCollector(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error)
}

type poplarEventCollectorClient struct {
	cc grpc.ClientConnInterface
}

func NewPoplarEventCollectorClient(cc grpc.ClientConnInterface) PoplarEventCollectorClient {
	return &poplarEventCollectorClient{cc}
}

func (c *poplarEventCollectorClient) CreateCollector(ctx context.Context, in *CreateOptions, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarEventCollector/CreateCollector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarEventCollectorClient) SendEvent(ctx context.Context, in *EventMetrics, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarEventCollector/SendEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarEventCollectorClient) RegisterStream(ctx context.Context, in *CollectorName, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarEventCollector/RegisterStream", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarEventCollectorClient) StreamEvents(ctx context.Context, opts ...grpc.CallOption) (PoplarEventCollector_StreamEventsClient, error) {
	stream, err := c.cc.NewStream(ctx, &PoplarEventCollector_ServiceDesc.Streams[0], "/poplar.PoplarEventCollector/StreamEvents", opts...)
	if err != nil {
		return nil, err
	}
	x := &poplarEventCollectorStreamEventsClient{stream}
	return x, nil
}

type PoplarEventCollector_StreamEventsClient interface {
	Send(*EventMetrics) error
	CloseAndRecv() (*PoplarResponse, error)
	grpc.ClientStream
}

type poplarEventCollectorStreamEventsClient struct {
	grpc.ClientStream
}

func (x *poplarEventCollectorStreamEventsClient) Send(m *EventMetrics) error {
	return x.ClientStream.SendMsg(m)
}

func (x *poplarEventCollectorStreamEventsClient) CloseAndRecv() (*PoplarResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(PoplarResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *poplarEventCollectorClient) CloseCollector(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarEventCollector/CloseCollector", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PoplarEventCollectorServer is the server API for PoplarEventCollector service.
// All implementations must embed UnimplementedPoplarEventCollectorServer
// for forward compatibility
type PoplarEventCollectorServer interface {
	CreateCollector(context.Context, *CreateOptions) (*PoplarResponse, error)
	SendEvent(context.Context, *EventMetrics) (*PoplarResponse, error)
	RegisterStream(context.Context, *CollectorName) (*PoplarResponse, error)
	StreamEvents(PoplarEventCollector_StreamEventsServer) error
	CloseCollector(context.Context, *PoplarID) (*PoplarResponse, error)
	mustEmbedUnimplementedPoplarEventCollectorServer()
}

// UnimplementedPoplarEventCollectorServer must be embedded to have forward compatible implementations.
type UnimplementedPoplarEventCollectorServer struct {
}

func (UnimplementedPoplarEventCollectorServer) CreateCollector(context.Context, *CreateOptions) (*PoplarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCollector not implemented")
}
func (UnimplementedPoplarEventCollectorServer) SendEvent(context.Context, *EventMetrics) (*PoplarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendEvent not implemented")
}
func (UnimplementedPoplarEventCollectorServer) RegisterStream(context.Context, *CollectorName) (*PoplarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterStream not implemented")
}
func (UnimplementedPoplarEventCollectorServer) StreamEvents(PoplarEventCollector_StreamEventsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamEvents not implemented")
}
func (UnimplementedPoplarEventCollectorServer) CloseCollector(context.Context, *PoplarID) (*PoplarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CloseCollector not implemented")
}
func (UnimplementedPoplarEventCollectorServer) mustEmbedUnimplementedPoplarEventCollectorServer() {}

// UnsafePoplarEventCollectorServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PoplarEventCollectorServer will
// result in compilation errors.
type UnsafePoplarEventCollectorServer interface {
	mustEmbedUnimplementedPoplarEventCollectorServer()
}

func RegisterPoplarEventCollectorServer(s grpc.ServiceRegistrar, srv PoplarEventCollectorServer) {
	s.RegisterService(&PoplarEventCollector_ServiceDesc, srv)
}

func _PoplarEventCollector_CreateCollector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarEventCollectorServer).CreateCollector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarEventCollector/CreateCollector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarEventCollectorServer).CreateCollector(ctx, req.(*CreateOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarEventCollector_SendEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventMetrics)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarEventCollectorServer).SendEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarEventCollector/SendEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarEventCollectorServer).SendEvent(ctx, req.(*EventMetrics))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarEventCollector_RegisterStream_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CollectorName)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarEventCollectorServer).RegisterStream(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarEventCollector/RegisterStream",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarEventCollectorServer).RegisterStream(ctx, req.(*CollectorName))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarEventCollector_StreamEvents_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(PoplarEventCollectorServer).StreamEvents(&poplarEventCollectorStreamEventsServer{stream})
}

type PoplarEventCollector_StreamEventsServer interface {
	SendAndClose(*PoplarResponse) error
	Recv() (*EventMetrics, error)
	grpc.ServerStream
}

type poplarEventCollectorStreamEventsServer struct {
	grpc.ServerStream
}

func (x *poplarEventCollectorStreamEventsServer) SendAndClose(m *PoplarResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *poplarEventCollectorStreamEventsServer) Recv() (*EventMetrics, error) {
	m := new(EventMetrics)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _PoplarEventCollector_CloseCollector_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoplarID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarEventCollectorServer).CloseCollector(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarEventCollector/CloseCollector",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarEventCollectorServer).CloseCollector(ctx, req.(*PoplarID))
	}
	return interceptor(ctx, in, info, handler)
}

// PoplarEventCollector_ServiceDesc is the grpc.ServiceDesc for PoplarEventCollector service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PoplarEventCollector_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "poplar.PoplarEventCollector",
	HandlerType: (*PoplarEventCollectorServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCollector",
			Handler:    _PoplarEventCollector_CreateCollector_Handler,
		},
		{
			MethodName: "SendEvent",
			Handler:    _PoplarEventCollector_SendEvent_Handler,
		},
		{
			MethodName: "RegisterStream",
			Handler:    _PoplarEventCollector_RegisterStream_Handler,
		},
		{
			MethodName: "CloseCollector",
			Handler:    _PoplarEventCollector_CloseCollector_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamEvents",
			Handler:       _PoplarEventCollector_StreamEvents_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "collector.proto",
}
