// Code generated by protoc-gen-go. DO NOT EDIT.
// source: main.proto

package registrygrpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type File struct {
	Name                 string   `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Data                 []byte   `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *File) Reset()         { *m = File{} }
func (m *File) String() string { return proto.CompactTextString(m) }
func (*File) ProtoMessage()    {}
func (*File) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{0}
}

func (m *File) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_File.Unmarshal(m, b)
}
func (m *File) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_File.Marshal(b, m, deterministic)
}
func (m *File) XXX_Merge(src proto.Message) {
	xxx_messageInfo_File.Merge(m, src)
}
func (m *File) XXX_Size() int {
	return xxx_messageInfo_File.Size(m)
}
func (m *File) XXX_DiscardUnknown() {
	xxx_messageInfo_File.DiscardUnknown(m)
}

var xxx_messageInfo_File proto.InternalMessageInfo

func (m *File) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *File) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type SaveFilesRequest struct {
	Cluster              string   `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	Channel              string   `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	Files                []*File  `protobuf:"bytes,3,rep,name=files,proto3" json:"files,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SaveFilesRequest) Reset()         { *m = SaveFilesRequest{} }
func (m *SaveFilesRequest) String() string { return proto.CompactTextString(m) }
func (*SaveFilesRequest) ProtoMessage()    {}
func (*SaveFilesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{1}
}

func (m *SaveFilesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveFilesRequest.Unmarshal(m, b)
}
func (m *SaveFilesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveFilesRequest.Marshal(b, m, deterministic)
}
func (m *SaveFilesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveFilesRequest.Merge(m, src)
}
func (m *SaveFilesRequest) XXX_Size() int {
	return xxx_messageInfo_SaveFilesRequest.Size(m)
}
func (m *SaveFilesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveFilesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SaveFilesRequest proto.InternalMessageInfo

func (m *SaveFilesRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *SaveFilesRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *SaveFilesRequest) GetFiles() []*File {
	if m != nil {
		return m.Files
	}
	return nil
}

type SaveFilesReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SaveFilesReply) Reset()         { *m = SaveFilesReply{} }
func (m *SaveFilesReply) String() string { return proto.CompactTextString(m) }
func (*SaveFilesReply) ProtoMessage()    {}
func (*SaveFilesReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{2}
}

func (m *SaveFilesReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SaveFilesReply.Unmarshal(m, b)
}
func (m *SaveFilesReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SaveFilesReply.Marshal(b, m, deterministic)
}
func (m *SaveFilesReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SaveFilesReply.Merge(m, src)
}
func (m *SaveFilesReply) XXX_Size() int {
	return xxx_messageInfo_SaveFilesReply.Size(m)
}
func (m *SaveFilesReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SaveFilesReply.DiscardUnknown(m)
}

var xxx_messageInfo_SaveFilesReply proto.InternalMessageInfo

type GetFilesRequest struct {
	Cluster              string   `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	Channel              string   `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFilesRequest) Reset()         { *m = GetFilesRequest{} }
func (m *GetFilesRequest) String() string { return proto.CompactTextString(m) }
func (*GetFilesRequest) ProtoMessage()    {}
func (*GetFilesRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{3}
}

func (m *GetFilesRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFilesRequest.Unmarshal(m, b)
}
func (m *GetFilesRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFilesRequest.Marshal(b, m, deterministic)
}
func (m *GetFilesRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFilesRequest.Merge(m, src)
}
func (m *GetFilesRequest) XXX_Size() int {
	return xxx_messageInfo_GetFilesRequest.Size(m)
}
func (m *GetFilesRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFilesRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetFilesRequest proto.InternalMessageInfo

func (m *GetFilesRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *GetFilesRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

type GetFilesReply struct {
	Files                []*File  `protobuf:"bytes,1,rep,name=files,proto3" json:"files,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetFilesReply) Reset()         { *m = GetFilesReply{} }
func (m *GetFilesReply) String() string { return proto.CompactTextString(m) }
func (*GetFilesReply) ProtoMessage()    {}
func (*GetFilesReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{4}
}

func (m *GetFilesReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetFilesReply.Unmarshal(m, b)
}
func (m *GetFilesReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetFilesReply.Marshal(b, m, deterministic)
}
func (m *GetFilesReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetFilesReply.Merge(m, src)
}
func (m *GetFilesReply) XXX_Size() int {
	return xxx_messageInfo_GetFilesReply.Size(m)
}
func (m *GetFilesReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetFilesReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetFilesReply proto.InternalMessageInfo

func (m *GetFilesReply) GetFiles() []*File {
	if m != nil {
		return m.Files
	}
	return nil
}

type SetMessageRequest struct {
	Cluster              string   `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	Channel              string   `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	Name                 string   `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetMessageRequest) Reset()         { *m = SetMessageRequest{} }
func (m *SetMessageRequest) String() string { return proto.CompactTextString(m) }
func (*SetMessageRequest) ProtoMessage()    {}
func (*SetMessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{5}
}

func (m *SetMessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetMessageRequest.Unmarshal(m, b)
}
func (m *SetMessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetMessageRequest.Marshal(b, m, deterministic)
}
func (m *SetMessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetMessageRequest.Merge(m, src)
}
func (m *SetMessageRequest) XXX_Size() int {
	return xxx_messageInfo_SetMessageRequest.Size(m)
}
func (m *SetMessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_SetMessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_SetMessageRequest proto.InternalMessageInfo

func (m *SetMessageRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *SetMessageRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

func (m *SetMessageRequest) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

type SetMessageReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SetMessageReply) Reset()         { *m = SetMessageReply{} }
func (m *SetMessageReply) String() string { return proto.CompactTextString(m) }
func (*SetMessageReply) ProtoMessage()    {}
func (*SetMessageReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{6}
}

func (m *SetMessageReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SetMessageReply.Unmarshal(m, b)
}
func (m *SetMessageReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SetMessageReply.Marshal(b, m, deterministic)
}
func (m *SetMessageReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SetMessageReply.Merge(m, src)
}
func (m *SetMessageReply) XXX_Size() int {
	return xxx_messageInfo_SetMessageReply.Size(m)
}
func (m *SetMessageReply) XXX_DiscardUnknown() {
	xxx_messageInfo_SetMessageReply.DiscardUnknown(m)
}

var xxx_messageInfo_SetMessageReply proto.InternalMessageInfo

type GetMessageRequest struct {
	Cluster              string   `protobuf:"bytes,1,opt,name=cluster,proto3" json:"cluster,omitempty"`
	Channel              string   `protobuf:"bytes,2,opt,name=channel,proto3" json:"channel,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMessageRequest) Reset()         { *m = GetMessageRequest{} }
func (m *GetMessageRequest) String() string { return proto.CompactTextString(m) }
func (*GetMessageRequest) ProtoMessage()    {}
func (*GetMessageRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{7}
}

func (m *GetMessageRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMessageRequest.Unmarshal(m, b)
}
func (m *GetMessageRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMessageRequest.Marshal(b, m, deterministic)
}
func (m *GetMessageRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMessageRequest.Merge(m, src)
}
func (m *GetMessageRequest) XXX_Size() int {
	return xxx_messageInfo_GetMessageRequest.Size(m)
}
func (m *GetMessageRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMessageRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetMessageRequest proto.InternalMessageInfo

func (m *GetMessageRequest) GetCluster() string {
	if m != nil {
		return m.Cluster
	}
	return ""
}

func (m *GetMessageRequest) GetChannel() string {
	if m != nil {
		return m.Channel
	}
	return ""
}

type GetMessageReply struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetMessageReply) Reset()         { *m = GetMessageReply{} }
func (m *GetMessageReply) String() string { return proto.CompactTextString(m) }
func (*GetMessageReply) ProtoMessage()    {}
func (*GetMessageReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_7ed94b0a22d11796, []int{8}
}

func (m *GetMessageReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetMessageReply.Unmarshal(m, b)
}
func (m *GetMessageReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetMessageReply.Marshal(b, m, deterministic)
}
func (m *GetMessageReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetMessageReply.Merge(m, src)
}
func (m *GetMessageReply) XXX_Size() int {
	return xxx_messageInfo_GetMessageReply.Size(m)
}
func (m *GetMessageReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetMessageReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetMessageReply proto.InternalMessageInfo

func (m *GetMessageReply) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func init() {
	proto.RegisterType((*File)(nil), "registrygrpc.File")
	proto.RegisterType((*SaveFilesRequest)(nil), "registrygrpc.SaveFilesRequest")
	proto.RegisterType((*SaveFilesReply)(nil), "registrygrpc.SaveFilesReply")
	proto.RegisterType((*GetFilesRequest)(nil), "registrygrpc.GetFilesRequest")
	proto.RegisterType((*GetFilesReply)(nil), "registrygrpc.GetFilesReply")
	proto.RegisterType((*SetMessageRequest)(nil), "registrygrpc.SetMessageRequest")
	proto.RegisterType((*SetMessageReply)(nil), "registrygrpc.SetMessageReply")
	proto.RegisterType((*GetMessageRequest)(nil), "registrygrpc.GetMessageRequest")
	proto.RegisterType((*GetMessageReply)(nil), "registrygrpc.GetMessageReply")
}

func init() { proto.RegisterFile("main.proto", fileDescriptor_7ed94b0a22d11796) }

var fileDescriptor_7ed94b0a22d11796 = []byte{
	// 323 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xac, 0x93, 0xcd, 0x4e, 0xeb, 0x30,
	0x10, 0x46, 0x6f, 0x9a, 0x5e, 0xa0, 0x43, 0x4b, 0xdb, 0x59, 0x45, 0x85, 0x42, 0x64, 0x09, 0x29,
	0xab, 0x2c, 0xca, 0x8a, 0x07, 0x00, 0x23, 0x21, 0x10, 0x4a, 0x97, 0xac, 0x4c, 0x19, 0x42, 0xa4,
	0x34, 0x0d, 0xb1, 0x8b, 0x94, 0x17, 0xe4, 0xb9, 0x90, 0xdd, 0xe6, 0xa7, 0x91, 0xa2, 0x2e, 0xca,
	0xce, 0x8e, 0xc7, 0x67, 0x3e, 0x1f, 0x3b, 0x00, 0x4b, 0x11, 0x25, 0x7e, 0x9a, 0xad, 0xd4, 0x0a,
	0xfb, 0x19, 0x85, 0x91, 0x54, 0x59, 0x1e, 0x66, 0xe9, 0x82, 0xf9, 0xd0, 0xbd, 0x8f, 0x62, 0x42,
	0x84, 0x6e, 0x22, 0x96, 0xe4, 0x74, 0x5c, 0xcb, 0xeb, 0x05, 0x66, 0xac, 0xbf, 0xbd, 0x0b, 0x25,
	0x1c, 0xdb, 0xb5, 0xbc, 0x7e, 0x60, 0xc6, 0x2c, 0x85, 0xd1, 0x5c, 0x7c, 0x93, 0xde, 0x23, 0x03,
	0xfa, 0x5a, 0x93, 0x54, 0xe8, 0xc0, 0xf1, 0x22, 0x5e, 0x4b, 0x45, 0x99, 0x63, 0x99, 0xed, 0xc5,
	0xd4, 0xac, 0x7c, 0x8a, 0x24, 0xa1, 0x78, 0x0b, 0x2e, 0xa6, 0xe8, 0xc1, 0xff, 0x0f, 0xcd, 0x70,
	0x6c, 0xd7, 0xf6, 0x4e, 0x67, 0xe8, 0xd7, 0x53, 0xf9, 0x1a, 0x1f, 0x6c, 0x0a, 0xd8, 0x08, 0xce,
	0x6a, 0x1d, 0xd3, 0x38, 0x67, 0x77, 0x30, 0xe4, 0xa4, 0x0e, 0x8d, 0xc0, 0x6e, 0x61, 0x50, 0x61,
	0xd2, 0x38, 0xaf, 0x32, 0x59, 0xfb, 0x32, 0xbd, 0xc2, 0x78, 0x4e, 0xea, 0x89, 0xa4, 0x14, 0x21,
	0x1d, 0xa2, 0xa1, 0xd0, 0x6e, 0x57, 0xda, 0xd9, 0x18, 0x86, 0x75, 0xb8, 0x3e, 0x31, 0x87, 0x31,
	0xff, 0x8b, 0x7e, 0xec, 0xda, 0xa8, 0xab, 0xb3, 0xcb, 0x08, 0x56, 0x15, 0x61, 0xf6, 0xd3, 0x81,
	0xc1, 0x8b, 0x7e, 0x2d, 0xc1, 0xd6, 0x00, 0x3e, 0x42, 0xaf, 0xbc, 0x05, 0xbc, 0xdc, 0x35, 0xd3,
	0x7c, 0x10, 0x93, 0x8b, 0xd6, 0x75, 0x7d, 0x98, 0x7f, 0xf8, 0x00, 0x27, 0x85, 0x79, 0x9c, 0xee,
	0xd6, 0x36, 0x2e, 0x76, 0x72, 0xde, 0xb6, 0xbc, 0x21, 0x3d, 0x03, 0x54, 0xae, 0xf0, 0xaa, 0xd1,
	0xb7, 0xa9, 0x6c, 0x32, 0x6d, 0x2f, 0x28, 0x79, 0xbc, 0x95, 0xc7, 0xf7, 0xf1, 0x78, 0x93, 0xf7,
	0x76, 0x64, 0xfe, 0xb9, 0x9b, 0xdf, 0x00, 0x00, 0x00, 0xff, 0xff, 0xb5, 0x72, 0x30, 0x40, 0x81,
	0x03, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// ProtoRegistryClient is the client API for ProtoRegistry service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ProtoRegistryClient interface {
	SaveFiles(ctx context.Context, in *SaveFilesRequest, opts ...grpc.CallOption) (*SaveFilesReply, error)
	GetFiles(ctx context.Context, in *GetFilesRequest, opts ...grpc.CallOption) (*GetFilesReply, error)
	SetMessage(ctx context.Context, in *SetMessageRequest, opts ...grpc.CallOption) (*SetMessageReply, error)
	GetMessage(ctx context.Context, in *GetMessageRequest, opts ...grpc.CallOption) (*GetMessageReply, error)
}

type protoRegistryClient struct {
	cc *grpc.ClientConn
}

func NewProtoRegistryClient(cc *grpc.ClientConn) ProtoRegistryClient {
	return &protoRegistryClient{cc}
}

func (c *protoRegistryClient) SaveFiles(ctx context.Context, in *SaveFilesRequest, opts ...grpc.CallOption) (*SaveFilesReply, error) {
	out := new(SaveFilesReply)
	err := c.cc.Invoke(ctx, "/registrygrpc.ProtoRegistry/SaveFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protoRegistryClient) GetFiles(ctx context.Context, in *GetFilesRequest, opts ...grpc.CallOption) (*GetFilesReply, error) {
	out := new(GetFilesReply)
	err := c.cc.Invoke(ctx, "/registrygrpc.ProtoRegistry/GetFiles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protoRegistryClient) SetMessage(ctx context.Context, in *SetMessageRequest, opts ...grpc.CallOption) (*SetMessageReply, error) {
	out := new(SetMessageReply)
	err := c.cc.Invoke(ctx, "/registrygrpc.ProtoRegistry/SetMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *protoRegistryClient) GetMessage(ctx context.Context, in *GetMessageRequest, opts ...grpc.CallOption) (*GetMessageReply, error) {
	out := new(GetMessageReply)
	err := c.cc.Invoke(ctx, "/registrygrpc.ProtoRegistry/GetMessage", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProtoRegistryServer is the server API for ProtoRegistry service.
type ProtoRegistryServer interface {
	SaveFiles(context.Context, *SaveFilesRequest) (*SaveFilesReply, error)
	GetFiles(context.Context, *GetFilesRequest) (*GetFilesReply, error)
	SetMessage(context.Context, *SetMessageRequest) (*SetMessageReply, error)
	GetMessage(context.Context, *GetMessageRequest) (*GetMessageReply, error)
}

func RegisterProtoRegistryServer(s *grpc.Server, srv ProtoRegistryServer) {
	s.RegisterService(&_ProtoRegistry_serviceDesc, srv)
}

func _ProtoRegistry_SaveFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SaveFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoRegistryServer).SaveFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registrygrpc.ProtoRegistry/SaveFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoRegistryServer).SaveFiles(ctx, req.(*SaveFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProtoRegistry_GetFiles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetFilesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoRegistryServer).GetFiles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registrygrpc.ProtoRegistry/GetFiles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoRegistryServer).GetFiles(ctx, req.(*GetFilesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProtoRegistry_SetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoRegistryServer).SetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registrygrpc.ProtoRegistry/SetMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoRegistryServer).SetMessage(ctx, req.(*SetMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ProtoRegistry_GetMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMessageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProtoRegistryServer).GetMessage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/registrygrpc.ProtoRegistry/GetMessage",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProtoRegistryServer).GetMessage(ctx, req.(*GetMessageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ProtoRegistry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "registrygrpc.ProtoRegistry",
	HandlerType: (*ProtoRegistryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SaveFiles",
			Handler:    _ProtoRegistry_SaveFiles_Handler,
		},
		{
			MethodName: "GetFiles",
			Handler:    _ProtoRegistry_GetFiles_Handler,
		},
		{
			MethodName: "SetMessage",
			Handler:    _ProtoRegistry_SetMessage_Handler,
		},
		{
			MethodName: "GetMessage",
			Handler:    _ProtoRegistry_GetMessage_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "main.proto",
}
