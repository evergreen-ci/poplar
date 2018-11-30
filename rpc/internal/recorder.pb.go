// Code generated by protoc-gen-go. DO NOT EDIT.
// source: recorder.proto

package internal

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import duration "github.com/golang/protobuf/ptypes/duration"
import timestamp "github.com/golang/protobuf/ptypes/timestamp"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type CreateOptions_RecorderType int32

const (
	CreateOptions_UNKNOWN          CreateOptions_RecorderType = 0
	CreateOptions_PERF             CreateOptions_RecorderType = 1
	CreateOptions_PERF_SINGLE      CreateOptions_RecorderType = 2
	CreateOptions_PERF_100MS       CreateOptions_RecorderType = 3
	CreateOptions_PERF_1S          CreateOptions_RecorderType = 4
	CreateOptions_HISTOGRAM_SINGLE CreateOptions_RecorderType = 6
	CreateOptions_HISTOGRAM_100MS  CreateOptions_RecorderType = 7
	CreateOptions_HISTOGRAM_1S     CreateOptions_RecorderType = 8
)

var CreateOptions_RecorderType_name = map[int32]string{
	0: "UNKNOWN",
	1: "PERF",
	2: "PERF_SINGLE",
	3: "PERF_100MS",
	4: "PERF_1S",
	6: "HISTOGRAM_SINGLE",
	7: "HISTOGRAM_100MS",
	8: "HISTOGRAM_1S",
}
var CreateOptions_RecorderType_value = map[string]int32{
	"UNKNOWN":          0,
	"PERF":             1,
	"PERF_SINGLE":      2,
	"PERF_100MS":       3,
	"PERF_1S":          4,
	"HISTOGRAM_SINGLE": 6,
	"HISTOGRAM_100MS":  7,
	"HISTOGRAM_1S":     8,
}

func (x CreateOptions_RecorderType) String() string {
	return proto.EnumName(CreateOptions_RecorderType_name, int32(x))
}
func (CreateOptions_RecorderType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{6, 0}
}

type RecorderID struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecorderID) Reset()         { *m = RecorderID{} }
func (m *RecorderID) String() string { return proto.CompactTextString(m) }
func (*RecorderID) ProtoMessage()    {}
func (*RecorderID) Descriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{0}
}
func (m *RecorderID) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecorderID.Unmarshal(m, b)
}
func (m *RecorderID) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecorderID.Marshal(b, m, deterministic)
}
func (dst *RecorderID) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecorderID.Merge(dst, src)
}
func (m *RecorderID) XXX_Size() int {
	return xxx_messageInfo_RecorderID.Size(m)
}
func (m *RecorderID) XXX_DiscardUnknown() {
	xxx_messageInfo_RecorderID.DiscardUnknown(m)
}

var xxx_messageInfo_RecorderID proto.InternalMessageInfo

func (m *RecorderID) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type EventSendTime struct {
	Name                 string               `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Time                 *timestamp.Timestamp `protobuf:"bytes,2,opt,name=time,proto3" json:"time,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *EventSendTime) Reset()         { *m = EventSendTime{} }
func (m *EventSendTime) String() string { return proto.CompactTextString(m) }
func (*EventSendTime) ProtoMessage()    {}
func (*EventSendTime) Descriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{1}
}
func (m *EventSendTime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendTime.Unmarshal(m, b)
}
func (m *EventSendTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendTime.Marshal(b, m, deterministic)
}
func (dst *EventSendTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendTime.Merge(dst, src)
}
func (m *EventSendTime) XXX_Size() int {
	return xxx_messageInfo_EventSendTime.Size(m)
}
func (m *EventSendTime) XXX_DiscardUnknown() {
	xxx_messageInfo_EventSendTime.DiscardUnknown(m)
}

var xxx_messageInfo_EventSendTime proto.InternalMessageInfo

func (m *EventSendTime) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EventSendTime) GetTime() *timestamp.Timestamp {
	if m != nil {
		return m.Time
	}
	return nil
}

type EventSendInt struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                int64    `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventSendInt) Reset()         { *m = EventSendInt{} }
func (m *EventSendInt) String() string { return proto.CompactTextString(m) }
func (*EventSendInt) ProtoMessage()    {}
func (*EventSendInt) Descriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{2}
}
func (m *EventSendInt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendInt.Unmarshal(m, b)
}
func (m *EventSendInt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendInt.Marshal(b, m, deterministic)
}
func (dst *EventSendInt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendInt.Merge(dst, src)
}
func (m *EventSendInt) XXX_Size() int {
	return xxx_messageInfo_EventSendInt.Size(m)
}
func (m *EventSendInt) XXX_DiscardUnknown() {
	xxx_messageInfo_EventSendInt.DiscardUnknown(m)
}

var xxx_messageInfo_EventSendInt proto.InternalMessageInfo

func (m *EventSendInt) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EventSendInt) GetValue() int64 {
	if m != nil {
		return m.Value
	}
	return 0
}

type EventSendBool struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Value                bool     `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventSendBool) Reset()         { *m = EventSendBool{} }
func (m *EventSendBool) String() string { return proto.CompactTextString(m) }
func (*EventSendBool) ProtoMessage()    {}
func (*EventSendBool) Descriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{3}
}
func (m *EventSendBool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendBool.Unmarshal(m, b)
}
func (m *EventSendBool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendBool.Marshal(b, m, deterministic)
}
func (dst *EventSendBool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendBool.Merge(dst, src)
}
func (m *EventSendBool) XXX_Size() int {
	return xxx_messageInfo_EventSendBool.Size(m)
}
func (m *EventSendBool) XXX_DiscardUnknown() {
	xxx_messageInfo_EventSendBool.DiscardUnknown(m)
}

var xxx_messageInfo_EventSendBool proto.InternalMessageInfo

func (m *EventSendBool) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EventSendBool) GetValue() bool {
	if m != nil {
		return m.Value
	}
	return false
}

type EventSendDuration struct {
	Name                 string             `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Duration             *duration.Duration `protobuf:"bytes,2,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *EventSendDuration) Reset()         { *m = EventSendDuration{} }
func (m *EventSendDuration) String() string { return proto.CompactTextString(m) }
func (*EventSendDuration) ProtoMessage()    {}
func (*EventSendDuration) Descriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{4}
}
func (m *EventSendDuration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendDuration.Unmarshal(m, b)
}
func (m *EventSendDuration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendDuration.Marshal(b, m, deterministic)
}
func (dst *EventSendDuration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendDuration.Merge(dst, src)
}
func (m *EventSendDuration) XXX_Size() int {
	return xxx_messageInfo_EventSendDuration.Size(m)
}
func (m *EventSendDuration) XXX_DiscardUnknown() {
	xxx_messageInfo_EventSendDuration.DiscardUnknown(m)
}

var xxx_messageInfo_EventSendDuration proto.InternalMessageInfo

func (m *EventSendDuration) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EventSendDuration) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

type EventResponse struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Status               bool     `protobuf:"varint,2,opt,name=status,proto3" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EventResponse) Reset()         { *m = EventResponse{} }
func (m *EventResponse) String() string { return proto.CompactTextString(m) }
func (*EventResponse) ProtoMessage()    {}
func (*EventResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{5}
}
func (m *EventResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventResponse.Unmarshal(m, b)
}
func (m *EventResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventResponse.Marshal(b, m, deterministic)
}
func (dst *EventResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventResponse.Merge(dst, src)
}
func (m *EventResponse) XXX_Size() int {
	return xxx_messageInfo_EventResponse.Size(m)
}
func (m *EventResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_EventResponse.DiscardUnknown(m)
}

var xxx_messageInfo_EventResponse proto.InternalMessageInfo

func (m *EventResponse) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *EventResponse) GetStatus() bool {
	if m != nil {
		return m.Status
	}
	return false
}

type CreateOptions struct {
	Name                 string                     `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Path                 string                     `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	ChunkSize            int32                      `protobuf:"varint,3,opt,name=chunkSize,proto3" json:"chunkSize,omitempty"`
	Streaming            bool                       `protobuf:"varint,4,opt,name=streaming,proto3" json:"streaming,omitempty"`
	Dynamic              bool                       `protobuf:"varint,5,opt,name=dynamic,proto3" json:"dynamic,omitempty"`
	Recorder             CreateOptions_RecorderType `protobuf:"varint,6,opt,name=recorder,proto3,enum=poplar.CreateOptions_RecorderType" json:"recorder,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *CreateOptions) Reset()         { *m = CreateOptions{} }
func (m *CreateOptions) String() string { return proto.CompactTextString(m) }
func (*CreateOptions) ProtoMessage()    {}
func (*CreateOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_recorder_a9549d9f8d4fbcd2, []int{6}
}
func (m *CreateOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateOptions.Unmarshal(m, b)
}
func (m *CreateOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateOptions.Marshal(b, m, deterministic)
}
func (dst *CreateOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateOptions.Merge(dst, src)
}
func (m *CreateOptions) XXX_Size() int {
	return xxx_messageInfo_CreateOptions.Size(m)
}
func (m *CreateOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateOptions.DiscardUnknown(m)
}

var xxx_messageInfo_CreateOptions proto.InternalMessageInfo

func (m *CreateOptions) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CreateOptions) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *CreateOptions) GetChunkSize() int32 {
	if m != nil {
		return m.ChunkSize
	}
	return 0
}

func (m *CreateOptions) GetStreaming() bool {
	if m != nil {
		return m.Streaming
	}
	return false
}

func (m *CreateOptions) GetDynamic() bool {
	if m != nil {
		return m.Dynamic
	}
	return false
}

func (m *CreateOptions) GetRecorder() CreateOptions_RecorderType {
	if m != nil {
		return m.Recorder
	}
	return CreateOptions_UNKNOWN
}

func init() {
	proto.RegisterType((*RecorderID)(nil), "poplar.RecorderID")
	proto.RegisterType((*EventSendTime)(nil), "poplar.EventSendTime")
	proto.RegisterType((*EventSendInt)(nil), "poplar.EventSendInt")
	proto.RegisterType((*EventSendBool)(nil), "poplar.EventSendBool")
	proto.RegisterType((*EventSendDuration)(nil), "poplar.EventSendDuration")
	proto.RegisterType((*EventResponse)(nil), "poplar.EventResponse")
	proto.RegisterType((*CreateOptions)(nil), "poplar.CreateOptions")
	proto.RegisterEnum("poplar.CreateOptions_RecorderType", CreateOptions_RecorderType_name, CreateOptions_RecorderType_value)
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// PoplarMetricsRecorderClient is the client API for PoplarMetricsRecorder service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PoplarMetricsRecorderClient interface {
	// Create builds a new recorder instance which creates a local file,
	// while the close recorder method flushes the contents of that
	// recorder and closes the file.
	CreateRecorder(ctx context.Context, in *CreateOptions, opts ...grpc.CallOption) (*EventResponse, error)
	CloseRecorder(ctx context.Context, in *RecorderID, opts ...grpc.CallOption) (*EventResponse, error)
	// Event Lifecycle methods
	BeginEvent(ctx context.Context, in *RecorderID, opts ...grpc.CallOption) (*EventResponse, error)
	ResetEvent(ctx context.Context, in *RecorderID, opts ...grpc.CallOption) (*EventResponse, error)
	EndEvent(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*EventResponse, error)
	// Timers
	SetTime(ctx context.Context, in *EventSendTime, opts ...grpc.CallOption) (*EventResponse, error)
	SetDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*EventResponse, error)
	SetTotalDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*EventResponse, error)
	// Guages
	SetState(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error)
	SetWorkers(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error)
	SetFailed(ctx context.Context, in *EventSendBool, opts ...grpc.CallOption) (*EventResponse, error)
	// Counters
	IncOps(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error)
	IncSize(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error)
	IncError(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error)
	IncIterations(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error)
}

type poplarMetricsRecorderClient struct {
	cc *grpc.ClientConn
}

func NewPoplarMetricsRecorderClient(cc *grpc.ClientConn) PoplarMetricsRecorderClient {
	return &poplarMetricsRecorderClient{cc}
}

func (c *poplarMetricsRecorderClient) CreateRecorder(ctx context.Context, in *CreateOptions, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/CreateRecorder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) CloseRecorder(ctx context.Context, in *RecorderID, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/CloseRecorder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) BeginEvent(ctx context.Context, in *RecorderID, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/BeginEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) ResetEvent(ctx context.Context, in *RecorderID, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/ResetEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) EndEvent(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/EndEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetTime(ctx context.Context, in *EventSendTime, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetDuration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetTotalDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetTotalDuration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetState(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetWorkers(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetWorkers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetFailed(ctx context.Context, in *EventSendBool, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetFailed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncOps(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/IncOps", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncSize(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/IncSize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncError(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/IncError", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncIterations(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*EventResponse, error) {
	out := new(EventResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/IncIterations", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PoplarMetricsRecorderServer is the server API for PoplarMetricsRecorder service.
type PoplarMetricsRecorderServer interface {
	// Create builds a new recorder instance which creates a local file,
	// while the close recorder method flushes the contents of that
	// recorder and closes the file.
	CreateRecorder(context.Context, *CreateOptions) (*EventResponse, error)
	CloseRecorder(context.Context, *RecorderID) (*EventResponse, error)
	// Event Lifecycle methods
	BeginEvent(context.Context, *RecorderID) (*EventResponse, error)
	ResetEvent(context.Context, *RecorderID) (*EventResponse, error)
	EndEvent(context.Context, *EventSendDuration) (*EventResponse, error)
	// Timers
	SetTime(context.Context, *EventSendTime) (*EventResponse, error)
	SetDuration(context.Context, *EventSendDuration) (*EventResponse, error)
	SetTotalDuration(context.Context, *EventSendDuration) (*EventResponse, error)
	// Guages
	SetState(context.Context, *EventSendInt) (*EventResponse, error)
	SetWorkers(context.Context, *EventSendInt) (*EventResponse, error)
	SetFailed(context.Context, *EventSendBool) (*EventResponse, error)
	// Counters
	IncOps(context.Context, *EventSendInt) (*EventResponse, error)
	IncSize(context.Context, *EventSendInt) (*EventResponse, error)
	IncError(context.Context, *EventSendInt) (*EventResponse, error)
	IncIterations(context.Context, *EventSendInt) (*EventResponse, error)
}

func RegisterPoplarMetricsRecorderServer(s *grpc.Server, srv PoplarMetricsRecorderServer) {
	s.RegisterService(&_PoplarMetricsRecorder_serviceDesc, srv)
}

func _PoplarMetricsRecorder_CreateRecorder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateOptions)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).CreateRecorder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/CreateRecorder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).CreateRecorder(ctx, req.(*CreateOptions))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_CloseRecorder_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecorderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).CloseRecorder(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/CloseRecorder",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).CloseRecorder(ctx, req.(*RecorderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_BeginEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecorderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).BeginEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/BeginEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).BeginEvent(ctx, req.(*RecorderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_ResetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecorderID)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).ResetEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/ResetEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).ResetEvent(ctx, req.(*RecorderID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_EndEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendDuration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).EndEvent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/EndEvent",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).EndEvent(ctx, req.(*EventSendDuration))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_SetTime_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendTime)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).SetTime(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/SetTime",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).SetTime(ctx, req.(*EventSendTime))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_SetDuration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendDuration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).SetDuration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/SetDuration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).SetDuration(ctx, req.(*EventSendDuration))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_SetTotalDuration_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendDuration)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).SetTotalDuration(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/SetTotalDuration",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).SetTotalDuration(ctx, req.(*EventSendDuration))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_SetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendInt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).SetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/SetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).SetState(ctx, req.(*EventSendInt))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_SetWorkers_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendInt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).SetWorkers(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/SetWorkers",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).SetWorkers(ctx, req.(*EventSendInt))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_SetFailed_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendBool)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).SetFailed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/SetFailed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).SetFailed(ctx, req.(*EventSendBool))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_IncOps_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendInt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).IncOps(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/IncOps",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).IncOps(ctx, req.(*EventSendInt))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_IncSize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendInt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).IncSize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/IncSize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).IncSize(ctx, req.(*EventSendInt))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_IncError_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendInt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).IncError(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/IncError",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).IncError(ctx, req.(*EventSendInt))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_IncIterations_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendInt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).IncIterations(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/IncIterations",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).IncIterations(ctx, req.(*EventSendInt))
	}
	return interceptor(ctx, in, info, handler)
}

var _PoplarMetricsRecorder_serviceDesc = grpc.ServiceDesc{
	ServiceName: "poplar.PoplarMetricsRecorder",
	HandlerType: (*PoplarMetricsRecorderServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateRecorder",
			Handler:    _PoplarMetricsRecorder_CreateRecorder_Handler,
		},
		{
			MethodName: "CloseRecorder",
			Handler:    _PoplarMetricsRecorder_CloseRecorder_Handler,
		},
		{
			MethodName: "BeginEvent",
			Handler:    _PoplarMetricsRecorder_BeginEvent_Handler,
		},
		{
			MethodName: "ResetEvent",
			Handler:    _PoplarMetricsRecorder_ResetEvent_Handler,
		},
		{
			MethodName: "EndEvent",
			Handler:    _PoplarMetricsRecorder_EndEvent_Handler,
		},
		{
			MethodName: "SetTime",
			Handler:    _PoplarMetricsRecorder_SetTime_Handler,
		},
		{
			MethodName: "SetDuration",
			Handler:    _PoplarMetricsRecorder_SetDuration_Handler,
		},
		{
			MethodName: "SetTotalDuration",
			Handler:    _PoplarMetricsRecorder_SetTotalDuration_Handler,
		},
		{
			MethodName: "SetState",
			Handler:    _PoplarMetricsRecorder_SetState_Handler,
		},
		{
			MethodName: "SetWorkers",
			Handler:    _PoplarMetricsRecorder_SetWorkers_Handler,
		},
		{
			MethodName: "SetFailed",
			Handler:    _PoplarMetricsRecorder_SetFailed_Handler,
		},
		{
			MethodName: "IncOps",
			Handler:    _PoplarMetricsRecorder_IncOps_Handler,
		},
		{
			MethodName: "IncSize",
			Handler:    _PoplarMetricsRecorder_IncSize_Handler,
		},
		{
			MethodName: "IncError",
			Handler:    _PoplarMetricsRecorder_IncError_Handler,
		},
		{
			MethodName: "IncIterations",
			Handler:    _PoplarMetricsRecorder_IncIterations_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "recorder.proto",
}

func init() { proto.RegisterFile("recorder.proto", fileDescriptor_recorder_a9549d9f8d4fbcd2) }

var fileDescriptor_recorder_a9549d9f8d4fbcd2 = []byte{
	// 649 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xd1, 0x4e, 0xdb, 0x3c,
	0x14, 0xc7, 0xbf, 0x42, 0x69, 0xd3, 0x53, 0x5a, 0xf2, 0x79, 0x30, 0x85, 0x6a, 0xda, 0xaa, 0x5c,
	0xf5, 0x2a, 0x30, 0x26, 0xc6, 0xd8, 0x10, 0xd3, 0x80, 0xc2, 0xa2, 0x8d, 0x16, 0xc5, 0x9d, 0x90,
	0x76, 0xb1, 0xc9, 0xa4, 0x67, 0x25, 0x22, 0xb5, 0x23, 0xc7, 0x45, 0x62, 0x4f, 0xb1, 0xcb, 0x3d,
	0xde, 0x1e, 0x65, 0x4a, 0xd2, 0xa4, 0x65, 0xa4, 0x9a, 0xa2, 0xdd, 0xc5, 0x7f, 0x9f, 0xdf, 0xf9,
	0x9f, 0x73, 0x62, 0x1b, 0x9a, 0x12, 0x5d, 0x21, 0x87, 0x28, 0xad, 0x40, 0x0a, 0x25, 0x48, 0x25,
	0x10, 0x81, 0xcf, 0x64, 0xeb, 0xd9, 0x48, 0x88, 0x91, 0x8f, 0x5b, 0xb1, 0x7a, 0x35, 0xf9, 0xb6,
	0xa5, 0xbc, 0x31, 0x86, 0x8a, 0x8d, 0x83, 0x24, 0xb0, 0xf5, 0xf4, 0xcf, 0x80, 0xe1, 0x44, 0x32,
	0xe5, 0x09, 0x9e, 0xec, 0x9b, 0x6d, 0x00, 0x67, 0x9a, 0xda, 0x3e, 0x21, 0x04, 0xca, 0x9c, 0x8d,
	0xd1, 0x28, 0xb5, 0x4b, 0x9d, 0x9a, 0x13, 0x7f, 0x9b, 0x14, 0x1a, 0xdd, 0x5b, 0xe4, 0x8a, 0x22,
	0x1f, 0x0e, 0xbc, 0x31, 0xe6, 0x05, 0x11, 0x0b, 0xca, 0x91, 0xb3, 0xb1, 0xd4, 0x2e, 0x75, 0xea,
	0x3b, 0x2d, 0x2b, 0x71, 0xb5, 0x52, 0x57, 0x6b, 0x90, 0x96, 0xe5, 0xc4, 0x71, 0xe6, 0x2b, 0x58,
	0xcd, 0x92, 0xda, 0x5c, 0xe5, 0xe6, 0x5c, 0x87, 0x95, 0x5b, 0xe6, 0x4f, 0x92, 0xa4, 0xcb, 0x4e,
	0xb2, 0x30, 0xf7, 0xe7, 0xca, 0x39, 0x12, 0xc2, 0xff, 0x3b, 0xaa, 0xa5, 0xe8, 0x17, 0xf8, 0x3f,
	0x43, 0x4f, 0xa6, 0x63, 0xc8, 0xc5, 0x77, 0x41, 0x4b, 0xc7, 0x34, 0xed, 0x68, 0xf3, 0x41, 0x47,
	0x69, 0x02, 0x27, 0x0b, 0x35, 0xdf, 0x4c, 0x4b, 0x73, 0x30, 0x0c, 0x04, 0x0f, 0xf3, 0x27, 0xf5,
	0x18, 0x2a, 0xa1, 0x62, 0x6a, 0x12, 0x4e, 0x6b, 0x9b, 0xae, 0xcc, 0x5f, 0x4b, 0xd0, 0x38, 0x96,
	0xc8, 0x14, 0xf6, 0x83, 0x28, 0x5b, 0x98, 0x4b, 0x13, 0x28, 0x07, 0x4c, 0x5d, 0xc7, 0x6c, 0xcd,
	0x89, 0xbf, 0xc9, 0x13, 0xa8, 0xb9, 0xd7, 0x13, 0x7e, 0x43, 0xbd, 0xef, 0x68, 0x2c, 0xb7, 0x4b,
	0x9d, 0x15, 0x67, 0x26, 0x44, 0xbb, 0xa1, 0x92, 0xc8, 0xc6, 0x1e, 0x1f, 0x19, 0xe5, 0xd8, 0x72,
	0x26, 0x10, 0x03, 0xaa, 0xc3, 0x3b, 0xce, 0xc6, 0x9e, 0x6b, 0xac, 0xc4, 0x7b, 0xe9, 0x92, 0x1c,
	0x82, 0x96, 0x9e, 0x39, 0xa3, 0xd2, 0x2e, 0x75, 0x9a, 0x3b, 0xa6, 0x95, 0x1c, 0x3a, 0xeb, 0x5e,
	0x99, 0x56, 0x7a, 0x7c, 0x06, 0x77, 0x01, 0x3a, 0x19, 0x63, 0xfe, 0x28, 0xc1, 0xea, 0xfc, 0x16,
	0xa9, 0x43, 0xf5, 0x53, 0xef, 0x43, 0xaf, 0x7f, 0xd9, 0xd3, 0xff, 0x23, 0x1a, 0x94, 0x2f, 0xba,
	0xce, 0xa9, 0x5e, 0x22, 0x6b, 0x50, 0x8f, 0xbe, 0xbe, 0x52, 0xbb, 0x77, 0xf6, 0xb1, 0xab, 0x2f,
	0x91, 0x26, 0x40, 0x2c, 0x3c, 0xdf, 0xde, 0x3e, 0xa7, 0xfa, 0x72, 0xc4, 0x25, 0x6b, 0xaa, 0x97,
	0xc9, 0x3a, 0xe8, 0xef, 0x6d, 0x3a, 0xe8, 0x9f, 0x39, 0xef, 0xce, 0x53, 0xa4, 0x42, 0x1e, 0xc1,
	0xda, 0x4c, 0x4d, 0xb8, 0x2a, 0xd1, 0x61, 0x75, 0x4e, 0xa4, 0xba, 0xb6, 0xf3, 0xb3, 0x0a, 0x1b,
	0x17, 0x71, 0x0b, 0xe7, 0xa8, 0xa4, 0xe7, 0x86, 0x69, 0x7d, 0xe4, 0x10, 0x9a, 0x49, 0x53, 0x99,
	0xb2, 0x91, 0xdb, 0x6c, 0x2b, 0x93, 0xef, 0xff, 0xe8, 0xd7, 0xd0, 0x38, 0xf6, 0x45, 0x38, 0xc3,
	0x49, 0x1a, 0x37, 0xbb, 0x5c, 0x8b, 0xd8, 0x3d, 0x80, 0x23, 0x1c, 0x79, 0x3c, 0x56, 0x0b, 0x82,
	0x0e, 0x86, 0xa8, 0x0a, 0x83, 0x07, 0xa0, 0x75, 0xf9, 0x30, 0xc1, 0x36, 0xef, 0x85, 0xcc, 0xdf,
	0x8c, 0xc5, 0xb6, 0x55, 0x8a, 0x2a, 0x7e, 0x09, 0x36, 0x1e, 0xc0, 0x91, 0xbc, 0x08, 0x7c, 0x0b,
	0x75, 0x8a, 0x2a, 0xbb, 0x78, 0xc5, 0x9d, 0x4f, 0x40, 0x8f, 0x9c, 0x85, 0x62, 0xfe, 0x3f, 0x64,
	0xd9, 0x03, 0x8d, 0xa2, 0xa2, 0x8a, 0x29, 0x24, 0xeb, 0x0f, 0x68, 0x9b, 0xab, 0x45, 0xe0, 0x3e,
	0x00, 0x45, 0x75, 0x29, 0xe4, 0x0d, 0xca, 0xb0, 0x28, 0x5a, 0xa3, 0xa8, 0x4e, 0x99, 0xe7, 0xe3,
	0x30, 0x67, 0x6a, 0xd1, 0x3b, 0xb6, 0x08, 0xdd, 0x85, 0x8a, 0xcd, 0xdd, 0x7e, 0x50, 0xd0, 0xf1,
	0x25, 0x54, 0x6d, 0xee, 0xc6, 0x2f, 0x40, 0x21, 0x6e, 0x0f, 0x34, 0x9b, 0xbb, 0x5d, 0x29, 0x85,
	0x2c, 0x06, 0x1e, 0x40, 0xc3, 0xe6, 0xae, 0xad, 0x30, 0x99, 0x7e, 0xb1, 0x72, 0x8f, 0xe0, 0xb3,
	0xe6, 0x71, 0x85, 0x92, 0x33, 0xff, 0xaa, 0x12, 0xbf, 0xb1, 0x2f, 0x7e, 0x07, 0x00, 0x00, 0xff,
	0xff, 0x54, 0xdc, 0x37, 0x6f, 0xf4, 0x06, 0x00, 0x00,
}
