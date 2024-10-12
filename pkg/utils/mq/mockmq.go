package mq

import (
	"reflect"

	"github.com/golang/mock/gomock"
)

// MockMq is a mock of Mq interface.
type MockMq[T any] struct {
	ctrl     *gomock.Controller
	recorder *MockMqMockRecorder[T]
}

// MockMqMockRecorder is the mock recorder for MockMq.
type MockMqMockRecorder[T any] struct {
	mock *MockMq[T]
}

// NewMockMq creates a new mock instance.
func NewMockMq[T any](ctrl *gomock.Controller) *MockMq[T] {
	mock := &MockMq[T]{ctrl: ctrl}
	mock.recorder = &MockMqMockRecorder[T]{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockMq[T]) EXPECT() *MockMqMockRecorder[T] {
	return m.recorder
}

// Close mocks base method.
func (m *MockMq[T]) Close() {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "Close")
}

// Close indicates an expected call of Close.
func (mr *MockMqMockRecorder[T]) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockMq)(nil).Close))
}

// GetPayLoad mocks base method.
func (m *MockMq[T]) GetPayLoad(sub <-chan T) T {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPayLoad", sub)
	ret0, _ := ret[0].(T)
	return ret0
}

// GetPayLoad indicates an expected call of GetPayLoad.
func (mr *MockMqMockRecorder[T]) GetPayLoad(sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPayLoad", reflect.TypeOf((*MockMq)(nil).GetPayLoad), sub)
}

// Publish mocks base method.
func (m *MockMq[T]) Publish(topic string, msg T) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", topic, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockMqMockRecorder[T]) Publish(topic, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockMq)(nil).Publish), topic, msg)
}

// PublishBytes mocks base method.
func (m *MockMq[T]) PublishBytes(topic string, msg []byte) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PublishBytes", topic, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// PublishBytes indicates an expected call of PublishBytes.
func (mr *MockMqMockRecorder[T]) PublishBytes(topic, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PublishBytes", reflect.TypeOf((*MockMq)(nil).PublishBytes), topic, msg)
}

// SetConditions mocks base method.
func (m *MockMq[T]) SetConditions(capacity int) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "SetConditions", capacity)
}

// SetConditions indicates an expected call of SetConditions.
func (mr *MockMqMockRecorder[T]) SetConditions(capacity interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetConditions", reflect.TypeOf((*MockMq)(nil).SetConditions), capacity)
}

// Subscribe mocks base method.
func (m *MockMq[T]) Subscribe(topic string) (<-chan T, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Subscribe", topic)
	ret0, _ := ret[0].(<-chan T)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Subscribe indicates an expected call of Subscribe.
func (mr *MockMqMockRecorder[T]) Subscribe(topic interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Subscribe", reflect.TypeOf((*MockMq)(nil).Subscribe), topic)
}

// Unsubscribe mocks base method.
func (m *MockMq[T]) Unsubscribe(topic string, sub <-chan T) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unsubscribe", topic, sub)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unsubscribe indicates an expected call of Unsubscribe.
func (mr *MockMqMockRecorder[T]) Unsubscribe(topic, sub interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unsubscribe", reflect.TypeOf((*MockMq)(nil).Unsubscribe), topic, sub)
}
