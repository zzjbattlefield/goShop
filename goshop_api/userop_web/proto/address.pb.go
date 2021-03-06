// Code generated by protoc-gen-go. DO NOT EDIT.
// source: address.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
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

type AddressRequest struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId               int32    `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Province             string   `protobuf:"bytes,3,opt,name=province,proto3" json:"province,omitempty"`
	City                 string   `protobuf:"bytes,4,opt,name=city,proto3" json:"city,omitempty"`
	District             string   `protobuf:"bytes,5,opt,name=district,proto3" json:"district,omitempty"`
	Address              string   `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
	SignerName           string   `protobuf:"bytes,7,opt,name=signerName,proto3" json:"signerName,omitempty"`
	SignerMobile         string   `protobuf:"bytes,8,opt,name=signerMobile,proto3" json:"signerMobile,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddressRequest) Reset()         { *m = AddressRequest{} }
func (m *AddressRequest) String() string { return proto.CompactTextString(m) }
func (*AddressRequest) ProtoMessage()    {}
func (*AddressRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_982c640dad8fe78e, []int{0}
}

func (m *AddressRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddressRequest.Unmarshal(m, b)
}
func (m *AddressRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddressRequest.Marshal(b, m, deterministic)
}
func (m *AddressRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddressRequest.Merge(m, src)
}
func (m *AddressRequest) XXX_Size() int {
	return xxx_messageInfo_AddressRequest.Size(m)
}
func (m *AddressRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AddressRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AddressRequest proto.InternalMessageInfo

func (m *AddressRequest) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *AddressRequest) GetUserId() int32 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *AddressRequest) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *AddressRequest) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *AddressRequest) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *AddressRequest) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AddressRequest) GetSignerName() string {
	if m != nil {
		return m.SignerName
	}
	return ""
}

func (m *AddressRequest) GetSignerMobile() string {
	if m != nil {
		return m.SignerMobile
	}
	return ""
}

type AddressResponse struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	UserId               int32    `protobuf:"varint,2,opt,name=userId,proto3" json:"userId,omitempty"`
	Province             string   `protobuf:"bytes,3,opt,name=province,proto3" json:"province,omitempty"`
	City                 string   `protobuf:"bytes,4,opt,name=city,proto3" json:"city,omitempty"`
	District             string   `protobuf:"bytes,5,opt,name=district,proto3" json:"district,omitempty"`
	Address              string   `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
	SignerName           string   `protobuf:"bytes,7,opt,name=signerName,proto3" json:"signerName,omitempty"`
	SignerMobile         string   `protobuf:"bytes,8,opt,name=signerMobile,proto3" json:"signerMobile,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AddressResponse) Reset()         { *m = AddressResponse{} }
func (m *AddressResponse) String() string { return proto.CompactTextString(m) }
func (*AddressResponse) ProtoMessage()    {}
func (*AddressResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_982c640dad8fe78e, []int{1}
}

func (m *AddressResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddressResponse.Unmarshal(m, b)
}
func (m *AddressResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddressResponse.Marshal(b, m, deterministic)
}
func (m *AddressResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddressResponse.Merge(m, src)
}
func (m *AddressResponse) XXX_Size() int {
	return xxx_messageInfo_AddressResponse.Size(m)
}
func (m *AddressResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddressResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddressResponse proto.InternalMessageInfo

func (m *AddressResponse) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *AddressResponse) GetUserId() int32 {
	if m != nil {
		return m.UserId
	}
	return 0
}

func (m *AddressResponse) GetProvince() string {
	if m != nil {
		return m.Province
	}
	return ""
}

func (m *AddressResponse) GetCity() string {
	if m != nil {
		return m.City
	}
	return ""
}

func (m *AddressResponse) GetDistrict() string {
	if m != nil {
		return m.District
	}
	return ""
}

func (m *AddressResponse) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *AddressResponse) GetSignerName() string {
	if m != nil {
		return m.SignerName
	}
	return ""
}

func (m *AddressResponse) GetSignerMobile() string {
	if m != nil {
		return m.SignerMobile
	}
	return ""
}

type AddressListResponse struct {
	Total                int32              `protobuf:"varint,1,opt,name=total,proto3" json:"total,omitempty"`
	Data                 []*AddressResponse `protobuf:"bytes,2,rep,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *AddressListResponse) Reset()         { *m = AddressListResponse{} }
func (m *AddressListResponse) String() string { return proto.CompactTextString(m) }
func (*AddressListResponse) ProtoMessage()    {}
func (*AddressListResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_982c640dad8fe78e, []int{2}
}

func (m *AddressListResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AddressListResponse.Unmarshal(m, b)
}
func (m *AddressListResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AddressListResponse.Marshal(b, m, deterministic)
}
func (m *AddressListResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AddressListResponse.Merge(m, src)
}
func (m *AddressListResponse) XXX_Size() int {
	return xxx_messageInfo_AddressListResponse.Size(m)
}
func (m *AddressListResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_AddressListResponse.DiscardUnknown(m)
}

var xxx_messageInfo_AddressListResponse proto.InternalMessageInfo

func (m *AddressListResponse) GetTotal() int32 {
	if m != nil {
		return m.Total
	}
	return 0
}

func (m *AddressListResponse) GetData() []*AddressResponse {
	if m != nil {
		return m.Data
	}
	return nil
}

func init() {
	proto.RegisterType((*AddressRequest)(nil), "AddressRequest")
	proto.RegisterType((*AddressResponse)(nil), "AddressResponse")
	proto.RegisterType((*AddressListResponse)(nil), "AddressListResponse")
}

func init() { proto.RegisterFile("address.proto", fileDescriptor_982c640dad8fe78e) }

var fileDescriptor_982c640dad8fe78e = []byte{
	// 337 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xdc, 0x52, 0x4d, 0x4b, 0xc3, 0x30,
	0x18, 0xa6, 0xdd, 0x47, 0xb7, 0x57, 0xb7, 0x49, 0x1c, 0x23, 0x4c, 0x90, 0x31, 0x3c, 0xec, 0x94,
	0xc1, 0x3c, 0x28, 0x78, 0xf2, 0x0b, 0x11, 0x54, 0xb0, 0xe0, 0xc5, 0x5b, 0xb6, 0xbc, 0x8e, 0x40,
	0xb7, 0xd4, 0x24, 0x13, 0xf6, 0x77, 0xf5, 0xe6, 0xaf, 0x90, 0x25, 0x6d, 0xd9, 0x9c, 0x08, 0x5e,
	0x3d, 0xb5, 0xcf, 0x57, 0x93, 0xa7, 0x3c, 0xd0, 0xe0, 0x42, 0x68, 0x34, 0x86, 0xa5, 0x5a, 0x59,
	0xd5, 0x3d, 0x98, 0x2a, 0x35, 0x4d, 0x70, 0xe8, 0xd0, 0x78, 0xf1, 0x32, 0xc4, 0x59, 0x6a, 0x97,
	0x5e, 0xec, 0xbf, 0x07, 0xd0, 0x3c, 0xf7, 0xf6, 0x18, 0x5f, 0x17, 0x68, 0x2c, 0x69, 0x42, 0x28,
	0x05, 0x0d, 0x7a, 0xc1, 0xa0, 0x12, 0x87, 0x52, 0x90, 0x0e, 0x54, 0x17, 0x06, 0xf5, 0xad, 0xa0,
	0xa1, 0xe3, 0x32, 0x44, 0xba, 0x50, 0x4b, 0xb5, 0x7a, 0x93, 0xf3, 0x09, 0xd2, 0x52, 0x2f, 0x18,
	0xd4, 0xe3, 0x02, 0x13, 0x02, 0xe5, 0x89, 0xb4, 0x4b, 0x5a, 0x76, 0xbc, 0x7b, 0x5f, 0xf9, 0x85,
	0x34, 0x56, 0xcb, 0x89, 0xa5, 0x15, 0xef, 0xcf, 0x31, 0xa1, 0x10, 0x65, 0x97, 0xa6, 0x55, 0x27,
	0xe5, 0x90, 0x1c, 0x02, 0x18, 0x39, 0x9d, 0xa3, 0x7e, 0xe0, 0x33, 0xa4, 0x91, 0x13, 0xd7, 0x18,
	0xd2, 0x87, 0x5d, 0x8f, 0xee, 0xd5, 0x58, 0x26, 0x48, 0x6b, 0xce, 0xb1, 0xc1, 0xf5, 0x3f, 0x02,
	0x68, 0x15, 0x25, 0x4d, 0xaa, 0xe6, 0x06, 0xff, 0x61, 0xcb, 0x47, 0xd8, 0xcf, 0x4a, 0xde, 0x49,
	0x63, 0x8b, 0xa2, 0x6d, 0xa8, 0x58, 0x65, 0x79, 0x92, 0x75, 0xf5, 0x80, 0x1c, 0x41, 0x59, 0x70,
	0xcb, 0x69, 0xd8, 0x2b, 0x0d, 0x76, 0x46, 0x7b, 0xec, 0xdb, 0xef, 0x89, 0x9d, 0x3a, 0xfa, 0x0c,
	0x20, 0xca, 0x14, 0x72, 0x02, 0xcd, 0x1b, 0xb4, 0x6b, 0x27, 0x90, 0x16, 0xdb, 0x5c, 0x4e, 0xb7,
	0xcd, 0x7e, 0xba, 0xc0, 0x08, 0x1a, 0x97, 0x1a, 0xb9, 0xc5, 0xfc, 0x4b, 0x5b, 0xb9, 0xad, 0xe3,
	0xc9, 0x29, 0x34, 0xae, 0x30, 0xc1, 0x5f, 0x32, 0x1d, 0xe6, 0x67, 0xcd, 0xf2, 0x59, 0xb3, 0xeb,
	0xd5, 0xac, 0x57, 0xc9, 0xa7, 0x54, 0xf0, 0xbf, 0x27, 0x2f, 0xea, 0xcf, 0x11, 0x3b, 0xf3, 0x5c,
	0xd5, 0x3d, 0x8e, 0xbf, 0x02, 0x00, 0x00, 0xff, 0xff, 0x34, 0x72, 0xc5, 0x22, 0x4a, 0x03, 0x00,
	0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AddressClient is the client API for Address service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AddressClient interface {
	GetAddressList(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*AddressListResponse, error)
	CreateAddress(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*AddressResponse, error)
	DeleteAddress(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*empty.Empty, error)
	UpdateAddress(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*empty.Empty, error)
}

type addressClient struct {
	cc grpc.ClientConnInterface
}

func NewAddressClient(cc grpc.ClientConnInterface) AddressClient {
	return &addressClient{cc}
}

func (c *addressClient) GetAddressList(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*AddressListResponse, error) {
	out := new(AddressListResponse)
	err := c.cc.Invoke(ctx, "/Address/GetAddressList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressClient) CreateAddress(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*AddressResponse, error) {
	out := new(AddressResponse)
	err := c.cc.Invoke(ctx, "/Address/CreateAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressClient) DeleteAddress(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/Address/DeleteAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *addressClient) UpdateAddress(ctx context.Context, in *AddressRequest, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/Address/UpdateAddress", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AddressServer is the server API for Address service.
type AddressServer interface {
	GetAddressList(context.Context, *AddressRequest) (*AddressListResponse, error)
	CreateAddress(context.Context, *AddressRequest) (*AddressResponse, error)
	DeleteAddress(context.Context, *AddressRequest) (*empty.Empty, error)
	UpdateAddress(context.Context, *AddressRequest) (*empty.Empty, error)
}

// UnimplementedAddressServer can be embedded to have forward compatible implementations.
type UnimplementedAddressServer struct {
}

func (*UnimplementedAddressServer) GetAddressList(ctx context.Context, req *AddressRequest) (*AddressListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAddressList not implemented")
}
func (*UnimplementedAddressServer) CreateAddress(ctx context.Context, req *AddressRequest) (*AddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateAddress not implemented")
}
func (*UnimplementedAddressServer) DeleteAddress(ctx context.Context, req *AddressRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAddress not implemented")
}
func (*UnimplementedAddressServer) UpdateAddress(ctx context.Context, req *AddressRequest) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAddress not implemented")
}

func RegisterAddressServer(s *grpc.Server, srv AddressServer) {
	s.RegisterService(&_Address_serviceDesc, srv)
}

func _Address_GetAddressList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddressServer).GetAddressList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Address/GetAddressList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddressServer).GetAddressList(ctx, req.(*AddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Address_CreateAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddressServer).CreateAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Address/CreateAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddressServer).CreateAddress(ctx, req.(*AddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Address_DeleteAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddressServer).DeleteAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Address/DeleteAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddressServer).DeleteAddress(ctx, req.(*AddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Address_UpdateAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AddressServer).UpdateAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Address/UpdateAddress",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AddressServer).UpdateAddress(ctx, req.(*AddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Address_serviceDesc = grpc.ServiceDesc{
	ServiceName: "Address",
	HandlerType: (*AddressServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAddressList",
			Handler:    _Address_GetAddressList_Handler,
		},
		{
			MethodName: "CreateAddress",
			Handler:    _Address_CreateAddress_Handler,
		},
		{
			MethodName: "DeleteAddress",
			Handler:    _Address_DeleteAddress_Handler,
		},
		{
			MethodName: "UpdateAddress",
			Handler:    _Address_UpdateAddress_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "address.proto",
}
