// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: distillationproto/pt100.proto

package distillationproto

import (
	context "context"

	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// PTClient is the client API for PT service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PTClient interface {
	PTGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PTConfigs, error)
	PTConfigure(ctx context.Context, in *PTConfig, opts ...grpc.CallOption) (*PTConfig, error)
	PTGetTemperatures(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PTTemperatures, error)
}

type pTClient struct {
	cc grpc.ClientConnInterface
}

func NewPTClient(cc grpc.ClientConnInterface) PTClient {
	return &pTClient{cc}
}

func (c *pTClient) PTGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PTConfigs, error) {
	out := new(PTConfigs)
	err := c.cc.Invoke(ctx, "/distillationproto.PT/PTGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pTClient) PTConfigure(ctx context.Context, in *PTConfig, opts ...grpc.CallOption) (*PTConfig, error) {
	out := new(PTConfig)
	err := c.cc.Invoke(ctx, "/distillationproto.PT/PTConfigure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *pTClient) PTGetTemperatures(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*PTTemperatures, error) {
	out := new(PTTemperatures)
	err := c.cc.Invoke(ctx, "/distillationproto.PT/PTGetTemperatures", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PTServer is the server API for PT service.
// All implementations must embed UnimplementedPTServer
// for forward compatibility
type PTServer interface {
	PTGet(context.Context, *empty.Empty) (*PTConfigs, error)
	PTConfigure(context.Context, *PTConfig) (*PTConfig, error)
	PTGetTemperatures(context.Context, *empty.Empty) (*PTTemperatures, error)
	mustEmbedUnimplementedPTServer()
}

// UnimplementedPTServer must be embedded to have forward compatible implementations.
type UnimplementedPTServer struct {
}

func (UnimplementedPTServer) PTGet(context.Context, *empty.Empty) (*PTConfigs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PTGet not implemented")
}
func (UnimplementedPTServer) PTConfigure(context.Context, *PTConfig) (*PTConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PTConfigure not implemented")
}
func (UnimplementedPTServer) PTGetTemperatures(context.Context, *empty.Empty) (*PTTemperatures, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PTGetTemperatures not implemented")
}
func (UnimplementedPTServer) mustEmbedUnimplementedPTServer() {}

// UnsafePTServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PTServer will
// result in compilation errors.
type UnsafePTServer interface {
	mustEmbedUnimplementedPTServer()
}

func RegisterPTServer(s grpc.ServiceRegistrar, srv PTServer) {
	s.RegisterService(&PT_ServiceDesc, srv)
}

func _PT_PTGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PTServer).PTGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.PT/PTGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PTServer).PTGet(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PT_PTConfigure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PTConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PTServer).PTConfigure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.PT/PTConfigure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PTServer).PTConfigure(ctx, req.(*PTConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _PT_PTGetTemperatures_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PTServer).PTGetTemperatures(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.PT/PTGetTemperatures",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PTServer).PTGetTemperatures(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// PT_ServiceDesc is the grpc.ServiceDesc for PT service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PT_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "distillationproto.PT",
	HandlerType: (*PTServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "PTGet",
			Handler:    _PT_PTGet_Handler,
		},
		{
			MethodName: "PTConfigure",
			Handler:    _PT_PTConfigure_Handler,
		},
		{
			MethodName: "PTGetTemperatures",
			Handler:    _PT_PTGetTemperatures_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "distillationproto/pt100.proto",
}
