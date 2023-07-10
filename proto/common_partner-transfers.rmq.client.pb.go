// Code generated by protoc-gen-go. DO NOT EDIT.
// source: proto/common_partner-transfers.rmq.client.proto

package partner_transfers_defs

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	client "gitlab.fbs-d.com/definitions/client"
	_ "gitlab.fbs-d.com/dev/amqp-rpc-generator"
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

type StartTwoPhasedPartnerTransferV1Request struct {
	Token                string                                      `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
	Client               *client.Client                              `protobuf:"bytes,2,opt,name=client,proto3" json:"client"`
	Body                 *StartTwoPhasedPartnerTransferV1RequestBody `protobuf:"bytes,3,opt,name=body,proto3" json:"body"`
	XXX_NoUnkeyedLiteral struct{}                                    `json:"-"`
	XXX_unrecognized     []byte                                      `json:"-"`
	XXX_sizecache        int32                                       `json:"-"`
}

func (m *StartTwoPhasedPartnerTransferV1Request) Reset() {
	*m = StartTwoPhasedPartnerTransferV1Request{}
}
func (m *StartTwoPhasedPartnerTransferV1Request) String() string { return proto.CompactTextString(m) }
func (*StartTwoPhasedPartnerTransferV1Request) ProtoMessage()    {}
func (*StartTwoPhasedPartnerTransferV1Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_c00dd1a11d825a35, []int{0}
}

func (m *StartTwoPhasedPartnerTransferV1Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1Request.Unmarshal(m, b)
}
func (m *StartTwoPhasedPartnerTransferV1Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1Request.Marshal(b, m, deterministic)
}
func (m *StartTwoPhasedPartnerTransferV1Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartTwoPhasedPartnerTransferV1Request.Merge(m, src)
}
func (m *StartTwoPhasedPartnerTransferV1Request) XXX_Size() int {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1Request.Size(m)
}
func (m *StartTwoPhasedPartnerTransferV1Request) XXX_DiscardUnknown() {
	xxx_messageInfo_StartTwoPhasedPartnerTransferV1Request.DiscardUnknown(m)
}

var xxx_messageInfo_StartTwoPhasedPartnerTransferV1Request proto.InternalMessageInfo

func (m *StartTwoPhasedPartnerTransferV1Request) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *StartTwoPhasedPartnerTransferV1Request) GetClient() *client.Client {
	if m != nil {
		return m.Client
	}
	return nil
}

func (m *StartTwoPhasedPartnerTransferV1Request) GetBody() *StartTwoPhasedPartnerTransferV1RequestBody {
	if m != nil {
		return m.Body
	}
	return nil
}

type StartTwoPhasedPartnerTransferV1RequestBody struct {
	TransferKey          string   `protobuf:"bytes,1,opt,name=transferKey,proto3" json:"transferKey"`
	FromAccountId        int64    `protobuf:"varint,2,opt,name=fromAccountId,proto3" json:"fromAccountId"`
	ToAccountId          int64    `protobuf:"varint,3,opt,name=toAccountId,proto3" json:"toAccountId"`
	Amount               int64    `protobuf:"varint,4,opt,name=amount,proto3" json:"amount"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartTwoPhasedPartnerTransferV1RequestBody) Reset() {
	*m = StartTwoPhasedPartnerTransferV1RequestBody{}
}
func (m *StartTwoPhasedPartnerTransferV1RequestBody) String() string {
	return proto.CompactTextString(m)
}
func (*StartTwoPhasedPartnerTransferV1RequestBody) ProtoMessage() {}
func (*StartTwoPhasedPartnerTransferV1RequestBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_c00dd1a11d825a35, []int{1}
}

func (m *StartTwoPhasedPartnerTransferV1RequestBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1RequestBody.Unmarshal(m, b)
}
func (m *StartTwoPhasedPartnerTransferV1RequestBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1RequestBody.Marshal(b, m, deterministic)
}
func (m *StartTwoPhasedPartnerTransferV1RequestBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartTwoPhasedPartnerTransferV1RequestBody.Merge(m, src)
}
func (m *StartTwoPhasedPartnerTransferV1RequestBody) XXX_Size() int {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1RequestBody.Size(m)
}
func (m *StartTwoPhasedPartnerTransferV1RequestBody) XXX_DiscardUnknown() {
	xxx_messageInfo_StartTwoPhasedPartnerTransferV1RequestBody.DiscardUnknown(m)
}

var xxx_messageInfo_StartTwoPhasedPartnerTransferV1RequestBody proto.InternalMessageInfo

func (m *StartTwoPhasedPartnerTransferV1RequestBody) GetTransferKey() string {
	if m != nil {
		return m.TransferKey
	}
	return ""
}

func (m *StartTwoPhasedPartnerTransferV1RequestBody) GetFromAccountId() int64 {
	if m != nil {
		return m.FromAccountId
	}
	return 0
}

func (m *StartTwoPhasedPartnerTransferV1RequestBody) GetToAccountId() int64 {
	if m != nil {
		return m.ToAccountId
	}
	return 0
}

func (m *StartTwoPhasedPartnerTransferV1RequestBody) GetAmount() int64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

type StartTwoPhasedPartnerTransferV1Response struct {
	OperationId          int64    `protobuf:"varint,1,opt,name=operationId,proto3" json:"operationId"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StartTwoPhasedPartnerTransferV1Response) Reset() {
	*m = StartTwoPhasedPartnerTransferV1Response{}
}
func (m *StartTwoPhasedPartnerTransferV1Response) String() string { return proto.CompactTextString(m) }
func (*StartTwoPhasedPartnerTransferV1Response) ProtoMessage()    {}
func (*StartTwoPhasedPartnerTransferV1Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_c00dd1a11d825a35, []int{2}
}

func (m *StartTwoPhasedPartnerTransferV1Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1Response.Unmarshal(m, b)
}
func (m *StartTwoPhasedPartnerTransferV1Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1Response.Marshal(b, m, deterministic)
}
func (m *StartTwoPhasedPartnerTransferV1Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StartTwoPhasedPartnerTransferV1Response.Merge(m, src)
}
func (m *StartTwoPhasedPartnerTransferV1Response) XXX_Size() int {
	return xxx_messageInfo_StartTwoPhasedPartnerTransferV1Response.Size(m)
}
func (m *StartTwoPhasedPartnerTransferV1Response) XXX_DiscardUnknown() {
	xxx_messageInfo_StartTwoPhasedPartnerTransferV1Response.DiscardUnknown(m)
}

var xxx_messageInfo_StartTwoPhasedPartnerTransferV1Response proto.InternalMessageInfo

func (m *StartTwoPhasedPartnerTransferV1Response) GetOperationId() int64 {
	if m != nil {
		return m.OperationId
	}
	return 0
}

type CompleteTwoPhasedPartnerTransferV1Request struct {
	Token                string                                         `protobuf:"bytes,1,opt,name=token,proto3" json:"token"`
	Client               *client.Client                                 `protobuf:"bytes,2,opt,name=client,proto3" json:"client"`
	Body                 *CompleteTwoPhasedPartnerTransferV1RequestBody `protobuf:"bytes,3,opt,name=body,proto3" json:"body"`
	XXX_NoUnkeyedLiteral struct{}                                       `json:"-"`
	XXX_unrecognized     []byte                                         `json:"-"`
	XXX_sizecache        int32                                          `json:"-"`
}

func (m *CompleteTwoPhasedPartnerTransferV1Request) Reset() {
	*m = CompleteTwoPhasedPartnerTransferV1Request{}
}
func (m *CompleteTwoPhasedPartnerTransferV1Request) String() string {
	return proto.CompactTextString(m)
}
func (*CompleteTwoPhasedPartnerTransferV1Request) ProtoMessage() {}
func (*CompleteTwoPhasedPartnerTransferV1Request) Descriptor() ([]byte, []int) {
	return fileDescriptor_c00dd1a11d825a35, []int{3}
}

func (m *CompleteTwoPhasedPartnerTransferV1Request) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Request.Unmarshal(m, b)
}
func (m *CompleteTwoPhasedPartnerTransferV1Request) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Request.Marshal(b, m, deterministic)
}
func (m *CompleteTwoPhasedPartnerTransferV1Request) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Request.Merge(m, src)
}
func (m *CompleteTwoPhasedPartnerTransferV1Request) XXX_Size() int {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Request.Size(m)
}
func (m *CompleteTwoPhasedPartnerTransferV1Request) XXX_DiscardUnknown() {
	xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Request.DiscardUnknown(m)
}

var xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Request proto.InternalMessageInfo

func (m *CompleteTwoPhasedPartnerTransferV1Request) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *CompleteTwoPhasedPartnerTransferV1Request) GetClient() *client.Client {
	if m != nil {
		return m.Client
	}
	return nil
}

func (m *CompleteTwoPhasedPartnerTransferV1Request) GetBody() *CompleteTwoPhasedPartnerTransferV1RequestBody {
	if m != nil {
		return m.Body
	}
	return nil
}

type CompleteTwoPhasedPartnerTransferV1RequestBody struct {
	OperationId          int64    `protobuf:"varint,1,opt,name=operationId,proto3" json:"operationId"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) Reset() {
	*m = CompleteTwoPhasedPartnerTransferV1RequestBody{}
}
func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) String() string {
	return proto.CompactTextString(m)
}
func (*CompleteTwoPhasedPartnerTransferV1RequestBody) ProtoMessage() {}
func (*CompleteTwoPhasedPartnerTransferV1RequestBody) Descriptor() ([]byte, []int) {
	return fileDescriptor_c00dd1a11d825a35, []int{4}
}

func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1RequestBody.Unmarshal(m, b)
}
func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1RequestBody.Marshal(b, m, deterministic)
}
func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1RequestBody.Merge(m, src)
}
func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) XXX_Size() int {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1RequestBody.Size(m)
}
func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) XXX_DiscardUnknown() {
	xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1RequestBody.DiscardUnknown(m)
}

var xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1RequestBody proto.InternalMessageInfo

func (m *CompleteTwoPhasedPartnerTransferV1RequestBody) GetOperationId() int64 {
	if m != nil {
		return m.OperationId
	}
	return 0
}

type CompleteTwoPhasedPartnerTransferV1Response struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CompleteTwoPhasedPartnerTransferV1Response) Reset() {
	*m = CompleteTwoPhasedPartnerTransferV1Response{}
}
func (m *CompleteTwoPhasedPartnerTransferV1Response) String() string {
	return proto.CompactTextString(m)
}
func (*CompleteTwoPhasedPartnerTransferV1Response) ProtoMessage() {}
func (*CompleteTwoPhasedPartnerTransferV1Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_c00dd1a11d825a35, []int{5}
}

func (m *CompleteTwoPhasedPartnerTransferV1Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Response.Unmarshal(m, b)
}
func (m *CompleteTwoPhasedPartnerTransferV1Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Response.Marshal(b, m, deterministic)
}
func (m *CompleteTwoPhasedPartnerTransferV1Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Response.Merge(m, src)
}
func (m *CompleteTwoPhasedPartnerTransferV1Response) XXX_Size() int {
	return xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Response.Size(m)
}
func (m *CompleteTwoPhasedPartnerTransferV1Response) XXX_DiscardUnknown() {
	xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Response.DiscardUnknown(m)
}

var xxx_messageInfo_CompleteTwoPhasedPartnerTransferV1Response proto.InternalMessageInfo

func init() {
	proto.RegisterType((*StartTwoPhasedPartnerTransferV1Request)(nil), "partner_transfers_defs.StartTwoPhasedPartnerTransferV1Request")
	proto.RegisterType((*StartTwoPhasedPartnerTransferV1RequestBody)(nil), "partner_transfers_defs.StartTwoPhasedPartnerTransferV1RequestBody")
	proto.RegisterType((*StartTwoPhasedPartnerTransferV1Response)(nil), "partner_transfers_defs.StartTwoPhasedPartnerTransferV1Response")
	proto.RegisterType((*CompleteTwoPhasedPartnerTransferV1Request)(nil), "partner_transfers_defs.CompleteTwoPhasedPartnerTransferV1Request")
	proto.RegisterType((*CompleteTwoPhasedPartnerTransferV1RequestBody)(nil), "partner_transfers_defs.CompleteTwoPhasedPartnerTransferV1RequestBody")
	proto.RegisterType((*CompleteTwoPhasedPartnerTransferV1Response)(nil), "partner_transfers_defs.CompleteTwoPhasedPartnerTransferV1Response")
}

func init() {
	proto.RegisterFile("proto/common_partner-transfers.rmq.client.proto", fileDescriptor_c00dd1a11d825a35)
}

var fileDescriptor_c00dd1a11d825a35 = []byte{
	// 497 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xb4, 0x54, 0xd1, 0x6a, 0x13, 0x41,
	0x14, 0x65, 0xdd, 0x1a, 0xe8, 0x14, 0x7d, 0x18, 0xa4, 0x84, 0xbc, 0x18, 0x16, 0xa9, 0xb5, 0xb8,
	0xb3, 0xa6, 0x15, 0xa5, 0x28, 0x91, 0xa6, 0xf8, 0x50, 0x8a, 0xd0, 0xae, 0xa5, 0xa0, 0x2f, 0x61,
	0x76, 0xe7, 0x6e, 0x5c, 0x9a, 0x9d, 0xbb, 0x99, 0x99, 0x44, 0xf2, 0x53, 0x7e, 0x82, 0x4f, 0x3e,
	0xf5, 0x07, 0x7c, 0xf2, 0x03, 0xfc, 0x0b, 0xc9, 0xec, 0x56, 0x37, 0x24, 0x9a, 0xd5, 0xe0, 0x53,
	0x32, 0x77, 0xce, 0x3d, 0xf7, 0x9c, 0x39, 0xcb, 0x25, 0x41, 0xae, 0xd0, 0x60, 0x10, 0x63, 0x96,
	0xa1, 0xec, 0xe7, 0x5c, 0x19, 0x09, 0xca, 0x37, 0x8a, 0x4b, 0x9d, 0x80, 0xd2, 0x4c, 0x65, 0x23,
	0x16, 0x0f, 0x53, 0x90, 0x86, 0x59, 0x24, 0xdd, 0x2e, 0x31, 0xfd, 0x9f, 0x98, 0xbe, 0x80, 0x44,
	0xb7, 0x0e, 0x06, 0xa9, 0x19, 0xf2, 0x88, 0x25, 0x91, 0xf6, 0x05, 0x8b, 0x31, 0x0b, 0x04, 0x4c,
	0x02, 0x9e, 0x8d, 0x72, 0x5f, 0xe5, 0xb1, 0x3f, 0x00, 0x09, 0x8a, 0x1b, 0x54, 0xc1, 0x68, 0x0c,
	0x63, 0x28, 0xc8, 0x5a, 0x4f, 0x96, 0x34, 0x25, 0xa9, 0x4c, 0x4d, 0x8a, 0x52, 0x07, 0xc5, 0xdc,
	0xa0, 0x3a, 0xde, 0xfb, 0xec, 0x90, 0x9d, 0xb7, 0x86, 0x2b, 0x73, 0xf1, 0x11, 0xcf, 0x3e, 0x70,
	0x0d, 0xe2, 0xac, 0xd0, 0x73, 0x51, 0xca, 0xb9, 0xec, 0x84, 0x30, 0x1a, 0x83, 0x36, 0xf4, 0x1e,
	0xb9, 0x6d, 0xf0, 0x0a, 0x64, 0xd3, 0x69, 0x3b, 0xbb, 0x9b, 0x61, 0x71, 0xa0, 0x3b, 0xa4, 0x51,
	0x10, 0x36, 0x6f, 0xb5, 0x9d, 0xdd, 0xad, 0xfd, 0xbb, 0x37, 0xf6, 0x8e, 0xed, 0x4f, 0x58, 0xde,
	0xd2, 0x4b, 0xb2, 0x11, 0xa1, 0x98, 0x36, 0x5d, 0x8b, 0xea, 0xb1, 0xe5, 0xb6, 0x59, 0x3d, 0x2d,
	0x3d, 0x14, 0xd3, 0xd0, 0xf2, 0x79, 0x9f, 0x1c, 0xb2, 0x57, 0xbf, 0x89, 0xb6, 0xc9, 0xd6, 0xcd,
	0xc4, 0x53, 0x98, 0x96, 0x56, 0xaa, 0x25, 0xfa, 0x80, 0xdc, 0x49, 0x14, 0x66, 0x47, 0x71, 0x8c,
	0x63, 0x69, 0x4e, 0x84, 0xf5, 0xe5, 0x86, 0xf3, 0x45, 0xcb, 0x83, 0xbf, 0x30, 0xae, 0xc5, 0x54,
	0x4b, 0x74, 0x9b, 0x34, 0x78, 0x36, 0xfb, 0xdf, 0xdc, 0xb0, 0x97, 0xe5, 0xc9, 0x3b, 0x25, 0x0f,
	0x57, 0xea, 0xd5, 0x39, 0x4a, 0x0d, 0xb3, 0x21, 0x98, 0xcf, 0x62, 0x4e, 0x51, 0x9e, 0x08, 0x2b,
	0xd6, 0x0d, 0xab, 0x25, 0xef, 0x8b, 0x43, 0x1e, 0x1d, 0x63, 0x96, 0x0f, 0xc1, 0xc0, 0xff, 0x4e,
	0xf0, 0xdd, 0x5c, 0x82, 0xaf, 0x7f, 0x97, 0x60, 0x6d, 0x39, 0x95, 0x10, 0xcf, 0x89, 0xff, 0x57,
	0x6d, 0x35, 0x5e, 0xe6, 0x31, 0xd9, 0xab, 0x43, 0x59, 0xbc, 0xf4, 0xfe, 0x57, 0x97, 0x6c, 0x86,
	0x6f, 0xce, 0x0b, 0xc7, 0xf4, 0x9b, 0x43, 0xee, 0xaf, 0xc8, 0x88, 0x76, 0xd7, 0xfb, 0x82, 0x5b,
	0xaf, 0xfe, 0xb9, 0xbf, 0x90, 0xec, 0x1d, 0x5e, 0x77, 0x9f, 0xd1, 0xa7, 0x8a, 0x2d, 0xee, 0x98,
	0x3f, 0x32, 0xb0, 0x49, 0x87, 0x7e, 0x77, 0x88, 0xb7, 0xfa, 0x71, 0xe8, 0xd1, 0xda, 0x11, 0xb7,
	0x7a, 0xeb, 0x50, 0x94, 0x46, 0x5f, 0x5e, 0x77, 0x0f, 0xe9, 0xf3, 0x65, 0x46, 0x57, 0x91, 0xb0,
	0x49, 0xa7, 0x17, 0xbf, 0xe7, 0x0b, 0x4b, 0x31, 0xe2, 0xf1, 0x15, 0x48, 0x11, 0x0c, 0x61, 0xc0,
	0xe3, 0x69, 0xb9, 0xac, 0x83, 0x05, 0xfe, 0xb9, 0xe5, 0x69, 0xd7, 0xe5, 0x8b, 0xe5, 0x6e, 0xa2,
	0x86, 0xbd, 0x3d, 0xf8, 0x11, 0x00, 0x00, 0xff, 0xff, 0x6a, 0xe1, 0x3e, 0x10, 0xfe, 0x05, 0x00,
	0x00,
}