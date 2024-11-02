// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: protobuf/core/data.proto

package core

import (
	context "context"
	meta "github.com/AEnjoy/IoT-lubricant/protobuf/meta"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	CoreService_Ping_FullMethodName          = "/lubricant.coreService/ping"
	CoreService_GetTask_FullMethodName       = "/lubricant.coreService/getTask"
	CoreService_PushMessageId_FullMethodName = "/lubricant.coreService/pushMessageId"
	CoreService_PushData_FullMethodName      = "/lubricant.coreService/pushData"
)

// CoreServiceClient is the client API for CoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CoreServiceClient interface {
	Ping(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[meta.Ping, meta.Ping], error)
	GetTask(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[Task, Task], error)
	PushMessageId(ctx context.Context, in *MessageIdInfo, opts ...grpc.CallOption) (*MessageIdInfo, error)
	PushData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[Data, Data], error)
}

type coreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCoreServiceClient(cc grpc.ClientConnInterface) CoreServiceClient {
	return &coreServiceClient{cc}
}

func (c *coreServiceClient) Ping(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[meta.Ping, meta.Ping], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &CoreService_ServiceDesc.Streams[0], CoreService_Ping_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[meta.Ping, meta.Ping]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CoreService_PingClient = grpc.BidiStreamingClient[meta.Ping, meta.Ping]

func (c *coreServiceClient) GetTask(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[Task, Task], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &CoreService_ServiceDesc.Streams[1], CoreService_GetTask_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Task, Task]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CoreService_GetTaskClient = grpc.BidiStreamingClient[Task, Task]

func (c *coreServiceClient) PushMessageId(ctx context.Context, in *MessageIdInfo, opts ...grpc.CallOption) (*MessageIdInfo, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(MessageIdInfo)
	err := c.cc.Invoke(ctx, CoreService_PushMessageId_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *coreServiceClient) PushData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[Data, Data], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &CoreService_ServiceDesc.Streams[2], CoreService_PushData_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[Data, Data]{ClientStream: stream}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CoreService_PushDataClient = grpc.BidiStreamingClient[Data, Data]

// CoreServiceServer is the server API for CoreService service.
// All implementations must embed UnimplementedCoreServiceServer
// for forward compatibility.
type CoreServiceServer interface {
	Ping(grpc.BidiStreamingServer[meta.Ping, meta.Ping]) error
	GetTask(grpc.BidiStreamingServer[Task, Task]) error
	PushMessageId(context.Context, *MessageIdInfo) (*MessageIdInfo, error)
	PushData(grpc.BidiStreamingServer[Data, Data]) error
	mustEmbedUnimplementedCoreServiceServer()
}

// UnimplementedCoreServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCoreServiceServer struct{}

func (UnimplementedCoreServiceServer) Ping(grpc.BidiStreamingServer[meta.Ping, meta.Ping]) error {
	return status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedCoreServiceServer) GetTask(grpc.BidiStreamingServer[Task, Task]) error {
	return status.Errorf(codes.Unimplemented, "method GetTask not implemented")
}
func (UnimplementedCoreServiceServer) PushMessageId(context.Context, *MessageIdInfo) (*MessageIdInfo, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PushMessageId not implemented")
}
func (UnimplementedCoreServiceServer) PushData(grpc.BidiStreamingServer[Data, Data]) error {
	return status.Errorf(codes.Unimplemented, "method PushData not implemented")
}
func (UnimplementedCoreServiceServer) mustEmbedUnimplementedCoreServiceServer() {}
func (UnimplementedCoreServiceServer) testEmbeddedByValue()                     {}

// UnsafeCoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CoreServiceServer will
// result in compilation errors.
type UnsafeCoreServiceServer interface {
	mustEmbedUnimplementedCoreServiceServer()
}

func RegisterCoreServiceServer(s grpc.ServiceRegistrar, srv CoreServiceServer) {
	// If the following call pancis, it indicates UnimplementedCoreServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CoreService_ServiceDesc, srv)
}

func _CoreService_Ping_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CoreServiceServer).Ping(&grpc.GenericServerStream[meta.Ping, meta.Ping]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CoreService_PingServer = grpc.BidiStreamingServer[meta.Ping, meta.Ping]

func _CoreService_GetTask_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CoreServiceServer).GetTask(&grpc.GenericServerStream[Task, Task]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CoreService_GetTaskServer = grpc.BidiStreamingServer[Task, Task]

func _CoreService_PushMessageId_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MessageIdInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CoreServiceServer).PushMessageId(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CoreService_PushMessageId_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CoreServiceServer).PushMessageId(ctx, req.(*MessageIdInfo))
	}
	return interceptor(ctx, in, info, handler)
}

func _CoreService_PushData_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(CoreServiceServer).PushData(&grpc.GenericServerStream[Data, Data]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type CoreService_PushDataServer = grpc.BidiStreamingServer[Data, Data]

// CoreService_ServiceDesc is the grpc.ServiceDesc for CoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lubricant.coreService",
	HandlerType: (*CoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "pushMessageId",
			Handler:    _CoreService_PushMessageId_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ping",
			Handler:       _CoreService_Ping_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "getTask",
			Handler:       _CoreService_GetTask_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "pushData",
			Handler:       _CoreService_PushData_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "protobuf/core/data.proto",
}
