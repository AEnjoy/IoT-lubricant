// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v6.31.0--rc2
// source: protobuf/svc/datastore.proto

package svc

import (
	context "context"
	meta "github.com/aenjoy/iot-lubricant/protobuf/meta"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DataStoreService_CheckLinker_FullMethodName = "/lubricant.svc.DataStoreService/CheckLinker"
	DataStoreService_StoreData_FullMethodName   = "/lubricant.svc.DataStoreService/StoreData"
	DataStoreService_Ping_FullMethodName        = "/lubricant.svc.DataStoreService/ping"
)

// DataStoreServiceClient is the client API for DataStoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataStoreServiceClient interface {
	CheckLinker(ctx context.Context, in *CheckLinkerRequest, opts ...grpc.CallOption) (*CheckLinkerResponse, error)
	StoreData(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[StoreDataRequest, StoreDataResponse], error)
	Ping(ctx context.Context, in *meta.Ping, opts ...grpc.CallOption) (*meta.Ping, error)
}

type dataStoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataStoreServiceClient(cc grpc.ClientConnInterface) DataStoreServiceClient {
	return &dataStoreServiceClient{cc}
}

func (c *dataStoreServiceClient) CheckLinker(ctx context.Context, in *CheckLinkerRequest, opts ...grpc.CallOption) (*CheckLinkerResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CheckLinkerResponse)
	err := c.cc.Invoke(ctx, DataStoreService_CheckLinker_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *dataStoreServiceClient) StoreData(ctx context.Context, opts ...grpc.CallOption) (grpc.ClientStreamingClient[StoreDataRequest, StoreDataResponse], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &DataStoreService_ServiceDesc.Streams[0], DataStoreService_StoreData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[StoreDataRequest, StoreDataResponse]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataStoreService_StoreDataClient = grpc.ClientStreamingClient[StoreDataRequest, StoreDataResponse]

func (c *dataStoreServiceClient) Ping(ctx context.Context, in *meta.Ping, opts ...grpc.CallOption) (*meta.Ping, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(meta.Ping)
	err := c.cc.Invoke(ctx, DataStoreService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataStoreServiceServer is the server API for DataStoreService service.
// All implementations must embed UnimplementedDataStoreServiceServer
// for forward compatibility.
type DataStoreServiceServer interface {
	CheckLinker(context.Context, *CheckLinkerRequest) (*CheckLinkerResponse, error)
	StoreData(grpc.ClientStreamingServer[StoreDataRequest, StoreDataResponse]) error
	Ping(context.Context, *meta.Ping) (*meta.Ping, error)
	mustEmbedUnimplementedDataStoreServiceServer()
}

// UnimplementedDataStoreServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDataStoreServiceServer struct{}

func (UnimplementedDataStoreServiceServer) CheckLinker(context.Context, *CheckLinkerRequest) (*CheckLinkerResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CheckLinker not implemented")
}
func (UnimplementedDataStoreServiceServer) StoreData(grpc.ClientStreamingServer[StoreDataRequest, StoreDataResponse]) error {
	return status.Errorf(codes.Unimplemented, "method StoreData not implemented")
}
func (UnimplementedDataStoreServiceServer) Ping(context.Context, *meta.Ping) (*meta.Ping, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDataStoreServiceServer) mustEmbedUnimplementedDataStoreServiceServer() {}
func (UnimplementedDataStoreServiceServer) testEmbeddedByValue()                          {}

// UnsafeDataStoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataStoreServiceServer will
// result in compilation errors.
type UnsafeDataStoreServiceServer interface {
	mustEmbedUnimplementedDataStoreServiceServer()
}

func RegisterDataStoreServiceServer(s grpc.ServiceRegistrar, srv DataStoreServiceServer) {
	// If the following call pancis, it indicates UnimplementedDataStoreServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DataStoreService_ServiceDesc, srv)
}

func _DataStoreService_CheckLinker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CheckLinkerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataStoreServiceServer).CheckLinker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataStoreService_CheckLinker_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataStoreServiceServer).CheckLinker(ctx, req.(*CheckLinkerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DataStoreService_StoreData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(DataStoreServiceServer).StoreData(&grpc.GenericServerStream[StoreDataRequest, StoreDataResponse]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type DataStoreService_StoreDataServer = grpc.ClientStreamingServer[StoreDataRequest, StoreDataResponse]

func _DataStoreService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(meta.Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataStoreServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataStoreService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataStoreServiceServer).Ping(ctx, req.(*meta.Ping))
	}
	return interceptor(ctx, in, info, handler)
}

// DataStoreService_ServiceDesc is the grpc.ServiceDesc for DataStoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataStoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lubricant.svc.DataStoreService",
	HandlerType: (*DataStoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CheckLinker",
			Handler:    _DataStoreService_CheckLinker_Handler,
		},
		{
			MethodName: "ping",
			Handler:    _DataStoreService_Ping_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StoreData",
			Handler:       _DataStoreService_StoreData_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "protobuf/svc/datastore.proto",
}

const (
	DataStoreDebugService_Ping_FullMethodName = "/lubricant.svc.DataStoreDebugService/ping"
)

// DataStoreDebugServiceClient is the client API for DataStoreDebugService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DataStoreDebugServiceClient interface {
	Ping(ctx context.Context, in *meta.Ping, opts ...grpc.CallOption) (*meta.Ping, error)
}

type dataStoreDebugServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDataStoreDebugServiceClient(cc grpc.ClientConnInterface) DataStoreDebugServiceClient {
	return &dataStoreDebugServiceClient{cc}
}

func (c *dataStoreDebugServiceClient) Ping(ctx context.Context, in *meta.Ping, opts ...grpc.CallOption) (*meta.Ping, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(meta.Ping)
	err := c.cc.Invoke(ctx, DataStoreDebugService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DataStoreDebugServiceServer is the server API for DataStoreDebugService service.
// All implementations must embed UnimplementedDataStoreDebugServiceServer
// for forward compatibility.
type DataStoreDebugServiceServer interface {
	Ping(context.Context, *meta.Ping) (*meta.Ping, error)
	mustEmbedUnimplementedDataStoreDebugServiceServer()
}

// UnimplementedDataStoreDebugServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDataStoreDebugServiceServer struct{}

func (UnimplementedDataStoreDebugServiceServer) Ping(context.Context, *meta.Ping) (*meta.Ping, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedDataStoreDebugServiceServer) mustEmbedUnimplementedDataStoreDebugServiceServer() {}
func (UnimplementedDataStoreDebugServiceServer) testEmbeddedByValue()                               {}

// UnsafeDataStoreDebugServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DataStoreDebugServiceServer will
// result in compilation errors.
type UnsafeDataStoreDebugServiceServer interface {
	mustEmbedUnimplementedDataStoreDebugServiceServer()
}

func RegisterDataStoreDebugServiceServer(s grpc.ServiceRegistrar, srv DataStoreDebugServiceServer) {
	// If the following call pancis, it indicates UnimplementedDataStoreDebugServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DataStoreDebugService_ServiceDesc, srv)
}

func _DataStoreDebugService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(meta.Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DataStoreDebugServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DataStoreDebugService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DataStoreDebugServiceServer).Ping(ctx, req.(*meta.Ping))
	}
	return interceptor(ctx, in, info, handler)
}

// DataStoreDebugService_ServiceDesc is the grpc.ServiceDesc for DataStoreDebugService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DataStoreDebugService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lubricant.svc.DataStoreDebugService",
	HandlerType: (*DataStoreDebugServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ping",
			Handler:    _DataStoreDebugService_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protobuf/svc/datastore.proto",
}
