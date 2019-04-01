// Code generated by protoc-gen-go. DO NOT EDIT.
// source: recorder.proto

package internal

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	grpc "google.golang.org/grpc"
	math "math"
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
	return fileDescriptor_b063ffe85a4e6395, []int{0}
}

func (m *EventSendTime) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendTime.Unmarshal(m, b)
}
func (m *EventSendTime) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendTime.Marshal(b, m, deterministic)
}
func (m *EventSendTime) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendTime.Merge(m, src)
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
	return fileDescriptor_b063ffe85a4e6395, []int{1}
}

func (m *EventSendInt) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendInt.Unmarshal(m, b)
}
func (m *EventSendInt) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendInt.Marshal(b, m, deterministic)
}
func (m *EventSendInt) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendInt.Merge(m, src)
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
	return fileDescriptor_b063ffe85a4e6395, []int{2}
}

func (m *EventSendBool) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendBool.Unmarshal(m, b)
}
func (m *EventSendBool) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendBool.Marshal(b, m, deterministic)
}
func (m *EventSendBool) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendBool.Merge(m, src)
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
	return fileDescriptor_b063ffe85a4e6395, []int{3}
}

func (m *EventSendDuration) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EventSendDuration.Unmarshal(m, b)
}
func (m *EventSendDuration) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EventSendDuration.Marshal(b, m, deterministic)
}
func (m *EventSendDuration) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EventSendDuration.Merge(m, src)
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

func init() {
	proto.RegisterType((*EventSendTime)(nil), "poplar.EventSendTime")
	proto.RegisterType((*EventSendInt)(nil), "poplar.EventSendInt")
	proto.RegisterType((*EventSendBool)(nil), "poplar.EventSendBool")
	proto.RegisterType((*EventSendDuration)(nil), "poplar.EventSendDuration")
}

func init() { proto.RegisterFile("recorder.proto", fileDescriptor_b063ffe85a4e6395) }

var fileDescriptor_b063ffe85a4e6395 = []byte{
	// 458 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xcf, 0x6f, 0xd3, 0x30,
	0x14, 0xc7, 0x55, 0xd8, 0xda, 0xec, 0x6d, 0x9d, 0x46, 0xb4, 0xa2, 0xad, 0x07, 0x98, 0x7a, 0xda,
	0x29, 0x93, 0x86, 0x06, 0x1d, 0xd2, 0x04, 0xea, 0x5a, 0x24, 0x1f, 0xd0, 0x50, 0x3c, 0x09, 0x89,
	0x03, 0x92, 0x97, 0x3c, 0x2a, 0x0b, 0xd7, 0x8e, 0xec, 0xd7, 0x1d, 0xf8, 0x4b, 0xf9, 0x73, 0x50,
	0xe2, 0x24, 0xfc, 0x68, 0x0b, 0x32, 0x3b, 0xb5, 0x7e, 0xf9, 0x7e, 0xbe, 0x7e, 0xbf, 0x12, 0xd8,
	0xb7, 0x98, 0x19, 0x9b, 0xa3, 0x4d, 0x0a, 0x6b, 0xc8, 0xc4, 0xdd, 0xc2, 0x14, 0x4a, 0xd8, 0xe1,
	0x9e, 0xff, 0xf5, 0xd1, 0xe1, 0xf3, 0xb9, 0x31, 0x73, 0x85, 0x67, 0xd5, 0xe9, 0x6e, 0xf9, 0xe5,
	0x8c, 0xe4, 0x02, 0x1d, 0x89, 0x45, 0x51, 0x0b, 0x9e, 0xfd, 0x29, 0xc8, 0x97, 0x56, 0x90, 0x34,
	0xda, 0x3f, 0x1f, 0x71, 0xe8, 0xcf, 0xee, 0x51, 0x13, 0x47, 0x9d, 0xdf, 0xca, 0x05, 0xc6, 0x31,
	0x6c, 0x69, 0xb1, 0xc0, 0xa3, 0xce, 0x49, 0xe7, 0x74, 0x27, 0xad, 0xfe, 0xc7, 0x09, 0x6c, 0x95,
	0xbe, 0x47, 0x8f, 0x4e, 0x3a, 0xa7, 0xbb, 0xe7, 0xc3, 0xc4, 0x7b, 0x26, 0x8d, 0x67, 0x72, 0xdb,
	0x5c, 0x9a, 0x56, 0xba, 0xd1, 0x18, 0xf6, 0x5a, 0x53, 0xa6, 0x69, 0xad, 0xe7, 0x21, 0x6c, 0xdf,
	0x0b, 0xb5, 0xf4, 0xa6, 0x8f, 0x53, 0x7f, 0x18, 0x5d, 0xfe, 0x92, 0xce, 0xc4, 0x18, 0xf5, 0x6f,
	0x34, 0x6a, 0xd0, 0xcf, 0xf0, 0xa4, 0x45, 0xa7, 0x75, 0x91, 0x6b, 0xf1, 0x0b, 0x88, 0x9a, 0x26,
	0xd4, 0x15, 0x1d, 0xaf, 0x54, 0xd4, 0x18, 0xa4, 0xad, 0xf4, 0xfc, 0x7b, 0x0f, 0x06, 0x1f, 0xaa,
	0xde, 0xbf, 0x47, 0xb2, 0x32, 0x73, 0x69, 0x3d, 0xa0, 0xf8, 0x0d, 0xec, 0x5f, 0x5b, 0x14, 0x84,
	0x6d, 0x64, 0x90, 0xd4, 0x53, 0xf2, 0xf1, 0x9b, 0xa2, 0x74, 0x70, 0xc3, 0xa7, 0x4d, 0xd8, 0xfb,
	0xa4, 0xe8, 0x0a, 0xa3, 0x1d, 0xc6, 0x97, 0xd0, 0xbf, 0x56, 0xc6, 0xfd, 0xe4, 0x0f, 0x7e, 0x17,
	0xb2, 0xe9, 0x46, 0xf4, 0x25, 0xc0, 0x04, 0xe7, 0x52, 0x57, 0xa5, 0x87, 0x71, 0x29, 0x3a, 0xa4,
	0x50, 0xee, 0x0a, 0xa2, 0x99, 0xce, 0x3d, 0x75, 0xdc, 0x68, 0x56, 0xfa, 0xbe, 0x11, 0xbf, 0x80,
	0x6d, 0x8e, 0xc4, 0xa6, 0xf1, 0xe1, 0x0a, 0xcb, 0x34, 0x6d, 0xc4, 0xc6, 0xd0, 0xe3, 0x48, 0xd5,
	0x7e, 0x0e, 0x56, 0xc0, 0x32, 0xbc, 0x91, 0x7c, 0x0b, 0xbb, 0x1c, 0xa9, 0xdd, 0x87, 0xff, 0x48,
	0x79, 0x06, 0x07, 0xe5, 0xdd, 0x86, 0x84, 0x7a, 0x88, 0xcd, 0x18, 0x22, 0x8e, 0xc4, 0x49, 0x10,
	0x06, 0x16, 0xff, 0x1a, 0x80, 0x23, 0x7d, 0x34, 0xf6, 0x2b, 0x5a, 0x17, 0xcc, 0xee, 0x70, 0xa4,
	0x77, 0x42, 0x2a, 0xcc, 0xd7, 0xb4, 0xae, 0x7c, 0xc5, 0xfe, 0xb2, 0x22, 0x5d, 0xa6, 0xb3, 0x9b,
	0x22, 0xf4, 0xce, 0x57, 0xd0, 0x63, 0x3a, 0xe3, 0xf2, 0x1b, 0x06, 0x4f, 0x39, 0x62, 0x3a, 0x9b,
	0x59, 0x6b, 0x6c, 0x20, 0x79, 0x05, 0x7d, 0xa6, 0x33, 0x46, 0xe8, 0x87, 0x10, 0x98, 0xf1, 0x04,
	0x3e, 0x45, 0x52, 0x13, 0x5a, 0x2d, 0xd4, 0x5d, 0xb7, 0xfa, 0x06, 0xbc, 0xf8, 0x11, 0x00, 0x00,
	0xff, 0xff, 0x32, 0x12, 0x0e, 0x13, 0x80, 0x05, 0x00, 0x00,
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
	CreateRecorder(ctx context.Context, in *CreateOptions, opts ...grpc.CallOption) (*PoplarResponse, error)
	CloseRecorder(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error)
	// Event Lifecycle methods
	BeginEvent(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error)
	ResetEvent(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error)
	EndEvent(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*PoplarResponse, error)
	SetID(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error)
	// Timers
	SetTime(ctx context.Context, in *EventSendTime, opts ...grpc.CallOption) (*PoplarResponse, error)
	SetDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*PoplarResponse, error)
	SetTotalDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*PoplarResponse, error)
	// Guages
	SetState(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error)
	SetWorkers(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error)
	SetFailed(ctx context.Context, in *EventSendBool, opts ...grpc.CallOption) (*PoplarResponse, error)
	// Counters
	IncOps(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error)
	IncSize(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error)
	IncError(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error)
	IncIterations(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error)
}

type poplarMetricsRecorderClient struct {
	cc *grpc.ClientConn
}

func NewPoplarMetricsRecorderClient(cc *grpc.ClientConn) PoplarMetricsRecorderClient {
	return &poplarMetricsRecorderClient{cc}
}

func (c *poplarMetricsRecorderClient) CreateRecorder(ctx context.Context, in *CreateOptions, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/CreateRecorder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) CloseRecorder(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/CloseRecorder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) BeginEvent(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/BeginEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) ResetEvent(ctx context.Context, in *PoplarID, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/ResetEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) EndEvent(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/EndEvent", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetID(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetTime(ctx context.Context, in *EventSendTime, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetTime", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetDuration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetTotalDuration(ctx context.Context, in *EventSendDuration, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetTotalDuration", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetState(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetWorkers(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetWorkers", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) SetFailed(ctx context.Context, in *EventSendBool, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/SetFailed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncOps(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/IncOps", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncSize(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/IncSize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncError(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
	err := c.cc.Invoke(ctx, "/poplar.PoplarMetricsRecorder/IncError", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *poplarMetricsRecorderClient) IncIterations(ctx context.Context, in *EventSendInt, opts ...grpc.CallOption) (*PoplarResponse, error) {
	out := new(PoplarResponse)
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
	CreateRecorder(context.Context, *CreateOptions) (*PoplarResponse, error)
	CloseRecorder(context.Context, *PoplarID) (*PoplarResponse, error)
	// Event Lifecycle methods
	BeginEvent(context.Context, *PoplarID) (*PoplarResponse, error)
	ResetEvent(context.Context, *PoplarID) (*PoplarResponse, error)
	EndEvent(context.Context, *EventSendDuration) (*PoplarResponse, error)
	SetID(context.Context, *EventSendInt) (*PoplarResponse, error)
	// Timers
	SetTime(context.Context, *EventSendTime) (*PoplarResponse, error)
	SetDuration(context.Context, *EventSendDuration) (*PoplarResponse, error)
	SetTotalDuration(context.Context, *EventSendDuration) (*PoplarResponse, error)
	// Guages
	SetState(context.Context, *EventSendInt) (*PoplarResponse, error)
	SetWorkers(context.Context, *EventSendInt) (*PoplarResponse, error)
	SetFailed(context.Context, *EventSendBool) (*PoplarResponse, error)
	// Counters
	IncOps(context.Context, *EventSendInt) (*PoplarResponse, error)
	IncSize(context.Context, *EventSendInt) (*PoplarResponse, error)
	IncError(context.Context, *EventSendInt) (*PoplarResponse, error)
	IncIterations(context.Context, *EventSendInt) (*PoplarResponse, error)
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
	in := new(PoplarID)
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
		return srv.(PoplarMetricsRecorderServer).CloseRecorder(ctx, req.(*PoplarID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_BeginEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoplarID)
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
		return srv.(PoplarMetricsRecorderServer).BeginEvent(ctx, req.(*PoplarID))
	}
	return interceptor(ctx, in, info, handler)
}

func _PoplarMetricsRecorder_ResetEvent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PoplarID)
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
		return srv.(PoplarMetricsRecorderServer).ResetEvent(ctx, req.(*PoplarID))
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

func _PoplarMetricsRecorder_SetID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EventSendInt)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PoplarMetricsRecorderServer).SetID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/poplar.PoplarMetricsRecorder/SetID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PoplarMetricsRecorderServer).SetID(ctx, req.(*EventSendInt))
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
			MethodName: "SetID",
			Handler:    _PoplarMetricsRecorder_SetID_Handler,
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
