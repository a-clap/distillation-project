// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: time.proto

package osproto

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	timestamp "github.com/golang/protobuf/ptypes/timestamp"
	wrappers "github.com/golang/protobuf/ptypes/wrappers"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TimeClient is the client API for Time service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimeClient interface {
	Now(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*timestamp.Timestamp, error)
	SetNow(ctx context.Context, in *timestamp.Timestamp, opts ...grpc.CallOption) (*empty.Empty, error)
	NTP(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*wrappers.BoolValue, error)
	SetNTP(ctx context.Context, in *wrappers.BoolValue, opts ...grpc.CallOption) (*empty.Empty, error)
}

type timeClient struct {
	cc grpc.ClientConnInterface
}

func NewTimeClient(cc grpc.ClientConnInterface) TimeClient {
	return &timeClient{cc}
}

func (c *timeClient) Now(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*timestamp.Timestamp, error) {
	out := new(timestamp.Timestamp)
	err := c.cc.Invoke(ctx, "/Time/Now", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *timeClient) SetNow(ctx context.Context, in *timestamp.Timestamp, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/Time/SetNow", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *timeClient) NTP(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*wrappers.BoolValue, error) {
	out := new(wrappers.BoolValue)
	err := c.cc.Invoke(ctx, "/Time/NTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *timeClient) SetNTP(ctx context.Context, in *wrappers.BoolValue, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/Time/SetNTP", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TimeServer is the server API for Time service.
// All implementations must embed UnimplementedTimeServer
// for forward compatibility
type TimeServer interface {
	Now(context.Context, *empty.Empty) (*timestamp.Timestamp, error)
	SetNow(context.Context, *timestamp.Timestamp) (*empty.Empty, error)
	NTP(context.Context, *empty.Empty) (*wrappers.BoolValue, error)
	SetNTP(context.Context, *wrappers.BoolValue) (*empty.Empty, error)
	mustEmbedUnimplementedTimeServer()
}

// UnimplementedTimeServer must be embedded to have forward compatible implementations.
type UnimplementedTimeServer struct {
}

func (UnimplementedTimeServer) Now(context.Context, *empty.Empty) (*timestamp.Timestamp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Now not implemented")
}
func (UnimplementedTimeServer) SetNow(context.Context, *timestamp.Timestamp) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNow not implemented")
}
func (UnimplementedTimeServer) NTP(context.Context, *empty.Empty) (*wrappers.BoolValue, error) {
	return nil, status.Errorf(codes.Unimplemented, "method NTP not implemented")
}
func (UnimplementedTimeServer) SetNTP(context.Context, *wrappers.BoolValue) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetNTP not implemented")
}
func (UnimplementedTimeServer) mustEmbedUnimplementedTimeServer() {}

// UnsafeTimeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimeServer will
// result in compilation errors.
type UnsafeTimeServer interface {
	mustEmbedUnimplementedTimeServer()
}

func RegisterTimeServer(s grpc.ServiceRegistrar, srv TimeServer) {
	s.RegisterService(&Time_ServiceDesc, srv)
}

func _Time_Now_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeServer).Now(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Time/Now",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeServer).Now(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Time_SetNow_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(timestamp.Timestamp)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeServer).SetNow(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Time/SetNow",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeServer).SetNow(ctx, req.(*timestamp.Timestamp))
	}
	return interceptor(ctx, in, info, handler)
}

func _Time_NTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeServer).NTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Time/NTP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeServer).NTP(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Time_SetNTP_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(wrappers.BoolValue)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TimeServer).SetNTP(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Time/SetNTP",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TimeServer).SetNTP(ctx, req.(*wrappers.BoolValue))
	}
	return interceptor(ctx, in, info, handler)
}

// Time_ServiceDesc is the grpc.ServiceDesc for Time service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Time_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Time",
	HandlerType: (*TimeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Now",
			Handler:    _Time_Now_Handler,
		},
		{
			MethodName: "SetNow",
			Handler:    _Time_SetNow_Handler,
		},
		{
			MethodName: "NTP",
			Handler:    _Time_NTP_Handler,
		},
		{
			MethodName: "SetNTP",
			Handler:    _Time_SetNTP_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "time.proto",
}