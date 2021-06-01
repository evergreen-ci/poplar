// Code generated by protoc-gen-go. DO NOT EDIT.
// source: google/ads/googleads/v0/services/topic_constant_service.proto

package services // import "google.golang.org/genproto/googleapis/ads/googleads/v0/services"

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import resources "google.golang.org/genproto/googleapis/ads/googleads/v0/resources"
import _ "google.golang.org/genproto/googleapis/api/annotations"

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

// Request message for [TopicConstantService.GetTopicConstant][google.ads.googleads.v0.services.TopicConstantService.GetTopicConstant].
type GetTopicConstantRequest struct {
	// Resource name of the Topic to fetch.
	ResourceName         string   `protobuf:"bytes,1,opt,name=resource_name,json=resourceName,proto3" json:"resource_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetTopicConstantRequest) Reset()         { *m = GetTopicConstantRequest{} }
func (m *GetTopicConstantRequest) String() string { return proto.CompactTextString(m) }
func (*GetTopicConstantRequest) ProtoMessage()    {}
func (*GetTopicConstantRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_topic_constant_service_8bf77bcf67597eb4, []int{0}
}
func (m *GetTopicConstantRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetTopicConstantRequest.Unmarshal(m, b)
}
func (m *GetTopicConstantRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetTopicConstantRequest.Marshal(b, m, deterministic)
}
func (dst *GetTopicConstantRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetTopicConstantRequest.Merge(dst, src)
}
func (m *GetTopicConstantRequest) XXX_Size() int {
	return xxx_messageInfo_GetTopicConstantRequest.Size(m)
}
func (m *GetTopicConstantRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetTopicConstantRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetTopicConstantRequest proto.InternalMessageInfo

func (m *GetTopicConstantRequest) GetResourceName() string {
	if m != nil {
		return m.ResourceName
	}
	return ""
}

func init() {
	proto.RegisterType((*GetTopicConstantRequest)(nil), "google.ads.googleads.v0.services.GetTopicConstantRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// TopicConstantServiceClient is the client API for TopicConstantService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type TopicConstantServiceClient interface {
	// Returns the requested topic constant in full detail.
	GetTopicConstant(ctx context.Context, in *GetTopicConstantRequest, opts ...grpc.CallOption) (*resources.TopicConstant, error)
}

type topicConstantServiceClient struct {
	cc *grpc.ClientConn
}

func NewTopicConstantServiceClient(cc *grpc.ClientConn) TopicConstantServiceClient {
	return &topicConstantServiceClient{cc}
}

func (c *topicConstantServiceClient) GetTopicConstant(ctx context.Context, in *GetTopicConstantRequest, opts ...grpc.CallOption) (*resources.TopicConstant, error) {
	out := new(resources.TopicConstant)
	err := c.cc.Invoke(ctx, "/google.ads.googleads.v0.services.TopicConstantService/GetTopicConstant", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TopicConstantServiceServer is the server API for TopicConstantService service.
type TopicConstantServiceServer interface {
	// Returns the requested topic constant in full detail.
	GetTopicConstant(context.Context, *GetTopicConstantRequest) (*resources.TopicConstant, error)
}

func RegisterTopicConstantServiceServer(s *grpc.Server, srv TopicConstantServiceServer) {
	s.RegisterService(&_TopicConstantService_serviceDesc, srv)
}

func _TopicConstantService_GetTopicConstant_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTopicConstantRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TopicConstantServiceServer).GetTopicConstant(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/google.ads.googleads.v0.services.TopicConstantService/GetTopicConstant",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TopicConstantServiceServer).GetTopicConstant(ctx, req.(*GetTopicConstantRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _TopicConstantService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "google.ads.googleads.v0.services.TopicConstantService",
	HandlerType: (*TopicConstantServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetTopicConstant",
			Handler:    _TopicConstantService_GetTopicConstant_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "google/ads/googleads/v0/services/topic_constant_service.proto",
}

func init() {
	proto.RegisterFile("google/ads/googleads/v0/services/topic_constant_service.proto", fileDescriptor_topic_constant_service_8bf77bcf67597eb4)
}

var fileDescriptor_topic_constant_service_8bf77bcf67597eb4 = []byte{
	// 354 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xb2, 0x4d, 0xcf, 0xcf, 0x4f,
	0xcf, 0x49, 0xd5, 0x4f, 0x4c, 0x29, 0xd6, 0x87, 0x30, 0x41, 0xac, 0x32, 0x03, 0xfd, 0xe2, 0xd4,
	0xa2, 0xb2, 0xcc, 0xe4, 0xd4, 0x62, 0xfd, 0x92, 0xfc, 0x82, 0xcc, 0xe4, 0xf8, 0xe4, 0xfc, 0xbc,
	0xe2, 0x92, 0xc4, 0xbc, 0x92, 0x78, 0xa8, 0xb8, 0x5e, 0x41, 0x51, 0x7e, 0x49, 0xbe, 0x90, 0x02,
	0x44, 0x8f, 0x5e, 0x62, 0x4a, 0xb1, 0x1e, 0x5c, 0xbb, 0x5e, 0x99, 0x81, 0x1e, 0x4c, 0xbb, 0x94,
	0x19, 0x2e, 0x0b, 0x8a, 0x52, 0x8b, 0xf3, 0x4b, 0x8b, 0x30, 0x6d, 0x80, 0x98, 0x2c, 0x25, 0x03,
	0xd3, 0x57, 0x90, 0xa9, 0x9f, 0x98, 0x97, 0x97, 0x5f, 0x92, 0x58, 0x92, 0x99, 0x9f, 0x57, 0x0c,
	0x91, 0x55, 0xb2, 0xe3, 0x12, 0x77, 0x4f, 0x2d, 0x09, 0x01, 0x69, 0x74, 0x86, 0xea, 0x0b, 0x4a,
	0x2d, 0x2c, 0x4d, 0x2d, 0x2e, 0x11, 0x52, 0xe6, 0xe2, 0x85, 0x19, 0x1d, 0x9f, 0x97, 0x98, 0x9b,
	0x2a, 0xc1, 0xa8, 0xc0, 0xa8, 0xc1, 0x19, 0xc4, 0x03, 0x13, 0xf4, 0x4b, 0xcc, 0x4d, 0x35, 0x3a,
	0xc6, 0xc8, 0x25, 0x82, 0xa2, 0x3b, 0x18, 0xe2, 0x5e, 0xa1, 0xb5, 0x8c, 0x5c, 0x02, 0xe8, 0x26,
	0x0b, 0x59, 0xea, 0x11, 0xf2, 0xa6, 0x1e, 0x0e, 0xd7, 0x48, 0x19, 0xe0, 0xd4, 0x0a, 0xf7, 0xbf,
	0x1e, 0x8a, 0x46, 0x25, 0x9d, 0xa6, 0xcb, 0x4f, 0x26, 0x33, 0xa9, 0x09, 0xa9, 0x80, 0x02, 0xa9,
	0x1a, 0xc5, 0x2b, 0xb6, 0x25, 0xc8, 0x2a, 0x8b, 0xf5, 0xb5, 0x6a, 0x9d, 0x1a, 0x98, 0xb8, 0x54,
	0x92, 0xf3, 0x73, 0x09, 0x3a, 0xd0, 0x49, 0x12, 0x9b, 0x77, 0x03, 0x40, 0x81, 0x19, 0xc0, 0x18,
	0xe5, 0x01, 0xd5, 0x9e, 0x9e, 0x9f, 0x93, 0x98, 0x97, 0xae, 0x97, 0x5f, 0x94, 0xae, 0x9f, 0x9e,
	0x9a, 0x07, 0x0e, 0x6a, 0x58, 0xa4, 0x15, 0x64, 0x16, 0xe3, 0x4e, 0x24, 0xd6, 0x30, 0xc6, 0x22,
	0x26, 0x66, 0x77, 0x47, 0xc7, 0x55, 0x4c, 0x0a, 0xee, 0x10, 0x03, 0x1d, 0x53, 0x8a, 0xf5, 0x20,
	0x4c, 0x10, 0x2b, 0xcc, 0x40, 0x0f, 0x6a, 0x71, 0xf1, 0x29, 0x98, 0x92, 0x18, 0xc7, 0x94, 0xe2,
	0x18, 0xb8, 0x92, 0x98, 0x30, 0x83, 0x18, 0x98, 0x92, 0x57, 0x4c, 0x2a, 0x10, 0x71, 0x2b, 0x2b,
	0xc7, 0x94, 0x62, 0x2b, 0x2b, 0xb8, 0x22, 0x2b, 0xab, 0x30, 0x03, 0x2b, 0x2b, 0x98, 0xb2, 0x24,
	0x36, 0xb0, 0x3b, 0x8d, 0x01, 0x01, 0x00, 0x00, 0xff, 0xff, 0x15, 0xb8, 0x80, 0x62, 0xcb, 0x02,
	0x00, 0x00,
}
