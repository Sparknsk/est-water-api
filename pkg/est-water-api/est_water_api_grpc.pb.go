// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package est_water_api

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

// EstWaterApiServiceClient is the client API for EstWaterApiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EstWaterApiServiceClient interface {
	// CreateWaterV1 - Create a water autotransport
	CreateWaterV1(ctx context.Context, in *CreateWaterV1Request, opts ...grpc.CallOption) (*CreateWaterV1Response, error)
	// DescribeWaterV1 - Describe a water autotransport
	DescribeWaterV1(ctx context.Context, in *DescribeWaterV1Request, opts ...grpc.CallOption) (*DescribeWaterV1Response, error)
	// ListWatersV1 - List of water autotransports
	ListWatersV1(ctx context.Context, in *ListWatersV1Request, opts ...grpc.CallOption) (*ListWatersV1Response, error)
	// RemoveWaterV1 - Remove a water autotransport
	RemoveWaterV1(ctx context.Context, in *RemoveWaterV1Request, opts ...grpc.CallOption) (*RemoveWaterV1Response, error)
}

type estWaterApiServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEstWaterApiServiceClient(cc grpc.ClientConnInterface) EstWaterApiServiceClient {
	return &estWaterApiServiceClient{cc}
}

func (c *estWaterApiServiceClient) CreateWaterV1(ctx context.Context, in *CreateWaterV1Request, opts ...grpc.CallOption) (*CreateWaterV1Response, error) {
	out := new(CreateWaterV1Response)
	err := c.cc.Invoke(ctx, "/ozonmp.est_water_api.v1.EstWaterApiService/CreateWaterV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *estWaterApiServiceClient) DescribeWaterV1(ctx context.Context, in *DescribeWaterV1Request, opts ...grpc.CallOption) (*DescribeWaterV1Response, error) {
	out := new(DescribeWaterV1Response)
	err := c.cc.Invoke(ctx, "/ozonmp.est_water_api.v1.EstWaterApiService/DescribeWaterV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *estWaterApiServiceClient) ListWatersV1(ctx context.Context, in *ListWatersV1Request, opts ...grpc.CallOption) (*ListWatersV1Response, error) {
	out := new(ListWatersV1Response)
	err := c.cc.Invoke(ctx, "/ozonmp.est_water_api.v1.EstWaterApiService/ListWatersV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *estWaterApiServiceClient) RemoveWaterV1(ctx context.Context, in *RemoveWaterV1Request, opts ...grpc.CallOption) (*RemoveWaterV1Response, error) {
	out := new(RemoveWaterV1Response)
	err := c.cc.Invoke(ctx, "/ozonmp.est_water_api.v1.EstWaterApiService/RemoveWaterV1", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EstWaterApiServiceServer is the server API for EstWaterApiService service.
// All implementations must embed UnimplementedEstWaterApiServiceServer
// for forward compatibility
type EstWaterApiServiceServer interface {
	// CreateWaterV1 - Create a water autotransport
	CreateWaterV1(context.Context, *CreateWaterV1Request) (*CreateWaterV1Response, error)
	// DescribeWaterV1 - Describe a water autotransport
	DescribeWaterV1(context.Context, *DescribeWaterV1Request) (*DescribeWaterV1Response, error)
	// ListWatersV1 - List of water autotransports
	ListWatersV1(context.Context, *ListWatersV1Request) (*ListWatersV1Response, error)
	// RemoveWaterV1 - Remove a water autotransport
	RemoveWaterV1(context.Context, *RemoveWaterV1Request) (*RemoveWaterV1Response, error)
	mustEmbedUnimplementedEstWaterApiServiceServer()
}

// UnimplementedEstWaterApiServiceServer must be embedded to have forward compatible implementations.
type UnimplementedEstWaterApiServiceServer struct {
}

func (UnimplementedEstWaterApiServiceServer) CreateWaterV1(context.Context, *CreateWaterV1Request) (*CreateWaterV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateWaterV1 not implemented")
}
func (UnimplementedEstWaterApiServiceServer) DescribeWaterV1(context.Context, *DescribeWaterV1Request) (*DescribeWaterV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DescribeWaterV1 not implemented")
}
func (UnimplementedEstWaterApiServiceServer) ListWatersV1(context.Context, *ListWatersV1Request) (*ListWatersV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListWatersV1 not implemented")
}
func (UnimplementedEstWaterApiServiceServer) RemoveWaterV1(context.Context, *RemoveWaterV1Request) (*RemoveWaterV1Response, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveWaterV1 not implemented")
}
func (UnimplementedEstWaterApiServiceServer) mustEmbedUnimplementedEstWaterApiServiceServer() {}

// UnsafeEstWaterApiServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EstWaterApiServiceServer will
// result in compilation errors.
type UnsafeEstWaterApiServiceServer interface {
	mustEmbedUnimplementedEstWaterApiServiceServer()
}

func RegisterEstWaterApiServiceServer(s grpc.ServiceRegistrar, srv EstWaterApiServiceServer) {
	s.RegisterService(&EstWaterApiService_ServiceDesc, srv)
}

func _EstWaterApiService_CreateWaterV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateWaterV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EstWaterApiServiceServer).CreateWaterV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozonmp.est_water_api.v1.EstWaterApiService/CreateWaterV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EstWaterApiServiceServer).CreateWaterV1(ctx, req.(*CreateWaterV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _EstWaterApiService_DescribeWaterV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DescribeWaterV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EstWaterApiServiceServer).DescribeWaterV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozonmp.est_water_api.v1.EstWaterApiService/DescribeWaterV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EstWaterApiServiceServer).DescribeWaterV1(ctx, req.(*DescribeWaterV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _EstWaterApiService_ListWatersV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListWatersV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EstWaterApiServiceServer).ListWatersV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozonmp.est_water_api.v1.EstWaterApiService/ListWatersV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EstWaterApiServiceServer).ListWatersV1(ctx, req.(*ListWatersV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

func _EstWaterApiService_RemoveWaterV1_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveWaterV1Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EstWaterApiServiceServer).RemoveWaterV1(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/ozonmp.est_water_api.v1.EstWaterApiService/RemoveWaterV1",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EstWaterApiServiceServer).RemoveWaterV1(ctx, req.(*RemoveWaterV1Request))
	}
	return interceptor(ctx, in, info, handler)
}

// EstWaterApiService_ServiceDesc is the grpc.ServiceDesc for EstWaterApiService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EstWaterApiService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "ozonmp.est_water_api.v1.EstWaterApiService",
	HandlerType: (*EstWaterApiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateWaterV1",
			Handler:    _EstWaterApiService_CreateWaterV1_Handler,
		},
		{
			MethodName: "DescribeWaterV1",
			Handler:    _EstWaterApiService_DescribeWaterV1_Handler,
		},
		{
			MethodName: "ListWatersV1",
			Handler:    _EstWaterApiService_ListWatersV1_Handler,
		},
		{
			MethodName: "RemoveWaterV1",
			Handler:    _EstWaterApiService_RemoveWaterV1_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ozonmp/est_water_api/v1/est_water_api.proto",
}
