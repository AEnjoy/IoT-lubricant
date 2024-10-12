package gateway

import (
	"context"
	"reflect"

	"github.com/AEnjoy/IoT-lubricant/protobuf/core"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc"
)

var _ core.CoreServiceClient = (*MockCoreServiceClient)(nil)

// MockCoreServiceClient is a mock of CoreServiceClient interface.
type MockCoreServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockCoreServiceClientMockRecorder
}

func (m *MockCoreServiceClient) Ping(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[core.Ping, core.Ping], error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCoreServiceClient) GetTask(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[core.Task, core.Task], error) {
	//TODO implement me
	panic("implement me")
}

func (m *MockCoreServiceClient) PushData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[core.Data, core.Data], error) {
	//TODO implement me
	panic("implement me")
}

// MockCoreServiceClientMockRecorder is the mock recorder for MockCoreServiceClient.
type MockCoreServiceClientMockRecorder struct {
	mock *MockCoreServiceClient
}

func (mr *MockCoreServiceClientMockRecorder) Ping(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[core.Ping, core.Ping], error) {
	//TODO implement me
	panic("implement me")
}

func (mr *MockCoreServiceClientMockRecorder) GetTask(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[core.Task, core.Task], error) {
	//TODO implement me
	panic("implement me")
}

func (mr *MockCoreServiceClientMockRecorder) PushData(ctx context.Context, opts ...grpc.CallOption) (grpc.BidiStreamingClient[core.Data, core.Data], error) {
	//TODO implement me
	panic("implement me")
}

// NewMockCoreServiceClient creates a new mock instance.
func NewMockCoreServiceClient(ctrl *gomock.Controller) *MockCoreServiceClient {
	mock := &MockCoreServiceClient{ctrl: ctrl}
	mock.recorder = &MockCoreServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoreServiceClient) EXPECT() *MockCoreServiceClientMockRecorder {
	return m.recorder
}

// PushMessageId mocks base method.
func (m *MockCoreServiceClient) PushMessageId(ctx context.Context, in *core.MessageIdInfo, opts ...grpc.CallOption) (*core.MessageIdInfo, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "PushMessageId", varargs...)
	ret0, _ := ret[0].(*core.MessageIdInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushMessageId indicates an expected call of PushMessageId.
func (mr *MockCoreServiceClientMockRecorder) PushMessageId(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushMessageId", reflect.TypeOf((*MockCoreServiceClient)(nil).PushMessageId), varargs...)
}

// MockCoreServiceServer is a mock of CoreServiceServer interface.
type MockCoreServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockCoreServiceServerMockRecorder
}

// MockCoreServiceServerMockRecorder is the mock recorder for MockCoreServiceServer.
type MockCoreServiceServerMockRecorder struct {
	mock *MockCoreServiceServer
}

// NewMockCoreServiceServer creates a new mock instance.
func NewMockCoreServiceServer(ctrl *gomock.Controller) *MockCoreServiceServer {
	mock := &MockCoreServiceServer{ctrl: ctrl}
	mock.recorder = &MockCoreServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCoreServiceServer) EXPECT() *MockCoreServiceServerMockRecorder {
	return m.recorder
}

// PushMessageId mocks base method.
func (m *MockCoreServiceServer) PushMessageId(arg0 context.Context, arg1 *core.MessageIdInfo) (*core.MessageIdInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PushMessageId", arg0, arg1)
	ret0, _ := ret[0].(*core.MessageIdInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// PushMessageId indicates an expected call of PushMessageId.
func (mr *MockCoreServiceServerMockRecorder) PushMessageId(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PushMessageId", reflect.TypeOf((*MockCoreServiceServer)(nil).PushMessageId), arg0, arg1)
}

// MockUnsafeCoreServiceServer is a mock of UnsafeCoreServiceServer interface.
type MockUnsafeCoreServiceServer struct {
	ctrl     *gomock.Controller
	recorder *MockUnsafeCoreServiceServerMockRecorder
}

// MockUnsafeCoreServiceServerMockRecorder is the mock recorder for MockUnsafeCoreServiceServer.
type MockUnsafeCoreServiceServerMockRecorder struct {
	mock *MockUnsafeCoreServiceServer
}

// NewMockUnsafeCoreServiceServer creates a new mock instance.
func NewMockUnsafeCoreServiceServer(ctrl *gomock.Controller) *MockUnsafeCoreServiceServer {
	mock := &MockUnsafeCoreServiceServer{ctrl: ctrl}
	mock.recorder = &MockUnsafeCoreServiceServerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUnsafeCoreServiceServer) EXPECT() *MockUnsafeCoreServiceServerMockRecorder {
	return m.recorder
}
