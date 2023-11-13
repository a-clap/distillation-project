// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: distillationproto/process.proto

package distillationproto

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	empty "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// ProcessClient is the client API for Process service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ProcessClient interface {
	GetGlobalConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessGlobalConfig, error)
	GetPhaseCount(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessPhaseCount, error)
	GetPhaseConfig(ctx context.Context, in *PhaseNumber, opts ...grpc.CallOption) (*ProcessPhaseConfig, error)
	ConfigurePhaseCount(ctx context.Context, in *ProcessPhaseCount, opts ...grpc.CallOption) (*ProcessPhaseCount, error)
	ConfigurePhase(ctx context.Context, in *ProcessPhaseConfig, opts ...grpc.CallOption) (*ProcessPhaseConfig, error)
	ValidateConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessConfigValidation, error)
	ConfigureGlobalGPIO(ctx context.Context, in *GlobalGPIOConfig, opts ...grpc.CallOption) (*GlobalGPIOConfig, error)
	EnableProcess(ctx context.Context, in *ProcessConfig, opts ...grpc.CallOption) (*ProcessConfig, error)
	Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessStatus, error)
}

type processClient struct {
	cc grpc.ClientConnInterface
}

func NewProcessClient(cc grpc.ClientConnInterface) ProcessClient {
	return &processClient{cc}
}

func (c *processClient) GetGlobalConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessGlobalConfig, error) {
	out := new(ProcessGlobalConfig)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/GetGlobalConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) GetPhaseCount(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessPhaseCount, error) {
	out := new(ProcessPhaseCount)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/GetPhaseCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) GetPhaseConfig(ctx context.Context, in *PhaseNumber, opts ...grpc.CallOption) (*ProcessPhaseConfig, error) {
	out := new(ProcessPhaseConfig)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/GetPhaseConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) ConfigurePhaseCount(ctx context.Context, in *ProcessPhaseCount, opts ...grpc.CallOption) (*ProcessPhaseCount, error) {
	out := new(ProcessPhaseCount)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/ConfigurePhaseCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) ConfigurePhase(ctx context.Context, in *ProcessPhaseConfig, opts ...grpc.CallOption) (*ProcessPhaseConfig, error) {
	out := new(ProcessPhaseConfig)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/ConfigurePhase", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) ValidateConfig(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessConfigValidation, error) {
	out := new(ProcessConfigValidation)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/ValidateConfig", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) ConfigureGlobalGPIO(ctx context.Context, in *GlobalGPIOConfig, opts ...grpc.CallOption) (*GlobalGPIOConfig, error) {
	out := new(GlobalGPIOConfig)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/ConfigureGlobalGPIO", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) EnableProcess(ctx context.Context, in *ProcessConfig, opts ...grpc.CallOption) (*ProcessConfig, error) {
	out := new(ProcessConfig)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/EnableProcess", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *processClient) Status(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ProcessStatus, error) {
	out := new(ProcessStatus)
	err := c.cc.Invoke(ctx, "/distillationproto.Process/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ProcessServer is the server API for Process service.
// All implementations must embed UnimplementedProcessServer
// for forward compatibility
type ProcessServer interface {
	GetGlobalConfig(context.Context, *empty.Empty) (*ProcessGlobalConfig, error)
	GetPhaseCount(context.Context, *empty.Empty) (*ProcessPhaseCount, error)
	GetPhaseConfig(context.Context, *PhaseNumber) (*ProcessPhaseConfig, error)
	ConfigurePhaseCount(context.Context, *ProcessPhaseCount) (*ProcessPhaseCount, error)
	ConfigurePhase(context.Context, *ProcessPhaseConfig) (*ProcessPhaseConfig, error)
	ValidateConfig(context.Context, *empty.Empty) (*ProcessConfigValidation, error)
	ConfigureGlobalGPIO(context.Context, *GlobalGPIOConfig) (*GlobalGPIOConfig, error)
	EnableProcess(context.Context, *ProcessConfig) (*ProcessConfig, error)
	Status(context.Context, *empty.Empty) (*ProcessStatus, error)
	mustEmbedUnimplementedProcessServer()
}

// UnimplementedProcessServer must be embedded to have forward compatible implementations.
type UnimplementedProcessServer struct {
}

func (UnimplementedProcessServer) GetGlobalConfig(context.Context, *empty.Empty) (*ProcessGlobalConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGlobalConfig not implemented")
}
func (UnimplementedProcessServer) GetPhaseCount(context.Context, *empty.Empty) (*ProcessPhaseCount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhaseCount not implemented")
}
func (UnimplementedProcessServer) GetPhaseConfig(context.Context, *PhaseNumber) (*ProcessPhaseConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPhaseConfig not implemented")
}
func (UnimplementedProcessServer) ConfigurePhaseCount(context.Context, *ProcessPhaseCount) (*ProcessPhaseCount, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigurePhaseCount not implemented")
}
func (UnimplementedProcessServer) ConfigurePhase(context.Context, *ProcessPhaseConfig) (*ProcessPhaseConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigurePhase not implemented")
}
func (UnimplementedProcessServer) ValidateConfig(context.Context, *empty.Empty) (*ProcessConfigValidation, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ValidateConfig not implemented")
}
func (UnimplementedProcessServer) ConfigureGlobalGPIO(context.Context, *GlobalGPIOConfig) (*GlobalGPIOConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ConfigureGlobalGPIO not implemented")
}
func (UnimplementedProcessServer) EnableProcess(context.Context, *ProcessConfig) (*ProcessConfig, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EnableProcess not implemented")
}
func (UnimplementedProcessServer) Status(context.Context, *empty.Empty) (*ProcessStatus, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedProcessServer) mustEmbedUnimplementedProcessServer() {}

// UnsafeProcessServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ProcessServer will
// result in compilation errors.
type UnsafeProcessServer interface {
	mustEmbedUnimplementedProcessServer()
}

func RegisterProcessServer(s grpc.ServiceRegistrar, srv ProcessServer) {
	s.RegisterService(&Process_ServiceDesc, srv)
}

func _Process_GetGlobalConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).GetGlobalConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/GetGlobalConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).GetGlobalConfig(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_GetPhaseCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).GetPhaseCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/GetPhaseCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).GetPhaseCount(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_GetPhaseConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PhaseNumber)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).GetPhaseConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/GetPhaseConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).GetPhaseConfig(ctx, req.(*PhaseNumber))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_ConfigurePhaseCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessPhaseCount)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).ConfigurePhaseCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/ConfigurePhaseCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).ConfigurePhaseCount(ctx, req.(*ProcessPhaseCount))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_ConfigurePhase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessPhaseConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).ConfigurePhase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/ConfigurePhase",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).ConfigurePhase(ctx, req.(*ProcessPhaseConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_ValidateConfig_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).ValidateConfig(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/ValidateConfig",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).ValidateConfig(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_ConfigureGlobalGPIO_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GlobalGPIOConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).ConfigureGlobalGPIO(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/ConfigureGlobalGPIO",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).ConfigureGlobalGPIO(ctx, req.(*GlobalGPIOConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_EnableProcess_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ProcessConfig)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).EnableProcess(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/EnableProcess",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).EnableProcess(ctx, req.(*ProcessConfig))
	}
	return interceptor(ctx, in, info, handler)
}

func _Process_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ProcessServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/distillationproto.Process/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ProcessServer).Status(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Process_ServiceDesc is the grpc.ServiceDesc for Process service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Process_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "distillationproto.Process",
	HandlerType: (*ProcessServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetGlobalConfig",
			Handler:    _Process_GetGlobalConfig_Handler,
		},
		{
			MethodName: "GetPhaseCount",
			Handler:    _Process_GetPhaseCount_Handler,
		},
		{
			MethodName: "GetPhaseConfig",
			Handler:    _Process_GetPhaseConfig_Handler,
		},
		{
			MethodName: "ConfigurePhaseCount",
			Handler:    _Process_ConfigurePhaseCount_Handler,
		},
		{
			MethodName: "ConfigurePhase",
			Handler:    _Process_ConfigurePhase_Handler,
		},
		{
			MethodName: "ValidateConfig",
			Handler:    _Process_ValidateConfig_Handler,
		},
		{
			MethodName: "ConfigureGlobalGPIO",
			Handler:    _Process_ConfigureGlobalGPIO_Handler,
		},
		{
			MethodName: "EnableProcess",
			Handler:    _Process_EnableProcess_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Process_Status_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "distillationproto/process.proto",
}
