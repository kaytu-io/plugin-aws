// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.12.4
// source: plugin/proto/aws_server.proto

package aws

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	Optimization_EC2InstanceOptimization_FullMethodName = "/pluginaws.optimization.v1.Optimization/EC2InstanceOptimization"
	Optimization_RDSInstanceOptimization_FullMethodName = "/pluginaws.optimization.v1.Optimization/RDSInstanceOptimization"
	Optimization_RDSClusterOptimization_FullMethodName  = "/pluginaws.optimization.v1.Optimization/RDSClusterOptimization"
)

// OptimizationClient is the client API for Optimization service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OptimizationClient interface {
	EC2InstanceOptimization(ctx context.Context, in *EC2InstanceOptimizationRequest, opts ...grpc.CallOption) (*EC2InstanceOptimizationResponse, error)
	RDSInstanceOptimization(ctx context.Context, in *RDSInstanceOptimizationRequest, opts ...grpc.CallOption) (*RDSInstanceOptimizationResponse, error)
	RDSClusterOptimization(ctx context.Context, in *RDSClusterOptimizationRequest, opts ...grpc.CallOption) (*RDSClusterOptimizationResponse, error)
}

type optimizationClient struct {
	cc grpc.ClientConnInterface
}

func NewOptimizationClient(cc grpc.ClientConnInterface) OptimizationClient {
	return &optimizationClient{cc}
}

func (c *optimizationClient) EC2InstanceOptimization(ctx context.Context, in *EC2InstanceOptimizationRequest, opts ...grpc.CallOption) (*EC2InstanceOptimizationResponse, error) {
	out := new(EC2InstanceOptimizationResponse)
	err := c.cc.Invoke(ctx, Optimization_EC2InstanceOptimization_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *optimizationClient) RDSInstanceOptimization(ctx context.Context, in *RDSInstanceOptimizationRequest, opts ...grpc.CallOption) (*RDSInstanceOptimizationResponse, error) {
	out := new(RDSInstanceOptimizationResponse)
	err := c.cc.Invoke(ctx, Optimization_RDSInstanceOptimization_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *optimizationClient) RDSClusterOptimization(ctx context.Context, in *RDSClusterOptimizationRequest, opts ...grpc.CallOption) (*RDSClusterOptimizationResponse, error) {
	out := new(RDSClusterOptimizationResponse)
	err := c.cc.Invoke(ctx, Optimization_RDSClusterOptimization_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OptimizationServer is the server API for Optimization service.
// All implementations must embed UnimplementedOptimizationServer
// for forward compatibility
type OptimizationServer interface {
	EC2InstanceOptimization(context.Context, *EC2InstanceOptimizationRequest) (*EC2InstanceOptimizationResponse, error)
	RDSInstanceOptimization(context.Context, *RDSInstanceOptimizationRequest) (*RDSInstanceOptimizationResponse, error)
	RDSClusterOptimization(context.Context, *RDSClusterOptimizationRequest) (*RDSClusterOptimizationResponse, error)
	mustEmbedUnimplementedOptimizationServer()
}

// UnimplementedOptimizationServer must be embedded to have forward compatible implementations.
type UnimplementedOptimizationServer struct {
}

func (UnimplementedOptimizationServer) EC2InstanceOptimization(context.Context, *EC2InstanceOptimizationRequest) (*EC2InstanceOptimizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method EC2InstanceOptimization not implemented")
}
func (UnimplementedOptimizationServer) RDSInstanceOptimization(context.Context, *RDSInstanceOptimizationRequest) (*RDSInstanceOptimizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RDSInstanceOptimization not implemented")
}
func (UnimplementedOptimizationServer) RDSClusterOptimization(context.Context, *RDSClusterOptimizationRequest) (*RDSClusterOptimizationResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RDSClusterOptimization not implemented")
}
func (UnimplementedOptimizationServer) mustEmbedUnimplementedOptimizationServer() {}

// UnsafeOptimizationServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OptimizationServer will
// result in compilation errors.
type UnsafeOptimizationServer interface {
	mustEmbedUnimplementedOptimizationServer()
}

func RegisterOptimizationServer(s grpc.ServiceRegistrar, srv OptimizationServer) {
	s.RegisterService(&Optimization_ServiceDesc, srv)
}

func _Optimization_EC2InstanceOptimization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EC2InstanceOptimizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OptimizationServer).EC2InstanceOptimization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Optimization_EC2InstanceOptimization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OptimizationServer).EC2InstanceOptimization(ctx, req.(*EC2InstanceOptimizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Optimization_RDSInstanceOptimization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSInstanceOptimizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OptimizationServer).RDSInstanceOptimization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Optimization_RDSInstanceOptimization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OptimizationServer).RDSInstanceOptimization(ctx, req.(*RDSInstanceOptimizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Optimization_RDSClusterOptimization_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RDSClusterOptimizationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OptimizationServer).RDSClusterOptimization(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Optimization_RDSClusterOptimization_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OptimizationServer).RDSClusterOptimization(ctx, req.(*RDSClusterOptimizationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Optimization_ServiceDesc is the grpc.ServiceDesc for Optimization service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Optimization_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pluginaws.optimization.v1.Optimization",
	HandlerType: (*OptimizationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "EC2InstanceOptimization",
			Handler:    _Optimization_EC2InstanceOptimization_Handler,
		},
		{
			MethodName: "RDSInstanceOptimization",
			Handler:    _Optimization_RDSInstanceOptimization_Handler,
		},
		{
			MethodName: "RDSClusterOptimization",
			Handler:    _Optimization_RDSClusterOptimization_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "plugin/proto/aws_server.proto",
}