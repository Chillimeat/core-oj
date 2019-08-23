// Code generated by protoc-gen-go. DO NOT EDIT.
// source: compiler.proto

package grpc

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
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

type CompileResponseCode int32

const (
	CompileResponseCode_Ok                     CompileResponseCode = 0
	CompileResponseCode_FileNotFound           CompileResponseCode = 1
	CompileResponseCode_CompileTimeLimitExceed CompileResponseCode = 2
	CompileResponseCode_CompileError           CompileResponseCode = 3
)

var CompileResponseCode_name = map[int32]string{
	0: "Ok",
	1: "FileNotFound",
	2: "CompileTimeLimitExceed",
	3: "CompileError",
}

var CompileResponseCode_value = map[string]int32{
	"Ok":                     0,
	"FileNotFound":           1,
	"CompileTimeLimitExceed": 2,
	"CompileError":           3,
}

func (x CompileResponseCode) String() string {
	return proto.EnumName(CompileResponseCode_name, int32(x))
}

func (CompileResponseCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_6a5727dbaeb66833, []int{0}
}

type CompileRequest struct {
	CompilerType         string   `protobuf:"bytes,1,opt,name=compiler_type,json=compilerType,proto3" json:"compiler_type,omitempty"`
	CodePath             string   `protobuf:"bytes,2,opt,name=code_path,json=codePath,proto3" json:"code_path,omitempty"`
	AimPath              string   `protobuf:"bytes,3,opt,name=aim_path,json=aimPath,proto3" json:"aim_path,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CompileRequest) Reset()         { *m = CompileRequest{} }
func (m *CompileRequest) String() string { return proto.CompactTextString(m) }
func (*CompileRequest) ProtoMessage()    {}
func (*CompileRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a5727dbaeb66833, []int{0}
}

func (m *CompileRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompileRequest.Unmarshal(m, b)
}
func (m *CompileRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompileRequest.Marshal(b, m, deterministic)
}
func (m *CompileRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompileRequest.Merge(m, src)
}
func (m *CompileRequest) XXX_Size() int {
	return xxx_messageInfo_CompileRequest.Size(m)
}
func (m *CompileRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CompileRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CompileRequest proto.InternalMessageInfo

func (m *CompileRequest) GetCompilerType() string {
	if m != nil {
		return m.CompilerType
	}
	return ""
}

func (m *CompileRequest) GetCodePath() string {
	if m != nil {
		return m.CodePath
	}
	return ""
}

func (m *CompileRequest) GetAimPath() string {
	if m != nil {
		return m.AimPath
	}
	return ""
}

type CompileReply struct {
	ResponseCode         CompileResponseCode `protobuf:"varint,1,opt,name=response_code,json=responseCode,proto3,enum=compilerrpc.CompileResponseCode" json:"response_code,omitempty"`
	Info                 []byte              `protobuf:"bytes,2,opt,name=info,proto3" json:"info,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *CompileReply) Reset()         { *m = CompileReply{} }
func (m *CompileReply) String() string { return proto.CompactTextString(m) }
func (*CompileReply) ProtoMessage()    {}
func (*CompileReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a5727dbaeb66833, []int{1}
}

func (m *CompileReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompileReply.Unmarshal(m, b)
}
func (m *CompileReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompileReply.Marshal(b, m, deterministic)
}
func (m *CompileReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompileReply.Merge(m, src)
}
func (m *CompileReply) XXX_Size() int {
	return xxx_messageInfo_CompileReply.Size(m)
}
func (m *CompileReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CompileReply.DiscardUnknown(m)
}

var xxx_messageInfo_CompileReply proto.InternalMessageInfo

func (m *CompileReply) GetResponseCode() CompileResponseCode {
	if m != nil {
		return m.ResponseCode
	}
	return CompileResponseCode_Ok
}

func (m *CompileReply) GetInfo() []byte {
	if m != nil {
		return m.Info
	}
	return nil
}

type InfoRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *InfoRequest) Reset()         { *m = InfoRequest{} }
func (m *InfoRequest) String() string { return proto.CompactTextString(m) }
func (*InfoRequest) ProtoMessage()    {}
func (*InfoRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a5727dbaeb66833, []int{2}
}

func (m *InfoRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoRequest.Unmarshal(m, b)
}
func (m *InfoRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoRequest.Marshal(b, m, deterministic)
}
func (m *InfoRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoRequest.Merge(m, src)
}
func (m *InfoRequest) XXX_Size() int {
	return xxx_messageInfo_InfoRequest.Size(m)
}
func (m *InfoRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoRequest.DiscardUnknown(m)
}

var xxx_messageInfo_InfoRequest proto.InternalMessageInfo

type CompilerToolInfo struct {
	Name                 string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Path                 string   `protobuf:"bytes,2,opt,name=path,proto3" json:"path,omitempty"`
	Version              string   `protobuf:"bytes,3,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CompilerToolInfo) Reset()         { *m = CompilerToolInfo{} }
func (m *CompilerToolInfo) String() string { return proto.CompactTextString(m) }
func (*CompilerToolInfo) ProtoMessage()    {}
func (*CompilerToolInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a5727dbaeb66833, []int{3}
}

func (m *CompilerToolInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CompilerToolInfo.Unmarshal(m, b)
}
func (m *CompilerToolInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CompilerToolInfo.Marshal(b, m, deterministic)
}
func (m *CompilerToolInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CompilerToolInfo.Merge(m, src)
}
func (m *CompilerToolInfo) XXX_Size() int {
	return xxx_messageInfo_CompilerToolInfo.Size(m)
}
func (m *CompilerToolInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_CompilerToolInfo.DiscardUnknown(m)
}

var xxx_messageInfo_CompilerToolInfo proto.InternalMessageInfo

func (m *CompilerToolInfo) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *CompilerToolInfo) GetPath() string {
	if m != nil {
		return m.Path
	}
	return ""
}

func (m *CompilerToolInfo) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type InfoReply struct {
	CompilerVersion      string              `protobuf:"bytes,1,opt,name=compiler_version,json=compilerVersion,proto3" json:"compiler_version,omitempty"`
	CompilerTools        []*CompilerToolInfo `protobuf:"bytes,2,rep,name=compiler_tools,json=compilerTools,proto3" json:"compiler_tools,omitempty"`
	XXX_NoUnkeyedLiteral struct{}            `json:"-"`
	XXX_unrecognized     []byte              `json:"-"`
	XXX_sizecache        int32               `json:"-"`
}

func (m *InfoReply) Reset()         { *m = InfoReply{} }
func (m *InfoReply) String() string { return proto.CompactTextString(m) }
func (*InfoReply) ProtoMessage()    {}
func (*InfoReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_6a5727dbaeb66833, []int{4}
}

func (m *InfoReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_InfoReply.Unmarshal(m, b)
}
func (m *InfoReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_InfoReply.Marshal(b, m, deterministic)
}
func (m *InfoReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_InfoReply.Merge(m, src)
}
func (m *InfoReply) XXX_Size() int {
	return xxx_messageInfo_InfoReply.Size(m)
}
func (m *InfoReply) XXX_DiscardUnknown() {
	xxx_messageInfo_InfoReply.DiscardUnknown(m)
}

var xxx_messageInfo_InfoReply proto.InternalMessageInfo

func (m *InfoReply) GetCompilerVersion() string {
	if m != nil {
		return m.CompilerVersion
	}
	return ""
}

func (m *InfoReply) GetCompilerTools() []*CompilerToolInfo {
	if m != nil {
		return m.CompilerTools
	}
	return nil
}

func init() {
	proto.RegisterEnum("compilerrpc.CompileResponseCode", CompileResponseCode_name, CompileResponseCode_value)
	proto.RegisterType((*CompileRequest)(nil), "compilerrpc.CompileRequest")
	proto.RegisterType((*CompileReply)(nil), "compilerrpc.CompileReply")
	proto.RegisterType((*InfoRequest)(nil), "compilerrpc.InfoRequest")
	proto.RegisterType((*CompilerToolInfo)(nil), "compilerrpc.CompilerToolInfo")
	proto.RegisterType((*InfoReply)(nil), "compilerrpc.InfoReply")
}

func init() { proto.RegisterFile("compiler.proto", fileDescriptor_6a5727dbaeb66833) }

var fileDescriptor_6a5727dbaeb66833 = []byte{
	// 421 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x6c, 0x52, 0xdf, 0x8b, 0xd3, 0x40,
	0x10, 0xb6, 0xed, 0x71, 0x6d, 0xa7, 0x69, 0x0d, 0x2b, 0x1c, 0xb9, 0x1e, 0x42, 0x89, 0x2f, 0xa7,
	0x70, 0x0d, 0xd6, 0x17, 0x5f, 0xb5, 0xd7, 0x03, 0xc1, 0x5f, 0x84, 0xe2, 0x83, 0x20, 0x25, 0x97,
	0xcc, 0x5d, 0x57, 0xb3, 0x99, 0x75, 0xb3, 0x15, 0x03, 0xbe, 0xfb, 0x6f, 0xcb, 0x6e, 0xb2, 0xb1,
	0x07, 0x79, 0x9b, 0xf9, 0xbe, 0x6f, 0xb2, 0xf3, 0x7d, 0x13, 0x98, 0xa5, 0x24, 0x24, 0xcf, 0x51,
	0x2d, 0xa5, 0x22, 0x4d, 0x6c, 0xe2, 0x7a, 0x25, 0xd3, 0x50, 0xc0, 0x6c, 0x5d, 0xb7, 0x31, 0xfe,
	0x3c, 0x60, 0xa9, 0xd9, 0x33, 0x98, 0x3a, 0xc1, 0x4e, 0x57, 0x12, 0x83, 0xde, 0xa2, 0x77, 0x39,
	0x8e, 0x3d, 0x07, 0x6e, 0x2b, 0x89, 0xec, 0x02, 0xc6, 0x29, 0x65, 0xb8, 0x93, 0x89, 0xde, 0x07,
	0x7d, 0x2b, 0x18, 0x19, 0xe0, 0x73, 0xa2, 0xf7, 0xec, 0x1c, 0x46, 0x09, 0x17, 0x35, 0x37, 0xb0,
	0xdc, 0x30, 0xe1, 0xc2, 0x50, 0x21, 0x07, 0xaf, 0x7d, 0x4e, 0xe6, 0x15, 0xdb, 0xc0, 0x54, 0x61,
	0x29, 0xa9, 0x28, 0x71, 0x67, 0xe6, 0xed, 0x63, 0xb3, 0xd5, 0x62, 0x79, 0xb4, 0xe3, 0xb2, 0x9d,
	0xa8, 0x85, 0x6b, 0xca, 0x30, 0xf6, 0xd4, 0x51, 0xc7, 0x18, 0x9c, 0xf0, 0xe2, 0x8e, 0xec, 0x26,
	0x5e, 0x6c, 0xeb, 0x70, 0x0a, 0x93, 0x77, 0xc5, 0x1d, 0x35, 0xb6, 0xc2, 0x2d, 0xf8, 0x6b, 0xe7,
	0x80, 0x28, 0x37, 0x94, 0x19, 0x2b, 0x12, 0xe1, 0x1c, 0xda, 0xda, 0x60, 0x47, 0xa6, 0x6c, 0xcd,
	0x02, 0x18, 0xfe, 0x42, 0x55, 0x72, 0x2a, 0x9c, 0x9f, 0xa6, 0x0d, 0xff, 0xc0, 0xb8, 0x7e, 0xc4,
	0x98, 0x79, 0x0e, 0x7e, 0x9b, 0x9c, 0xd3, 0xd7, 0x9f, 0x7e, 0xec, 0xf0, 0x2f, 0x35, 0xcc, 0xae,
	0xff, 0x5f, 0x65, 0xa7, 0x89, 0xf2, 0x32, 0xe8, 0x2f, 0x06, 0x97, 0x93, 0xd5, 0xd3, 0x2e, 0xe3,
	0xed, 0xc2, 0x71, 0x7b, 0x19, 0x83, 0x94, 0x2f, 0xbe, 0xc1, 0x93, 0x8e, 0x6c, 0xd8, 0x29, 0xf4,
	0x3f, 0xfd, 0xf0, 0x1f, 0x31, 0x1f, 0xbc, 0x1b, 0x9e, 0xe3, 0x47, 0xd2, 0x37, 0x74, 0x28, 0x32,
	0xbf, 0xc7, 0xe6, 0x70, 0xd6, 0x0c, 0x6c, 0xb9, 0xc0, 0xf7, 0x5c, 0x70, 0xbd, 0xf9, 0x9d, 0x22,
	0x66, 0x7e, 0xdf, 0xa8, 0x1b, 0x6e, 0xa3, 0x14, 0x29, 0x7f, 0xb0, 0xfa, 0xdb, 0x83, 0x91, 0x5b,
	0x81, 0xbd, 0x81, 0x61, 0x53, 0xb3, 0x8b, 0xee, 0xeb, 0xd8, 0x9c, 0xe7, 0xe7, 0xdd, 0xa4, 0xc9,
	0xe7, 0x35, 0x9c, 0xd8, 0xd8, 0x83, 0x07, 0x92, 0xa3, 0x23, 0xcd, 0xcf, 0x3a, 0x18, 0x99, 0x57,
	0x6f, 0x5f, 0x7e, 0x8d, 0xee, 0xb9, 0xde, 0x1f, 0x6e, 0x0d, 0x1f, 0x7d, 0xa8, 0x14, 0x4f, 0xb2,
	0xab, 0x6b, 0x85, 0x89, 0xe0, 0x45, 0x94, 0x92, 0xc2, 0x2b, 0xfa, 0x1e, 0xb9, 0xd1, 0xe8, 0x5e,
	0xc9, 0xf4, 0xf6, 0xd4, 0xfe, 0xec, 0xaf, 0xfe, 0x05, 0x00, 0x00, 0xff, 0xff, 0x0b, 0xcc, 0x32,
	0xee, 0xfe, 0x02, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// CompilerClient is the client API for Compiler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type CompilerClient interface {
	Compile(ctx context.Context, in *CompileRequest, opts ...grpc.CallOption) (*CompileReply, error)
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoReply, error)
}

type compilerClient struct {
	cc *grpc.ClientConn
}

func NewCompilerClient(cc *grpc.ClientConn) CompilerClient {
	return &compilerClient{cc}
}

func (c *compilerClient) Compile(ctx context.Context, in *CompileRequest, opts ...grpc.CallOption) (*CompileReply, error) {
	out := new(CompileReply)
	err := c.cc.Invoke(ctx, "/compilerrpc.Compiler/Compile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *compilerClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoReply, error) {
	out := new(InfoReply)
	err := c.cc.Invoke(ctx, "/compilerrpc.Compiler/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CompilerServer is the server API for Compiler service.
type CompilerServer interface {
	Compile(context.Context, *CompileRequest) (*CompileReply, error)
	Info(context.Context, *InfoRequest) (*InfoReply, error)
}

// UnimplementedCompilerServer can be embedded to have forward compatible implementations.
type UnimplementedCompilerServer struct {
}

func (*UnimplementedCompilerServer) Compile(ctx context.Context, req *CompileRequest) (*CompileReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Compile not implemented")
}
func (*UnimplementedCompilerServer) Info(ctx context.Context, req *InfoRequest) (*InfoReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}

func RegisterCompilerServer(s *grpc.Server, srv CompilerServer) {
	s.RegisterService(&_Compiler_serviceDesc, srv)
}

func _Compiler_Compile_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CompileRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompilerServer).Compile(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/compilerrpc.Compiler/Compile",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompilerServer).Compile(ctx, req.(*CompileRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Compiler_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CompilerServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/compilerrpc.Compiler/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CompilerServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Compiler_serviceDesc = grpc.ServiceDesc{
	ServiceName: "compilerrpc.Compiler",
	HandlerType: (*CompilerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Compile",
			Handler:    _Compiler_Compile_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _Compiler_Info_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "compiler.proto",
}
