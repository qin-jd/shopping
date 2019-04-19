// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/order/order.proto

package go_micro_srv_order

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type SubmitRequest struct {
	ProductId            uint32   `protobuf:"varint,1,opt,name=productId,proto3" json:"productId,omitempty"`
	Count                uint32   `protobuf:"varint,2,opt,name=count,proto3" json:"count,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SubmitRequest) Reset()         { *m = SubmitRequest{} }
func (m *SubmitRequest) String() string { return proto.CompactTextString(m) }
func (*SubmitRequest) ProtoMessage()    {}
func (*SubmitRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_986e030a471601a2, []int{0}
}

func (m *SubmitRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SubmitRequest.Unmarshal(m, b)
}
func (m *SubmitRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SubmitRequest.Marshal(b, m, deterministic)
}
func (m *SubmitRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SubmitRequest.Merge(m, src)
}
func (m *SubmitRequest) XXX_Size() int {
	return xxx_messageInfo_SubmitRequest.Size(m)
}
func (m *SubmitRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SubmitRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SubmitRequest proto.InternalMessageInfo

func (m *SubmitRequest) GetProductId() uint32 {
	if m != nil {
		return m.ProductId
	}
	return 0
}

func (m *SubmitRequest) GetCount() uint32 {
	if m != nil {
		return m.Count
	}
	return 0
}

type Response struct {
	Code                 string   `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Msg                  string   `protobuf:"bytes,2,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_986e030a471601a2, []int{1}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetCode() string {
	if m != nil {
		return m.Code
	}
	return ""
}

func (m *Response) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type OrderDetailRequest struct {
	OrderId              string   `protobuf:"bytes,1,opt,name=orderId,proto3" json:"orderId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *OrderDetailRequest) Reset()         { *m = OrderDetailRequest{} }
func (m *OrderDetailRequest) String() string { return proto.CompactTextString(m) }
func (*OrderDetailRequest) ProtoMessage()    {}
func (*OrderDetailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_986e030a471601a2, []int{2}
}

func (m *OrderDetailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OrderDetailRequest.Unmarshal(m, b)
}
func (m *OrderDetailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OrderDetailRequest.Marshal(b, m, deterministic)
}
func (m *OrderDetailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OrderDetailRequest.Merge(m, src)
}
func (m *OrderDetailRequest) XXX_Size() int {
	return xxx_messageInfo_OrderDetailRequest.Size(m)
}
func (m *OrderDetailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_OrderDetailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_OrderDetailRequest proto.InternalMessageInfo

func (m *OrderDetailRequest) GetOrderId() string {
	if m != nil {
		return m.OrderId
	}
	return ""
}

func init() {
	proto.RegisterType((*SubmitRequest)(nil), "go.micro.srv.order.SubmitRequest")
	proto.RegisterType((*Response)(nil), "go.micro.srv.order.Response")
	proto.RegisterType((*OrderDetailRequest)(nil), "go.micro.srv.order.OrderDetailRequest")
}

func init() { proto.RegisterFile("proto/order/order.proto", fileDescriptor_986e030a471601a2) }

var fileDescriptor_986e030a471601a2 = []byte{
	// 233 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x84, 0x90, 0xcd, 0x4a, 0x03, 0x31,
	0x10, 0xc7, 0x5d, 0x3f, 0xaa, 0x19, 0x2d, 0xc8, 0x20, 0xb8, 0x48, 0x0f, 0x9a, 0x83, 0x78, 0x8a,
	0xa2, 0x8f, 0xa0, 0x17, 0xf1, 0x20, 0xa4, 0xf8, 0x00, 0x36, 0x19, 0x4a, 0xc0, 0xed, 0xac, 0xf9,
	0xe8, 0x83, 0xf9, 0x84, 0xd2, 0x59, 0x17, 0x95, 0x5d, 0xe8, 0x25, 0xcc, 0xd7, 0x3f, 0xff, 0xdf,
	0x0c, 0x9c, 0xb7, 0x91, 0x33, 0xdf, 0x72, 0xf4, 0x14, 0xbb, 0xd7, 0x48, 0x05, 0x71, 0xc9, 0xa6,
	0x09, 0x2e, 0xb2, 0x49, 0x71, 0x6d, 0xa4, 0xa3, 0x1f, 0x61, 0x3a, 0x2f, 0x8b, 0x26, 0x64, 0x4b,
	0x9f, 0x85, 0x52, 0xc6, 0x19, 0xa8, 0x36, 0xb2, 0x2f, 0x2e, 0x3f, 0xfb, 0xba, 0xba, 0xac, 0x6e,
	0xa6, 0xf6, 0xb7, 0x80, 0x67, 0x70, 0xe0, 0xb8, 0xac, 0x72, 0xbd, 0x2b, 0x9d, 0x2e, 0xd1, 0x77,
	0x70, 0x64, 0x29, 0xb5, 0xbc, 0x4a, 0x84, 0x08, 0xfb, 0x8e, 0x3d, 0x89, 0x54, 0x59, 0x89, 0xf1,
	0x14, 0xf6, 0x9a, 0xb4, 0x14, 0x8d, 0xb2, 0x9b, 0x50, 0x1b, 0xc0, 0xd7, 0x8d, 0xff, 0x13, 0xe5,
	0xf7, 0xf0, 0xd1, 0x7b, 0xd7, 0x70, 0x28, 0x54, 0x3f, 0xce, 0xca, 0xf6, 0xe9, 0xfd, 0x57, 0x05,
	0x27, 0x22, 0x98, 0x53, 0x5c, 0x07, 0x47, 0xf8, 0x02, 0x93, 0x8e, 0x1b, 0xaf, 0xcc, 0x70, 0x2d,
	0xf3, 0x6f, 0xa7, 0x8b, 0xd9, 0xd8, 0x48, 0x4f, 0xac, 0x77, 0xf0, 0x0d, 0x8e, 0xff, 0xd0, 0xe0,
	0xf5, 0xd8, 0xf8, 0x10, 0x77, 0xdb, 0xb7, 0x8b, 0x89, 0x9c, 0xfd, 0xe1, 0x3b, 0x00, 0x00, 0xff,
	0xff, 0x31, 0xe3, 0xa7, 0x27, 0x91, 0x01, 0x00, 0x00,
}