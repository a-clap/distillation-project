// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: distillationproto/gpio.proto

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

// GPIOClient is the client API for GPIO service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type GPIOClient interface {
	GPIOGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GPIOConfigs, error)
	GPIOConfigure(ctx context.Context, in *GPIOConfig, opts ...grpc.CallOption) (*GPIOConfig, error)
}

type gPIOClient struct {
	cc grpc.ClientConnInterface
}

func NewGPIOClient(cc grpc.ClientConnInterface) GPIOClient {
	return &gPIOClient{cc}
}

func (c *gPIOClient) GPIOGet(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*GPIOConfigs, error) {
	out := new(GPIOConfigs)
	err := c.cc.Invoke(ctx, "/distillationproto.GPIO/GPIOGet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *gPIOClient) GPIOConfigure(ctx context.Context, in *GPIOConfig, opts ...grpc.CallOption) (*GPIOConfig, error) {
	out := new(GPIOConfig)
	err := c.cc.Invoke(ctx, "/distillationproto.GPIO/GPIOConfigure", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GPIOServer is the server API for GPIO service.
// All implementations must embed UnimplementedGPIOServer
// for forward compatibility
type GPIOServer interface {
	GPIOGet(context.Context, *empty.Empty) (*GPIOConfigs, error)
	GPIOConfigure(context.Context, *GPIOConfig) (*GPIOConfig, error)
	mustEmbedUnimplementedGPIOServer()
}

// UnimplementedGPIOServer must be embedded to have forward compatible implementations.
type UnimplementedGPIOServer struct {
}

func (UnimplementedGPIOServer) GPIOGet(context.Context, *empty.Empty) (*GPIOConfigs, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GPIOGet not implemented")
}
func (UnimplementedGPIOServer) GPIOConfigure(context.Context, *GPIOConfig) (*GPIOConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GPIOConfigure not implemented")
}
func (UnimplementedGPIOServer) mustEmbedUnimplementedGPIOServer() {}

// UnsafeGPIOServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GPIOServer will
// result in compilation errors.
type UnsafeGPIOServer interface {
	mustEmbedUnimplementedGPIOServer()
}

func RegisterGPIOServer(s grpc.ServiceRegistrar, srv GPIOServer) {
	s.RegisterService(&GPIO_ServiceDesc, srv)
}

func _GPIO_GPIOGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPIOServer).GPIOGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.GPIO/GPIOGet",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPIOServer).GPIOGet(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _GPIO_GPIOConfigure_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GPIOConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GPIOServer).GPIOConfigure(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.GPIO/GPIOConfigure",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GPIOServer).GPIOConfigure(ctx, req.(*GPIOConfig))
	}
	return interceptor(ctx, in, info, handler)
}

// GPIO_ServiceDesc is the grpc.ServiceDesc for GPIO service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GPIO_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "distillationproto.GPIO",
	HandlerType: (*GPIOServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GPIOGet",
			Handler:    _GPIO_GPIOGet_Handler,
		},
		{
			MethodName: "GPIOConfigure",
			Handler:    _GPIO_GPIOConfigure_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "distillationproto/gpio.proto",
}