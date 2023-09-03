// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: embeddedproto/heaters.proto

package embeddedproto

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

// HeaterClient is the client API for Heater service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HeaterClient interface {
	HeaterGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HeaterConfigs, error)
	HeaterConfigure(ctx context.Context, in *HeaterConfig, opts ...grpc.CallOption) (*HeaterConfig, error)
}

type heaterClient struct {
	cc grpc.ClientConnInterface
}

func NewHeaterClient(cc grpc.ClientConnInterface) HeaterClient {
	return &heaterClient{cc}
}

func (c *heaterClient) HeaterGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*HeaterConfigs, error) {
	out := new(HeaterConfigs)
	err := c.cc.Invoke(ctx, "/embeddedproto.Heater/HeaterGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *heaterClient) HeaterConfigure(ctx context.Context, in *HeaterConfig, opts ...grpc.CallOption) (*HeaterConfig, error) {
	out := new(HeaterConfig)
	err := c.cc.Invoke(ctx, "/embeddedproto.Heater/HeaterConfigure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HeaterServer is the server API for Heater service.
// All implementations must embed UnimplementedHeaterServer
// for forward compatibility
type HeaterServer interface {
	HeaterGet(context.Context, *empty.Empty) (*HeaterConfigs, error)
	HeaterConfigure(context.Context, *HeaterConfig) (*HeaterConfig, error)
	mustEmbedUnimplementedHeaterServer()
}

// UnimplementedHeaterServer must be embedded to have forward compatible implementations.
type UnimplementedHeaterServer struct {
}

func (UnimplementedHeaterServer) HeaterGet(context.Context, *empty.Empty) (*HeaterConfigs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeaterGet not implemented")
}
func (UnimplementedHeaterServer) HeaterConfigure(context.Context, *HeaterConfig) (*HeaterConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HeaterConfigure not implemented")
}
func (UnimplementedHeaterServer) mustEmbedUnimplementedHeaterServer() {}

// UnsafeHeaterServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HeaterServer will
// result in compilation errors.
type UnsafeHeaterServer interface {
	mustEmbedUnimplementedHeaterServer()
}

func RegisterHeaterServer(s grpc.ServiceRegistrar, srv HeaterServer) {
	s.RegisterService(&Heater_ServiceDesc, srv)
}

func _Heater_HeaterGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HeaterServer).HeaterGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/embeddedproto.Heater/HeaterGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HeaterServer).HeaterGet(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Heater_HeaterConfigure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HeaterConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HeaterServer).HeaterConfigure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/embeddedproto.Heater/HeaterConfigure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HeaterServer).HeaterConfigure(ctx, req.(*HeaterConfig))
	}
	return interceptor(ctx, in, info, handler)
}

// Heater_ServiceDesc is the grpc.ServiceDesc for Heater service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Heater_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "embeddedproto.Heater",
	HandlerType: (*HeaterServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "HeaterGet",
			Handler:    _Heater_HeaterGet_Handler,
		},
		{
			MethodName: "HeaterConfigure",
			Handler:    _Heater_HeaterConfigure_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "embeddedproto/heaters.proto",
}
