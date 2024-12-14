// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: protobuf/agent/agent.proto

package agent

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
	EdgeService_Ping_FullMethodName            = "/lubricant.agent.edgeService/ping"
	EdgeService_RegisterGateway_FullMethodName = "/lubricant.agent.edgeService/registerGateway"
	EdgeService_SetAgent_FullMethodName        = "/lubricant.agent.edgeService/setAgent"
	EdgeService_GetOpenapiDoc_FullMethodName   = "/lubricant.agent.edgeService/getOpenapiDoc"
	EdgeService_GetAgentInfo_FullMethodName    = "/lubricant.agent.edgeService/getAgentInfo"
	EdgeService_GetGatherData_FullMethodName   = "/lubricant.agent.edgeService/getGatherData"
	EdgeService_GetDataStream_FullMethodName   = "/lubricant.agent.edgeService/GetDataStream"
	EdgeService_SendHttpMethod_FullMethodName  = "/lubricant.agent.edgeService/sendHttpMethod"
	EdgeService_StartGather_FullMethodName     = "/lubricant.agent.edgeService/startGather"
	EdgeService_StopGather_FullMethodName      = "/lubricant.agent.edgeService/stopGather"
)

// EdgeServiceClient is the client API for EdgeService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type EdgeServiceClient interface {
	Ping(ctx context.Context, in *meta.Ping, opts ...grpc.CallOption) (*meta.Ping, error)
	RegisterGateway(ctx context.Context, in *RegisterGatewayRequest, opts ...grpc.CallOption) (*RegisterGatewayResponse, error)
	SetAgent(ctx context.Context, in *SetAgentRequest, opts ...grpc.CallOption) (*SetAgentResponse, error)
	GetOpenapiDoc(ctx context.Context, in *GetOpenapiDocRequest, opts ...grpc.CallOption) (*OpenapiDoc, error)
	GetAgentInfo(ctx context.Context, in *GetAgentInfoRequest, opts ...grpc.CallOption) (*GetAgentInfoResponse, error)
	GetGatherData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*DataMessage, error)
	GetDataStream(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DataChunk], error)
	SendHttpMethod(ctx context.Context, in *SendHttpMethodRequest, opts ...grpc.CallOption) (*SendHttpMethodResponse, error)
	StartGather(ctx context.Context, in *StartGatherRequest, opts ...grpc.CallOption) (*meta.CommonResponse, error)
	StopGather(ctx context.Context, in *StopGatherRequest, opts ...grpc.CallOption) (*meta.CommonResponse, error)
}

type edgeServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewEdgeServiceClient(cc grpc.ClientConnInterface) EdgeServiceClient {
	return &edgeServiceClient{cc}
}

func (c *edgeServiceClient) Ping(ctx context.Context, in *meta.Ping, opts ...grpc.CallOption) (*meta.Ping, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(meta.Ping)
	err := c.cc.Invoke(ctx, EdgeService_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) RegisterGateway(ctx context.Context, in *RegisterGatewayRequest, opts ...grpc.CallOption) (*RegisterGatewayResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterGatewayResponse)
	err := c.cc.Invoke(ctx, EdgeService_RegisterGateway_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) SetAgent(ctx context.Context, in *SetAgentRequest, opts ...grpc.CallOption) (*SetAgentResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SetAgentResponse)
	err := c.cc.Invoke(ctx, EdgeService_SetAgent_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) GetOpenapiDoc(ctx context.Context, in *GetOpenapiDocRequest, opts ...grpc.CallOption) (*OpenapiDoc, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(OpenapiDoc)
	err := c.cc.Invoke(ctx, EdgeService_GetOpenapiDoc_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) GetAgentInfo(ctx context.Context, in *GetAgentInfoRequest, opts ...grpc.CallOption) (*GetAgentInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAgentInfoResponse)
	err := c.cc.Invoke(ctx, EdgeService_GetAgentInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) GetGatherData(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (*DataMessage, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DataMessage)
	err := c.cc.Invoke(ctx, EdgeService_GetGatherData_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) GetDataStream(ctx context.Context, in *GetDataRequest, opts ...grpc.CallOption) (grpc.ServerStreamingClient[DataChunk], error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	stream, err := c.cc.NewStream(ctx, &EdgeService_ServiceDesc.Streams[0], EdgeService_GetDataStream_FullMethodName, cOpts...)
	if err != nil {
		return nil, err
	}
	x := &grpc.GenericClientStream[GetDataRequest, DataChunk]{ClientStream: stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EdgeService_GetDataStreamClient = grpc.ServerStreamingClient[DataChunk]

func (c *edgeServiceClient) SendHttpMethod(ctx context.Context, in *SendHttpMethodRequest, opts ...grpc.CallOption) (*SendHttpMethodResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(SendHttpMethodResponse)
	err := c.cc.Invoke(ctx, EdgeService_SendHttpMethod_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) StartGather(ctx context.Context, in *StartGatherRequest, opts ...grpc.CallOption) (*meta.CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(meta.CommonResponse)
	err := c.cc.Invoke(ctx, EdgeService_StartGather_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *edgeServiceClient) StopGather(ctx context.Context, in *StopGatherRequest, opts ...grpc.CallOption) (*meta.CommonResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(meta.CommonResponse)
	err := c.cc.Invoke(ctx, EdgeService_StopGather_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// EdgeServiceServer is the server API for EdgeService service.
// All implementations must embed UnimplementedEdgeServiceServer
// for forward compatibility.
type EdgeServiceServer interface {
	Ping(context.Context, *meta.Ping) (*meta.Ping, error)
	RegisterGateway(context.Context, *RegisterGatewayRequest) (*RegisterGatewayResponse, error)
	SetAgent(context.Context, *SetAgentRequest) (*SetAgentResponse, error)
	GetOpenapiDoc(context.Context, *GetOpenapiDocRequest) (*OpenapiDoc, error)
	GetAgentInfo(context.Context, *GetAgentInfoRequest) (*GetAgentInfoResponse, error)
	GetGatherData(context.Context, *GetDataRequest) (*DataMessage, error)
	GetDataStream(*GetDataRequest, grpc.ServerStreamingServer[DataChunk]) error
	SendHttpMethod(context.Context, *SendHttpMethodRequest) (*SendHttpMethodResponse, error)
	StartGather(context.Context, *StartGatherRequest) (*meta.CommonResponse, error)
	StopGather(context.Context, *StopGatherRequest) (*meta.CommonResponse, error)
	mustEmbedUnimplementedEdgeServiceServer()
}

// UnimplementedEdgeServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedEdgeServiceServer struct{}

func (UnimplementedEdgeServiceServer) Ping(context.Context, *meta.Ping) (*meta.Ping, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedEdgeServiceServer) RegisterGateway(context.Context, *RegisterGatewayRequest) (*RegisterGatewayResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RegisterGateway not implemented")
}
func (UnimplementedEdgeServiceServer) SetAgent(context.Context, *SetAgentRequest) (*SetAgentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetAgent not implemented")
}
func (UnimplementedEdgeServiceServer) GetOpenapiDoc(context.Context, *GetOpenapiDocRequest) (*OpenapiDoc, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetOpenapiDoc not implemented")
}
func (UnimplementedEdgeServiceServer) GetAgentInfo(context.Context, *GetAgentInfoRequest) (*GetAgentInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAgentInfo not implemented")
}
func (UnimplementedEdgeServiceServer) GetGatherData(context.Context, *GetDataRequest) (*DataMessage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetGatherData not implemented")
}
func (UnimplementedEdgeServiceServer) GetDataStream(*GetDataRequest, grpc.ServerStreamingServer[DataChunk]) error {
	return status.Errorf(codes.Unimplemented, "method GetDataStream not implemented")
}
func (UnimplementedEdgeServiceServer) SendHttpMethod(context.Context, *SendHttpMethodRequest) (*SendHttpMethodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendHttpMethod not implemented")
}
func (UnimplementedEdgeServiceServer) StartGather(context.Context, *StartGatherRequest) (*meta.CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StartGather not implemented")
}
func (UnimplementedEdgeServiceServer) StopGather(context.Context, *StopGatherRequest) (*meta.CommonResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopGather not implemented")
}
func (UnimplementedEdgeServiceServer) mustEmbedUnimplementedEdgeServiceServer() {}
func (UnimplementedEdgeServiceServer) testEmbeddedByValue()                     {}

// UnsafeEdgeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to EdgeServiceServer will
// result in compilation errors.
type UnsafeEdgeServiceServer interface {
	mustEmbedUnimplementedEdgeServiceServer()
}

func RegisterEdgeServiceServer(s grpc.ServiceRegistrar, srv EdgeServiceServer) {
	// If the following call pancis, it indicates UnimplementedEdgeServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&EdgeService_ServiceDesc, srv)
}

func _EdgeService_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(meta.Ping)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).Ping(ctx, req.(*meta.Ping))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_RegisterGateway_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterGatewayRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).RegisterGateway(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_RegisterGateway_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).RegisterGateway(ctx, req.(*RegisterGatewayRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_SetAgent_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetAgentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).SetAgent(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_SetAgent_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).SetAgent(ctx, req.(*SetAgentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_GetOpenapiDoc_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetOpenapiDocRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).GetOpenapiDoc(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_GetOpenapiDoc_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).GetOpenapiDoc(ctx, req.(*GetOpenapiDocRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_GetAgentInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAgentInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).GetAgentInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_GetAgentInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).GetAgentInfo(ctx, req.(*GetAgentInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_GetGatherData_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDataRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).GetGatherData(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_GetGatherData_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).GetGatherData(ctx, req.(*GetDataRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_GetDataStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(GetDataRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(EdgeServiceServer).GetDataStream(m, &grpc.GenericServerStream[GetDataRequest, DataChunk]{ServerStream: stream})
}

// This type alias is provided for backwards compatibility with existing code that references the prior non-generic stream type by name.
type EdgeService_GetDataStreamServer = grpc.ServerStreamingServer[DataChunk]

func _EdgeService_SendHttpMethod_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SendHttpMethodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).SendHttpMethod(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_SendHttpMethod_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).SendHttpMethod(ctx, req.(*SendHttpMethodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_StartGather_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StartGatherRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).StartGather(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_StartGather_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).StartGather(ctx, req.(*StartGatherRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _EdgeService_StopGather_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StopGatherRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EdgeServiceServer).StopGather(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: EdgeService_StopGather_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EdgeServiceServer).StopGather(ctx, req.(*StopGatherRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// EdgeService_ServiceDesc is the grpc.ServiceDesc for EdgeService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var EdgeService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lubricant.agent.edgeService",
	HandlerType: (*EdgeServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ping",
			Handler:    _EdgeService_Ping_Handler,
		},
		{
			MethodName: "registerGateway",
			Handler:    _EdgeService_RegisterGateway_Handler,
		},
		{
			MethodName: "setAgent",
			Handler:    _EdgeService_SetAgent_Handler,
		},
		{
			MethodName: "getOpenapiDoc",
			Handler:    _EdgeService_GetOpenapiDoc_Handler,
		},
		{
			MethodName: "getAgentInfo",
			Handler:    _EdgeService_GetAgentInfo_Handler,
		},
		{
			MethodName: "getGatherData",
			Handler:    _EdgeService_GetGatherData_Handler,
		},
		{
			MethodName: "sendHttpMethod",
			Handler:    _EdgeService_SendHttpMethod_Handler,
		},
		{
			MethodName: "startGather",
			Handler:    _EdgeService_StartGather_Handler,
		},
		{
			MethodName: "stopGather",
			Handler:    _EdgeService_StopGather_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "GetDataStream",
			Handler:       _EdgeService_GetDataStream_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "protobuf/agent/agent.proto",
}
