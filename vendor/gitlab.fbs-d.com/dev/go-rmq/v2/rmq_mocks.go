// Code generated by MockGen. DO NOT EDIT.
// Source: rmq.go

// Package rmq is a generated GoMock package.
package rmq

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	amqp "github.com/streadway/amqp"
)

// MockRmq is a mock of Rmq interface.
type MockRmq struct {
	ctrl     *gomock.Controller
	recorder *MockRmqMockRecorder
}

// MockRmqMockRecorder is the mock recorder for MockRmq.
type MockRmqMockRecorder struct {
	mock *MockRmq
}

// NewMockRmq creates a new mock instance.
func NewMockRmq(ctrl *gomock.Controller) *MockRmq {
	mock := &MockRmq{ctrl: ctrl}
	mock.recorder = &MockRmqMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRmq) EXPECT() *MockRmqMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRmq) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRmqMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRmq)(nil).Close))
}

// Connect mocks base method.
func (m *MockRmq) Connect(arg0 func() error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Connect", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Connect indicates an expected call of Connect.
func (mr *MockRmqMockRecorder) Connect(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Connect", reflect.TypeOf((*MockRmq)(nil).Connect), arg0)
}

// ConnectionIsInitialize mocks base method.
func (m *MockRmq) IsConnInitialized() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsConnInitialized")
	ret0, _ := ret[0].(bool)
	return ret0
}

// ConnectionIsInitialize indicates an expected call of ConnectionIsInitialize.
func (mr *MockRmqMockRecorder) ConnectionIsInitialize() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsConnInitialized", reflect.TypeOf((*MockRmq)(nil).IsConnInitialized))
}

// Consume mocks base method.
func (m *MockRmq) Consume(queueName string, worker WorkerFunc, skipErrors bool) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Consume", queueName, worker, skipErrors)
	ret0, _ := ret[0].(error)
	return ret0
}

// Consume indicates an expected call of Consume.
func (mr *MockRmqMockRecorder) Consume(queueName, worker, skipErrors interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Consume", reflect.TypeOf((*MockRmq)(nil).Consume), queueName, worker, skipErrors)
}

// ConsumeBatch mocks base method.
func (m *MockRmq) ConsumeBatch(exchangeName, routingKey, queueName, deadQueueName string, worker WorkerBatchFunc, trials int) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeBatch", exchangeName, routingKey, queueName, deadQueueName, worker, trials)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumeBatch indicates an expected call of ConsumeBatch.
func (mr *MockRmqMockRecorder) ConsumeBatch(exchangeName, routingKey, queueName, deadQueueName, worker, trials interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeBatch", reflect.TypeOf((*MockRmq)(nil).ConsumeBatch), exchangeName, routingKey, queueName, deadQueueName, worker, trials)
}

// ConsumeRpc mocks base method.
func (m *MockRmq) ConsumeRpc(queueName string, worker RpcWorkerFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeRpc", queueName, worker)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumeRpc indicates an expected call of ConsumeRpc.
func (mr *MockRmqMockRecorder) ConsumeRpc(queueName, worker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeRpc", reflect.TypeOf((*MockRmq)(nil).ConsumeRpc), queueName, worker)
}

// ConsumeTrialDelay mocks base method.
func (m *MockRmq) ConsumeTrialDelay(delay, trials int64, exchangeName, routingKey, queueName string, worker WorkerFunc) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ConsumeTrialDelay", delay, trials, exchangeName, routingKey, queueName, worker)
	ret0, _ := ret[0].(error)
	return ret0
}

// ConsumeTrialDelay indicates an expected call of ConsumeTrialDelay.
func (mr *MockRmqMockRecorder) ConsumeTrialDelay(delay, trials, exchangeName, routingKey, queueName, worker interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ConsumeTrialDelay", reflect.TypeOf((*MockRmq)(nil).ConsumeTrialDelay), delay, trials, exchangeName, routingKey, queueName, worker)
}

// ExchangeBind mocks base method.
func (m *MockRmq) ExchangeBind(from, to, key string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeBind", from, to, key)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangeBind indicates an expected call of ExchangeBind.
func (mr *MockRmqMockRecorder) ExchangeBind(from, to, key interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeBind", reflect.TypeOf((*MockRmq)(nil).ExchangeBind), from, to, key)
}

// ExchangeDeclare mocks base method.
func (m *MockRmq) ExchangeDeclare(exchangeName, exchangeType string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangeDeclare", exchangeName, exchangeType)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangeDeclare indicates an expected call of ExchangeDeclare.
func (mr *MockRmqMockRecorder) ExchangeDeclare(exchangeName, exchangeType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangeDeclare", reflect.TypeOf((*MockRmq)(nil).ExchangeDeclare), exchangeName, exchangeType)
}

// ExchangePublish mocks base method.
func (m *MockRmq) ExchangePublish(ctx context.Context, exchangeName, exchangeType, routingKey string, msg *amqp.Publishing) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ExchangePublish", ctx, exchangeName, exchangeType, routingKey, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// ExchangePublish indicates an expected call of ExchangePublish.
func (mr *MockRmqMockRecorder) ExchangePublish(ctx, exchangeName, exchangeType, routingKey, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExchangePublish", reflect.TypeOf((*MockRmq)(nil).ExchangePublish), ctx, exchangeName, exchangeType, routingKey, msg)
}

// IsClosed mocks base method.
func (m *MockRmq) IsClosed() bool {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "IsClosed")
	ret0, _ := ret[0].(bool)
	return ret0
}

// IsClosed indicates an expected call of IsClosed.
func (mr *MockRmqMockRecorder) IsClosed() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "IsClosed", reflect.TypeOf((*MockRmq)(nil).IsClosed))
}

// Publish mocks base method.
func (m *MockRmq) Publish(ctx context.Context, queueName string, msg *amqp.Publishing) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Publish", ctx, queueName, msg)
	ret0, _ := ret[0].(error)
	return ret0
}

// Publish indicates an expected call of Publish.
func (mr *MockRmqMockRecorder) Publish(ctx, queueName, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Publish", reflect.TypeOf((*MockRmq)(nil).Publish), ctx, queueName, msg)
}

// QueueBind mocks base method.
func (m *MockRmq) QueueBind(name, key, exchange string, args amqp.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueBind", name, key, exchange, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueueBind indicates an expected call of QueueBind.
func (mr *MockRmqMockRecorder) QueueBind(name, key, exchange, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueBind", reflect.TypeOf((*MockRmq)(nil).QueueBind), name, key, exchange, args)
}

// QueueDeclare mocks base method.
func (m *MockRmq) QueueDeclare(name string, args amqp.Table) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "QueueDeclare", name, args)
	ret0, _ := ret[0].(error)
	return ret0
}

// QueueDeclare indicates an expected call of QueueDeclare.
func (mr *MockRmqMockRecorder) QueueDeclare(name, args interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "QueueDeclare", reflect.TypeOf((*MockRmq)(nil).QueueDeclare), name, args)
}

// Rpc mocks base method.
func (m *MockRmq) Rpc(ctx context.Context, queueName string, msg *amqp.Publishing) (*amqp.Delivery, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rpc", ctx, queueName, msg)
	ret0, _ := ret[0].(*amqp.Delivery)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Rpc indicates an expected call of Rpc.
func (mr *MockRmqMockRecorder) Rpc(ctx, queueName, msg interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rpc", reflect.TypeOf((*MockRmq)(nil).Rpc), ctx, queueName, msg)
}

// State mocks base method.
func (m *MockRmq) State() *State {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "State")
	ret0, _ := ret[0].(*State)
	return ret0
}

// State indicates an expected call of State.
func (mr *MockRmqMockRecorder) State() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "State", reflect.TypeOf((*MockRmq)(nil).State))
}
