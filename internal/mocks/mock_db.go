// Code generated by MockGen. DO NOT EDIT.
// Source: db.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockPersistentDB is a mock of PersistentDB interface.
type MockPersistentDB struct {
	ctrl     *gomock.Controller
	recorder *MockPersistentDBMockRecorder
}

// MockPersistentDBMockRecorder is the mock recorder for MockPersistentDB.
type MockPersistentDBMockRecorder struct {
	mock *MockPersistentDB
}

// NewMockPersistentDB creates a new mock instance.
func NewMockPersistentDB(ctrl *gomock.Controller) *MockPersistentDB {
	mock := &MockPersistentDB{ctrl: ctrl}
	mock.recorder = &MockPersistentDBMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPersistentDB) EXPECT() *MockPersistentDBMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockPersistentDB) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockPersistentDBMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockPersistentDB)(nil).Close))
}

// GetCounter mocks base method.
func (m *MockPersistentDB) GetCounter(ctx context.Context, c string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", ctx, c)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockPersistentDBMockRecorder) GetCounter(ctx, c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockPersistentDB)(nil).GetCounter), ctx, c)
}

// GetGauge mocks base method.
func (m *MockPersistentDB) GetGauge(ctx context.Context, g string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", ctx, g)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockPersistentDBMockRecorder) GetGauge(ctx, g interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockPersistentDB)(nil).GetGauge), ctx, g)
}

// Ping mocks base method.
func (m *MockPersistentDB) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockPersistentDBMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockPersistentDB)(nil).Ping))
}

// UpdateCounter mocks base method.
func (m *MockPersistentDB) UpdateCounter(ctx context.Context, c string, v int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounter", ctx, c, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounter indicates an expected call of UpdateCounter.
func (mr *MockPersistentDBMockRecorder) UpdateCounter(ctx, c, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockPersistentDB)(nil).UpdateCounter), ctx, c, v)
}

// UpdateGauge mocks base method.
func (m *MockPersistentDB) UpdateGauge(ctx context.Context, g string, v float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGauge", ctx, g, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGauge indicates an expected call of UpdateGauge.
func (mr *MockPersistentDBMockRecorder) UpdateGauge(ctx, g, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGauge", reflect.TypeOf((*MockPersistentDB)(nil).UpdateGauge), ctx, g, v)
}