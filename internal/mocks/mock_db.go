// Code generated by MockGen. DO NOT EDIT.
// Source: internal/repo/repo.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockRepository is a mock of Repository interface.
type MockRepository struct {
	ctrl     *gomock.Controller
	recorder *MockRepositoryMockRecorder
}

// MockRepositoryMockRecorder is the mock recorder for MockRepository.
type MockRepositoryMockRecorder struct {
	mock *MockRepository
}

// NewMockRepository creates a new mock instance.
func NewMockRepository(ctrl *gomock.Controller) *MockRepository {
	mock := &MockRepository{ctrl: ctrl}
	mock.recorder = &MockRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRepository) EXPECT() *MockRepositoryMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockRepository) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockRepositoryMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockRepository)(nil).Close))
}

// GetCounter mocks base method.
func (m *MockRepository) GetCounter(ctx context.Context, c string) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounter", ctx, c)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounter indicates an expected call of GetCounter.
func (mr *MockRepositoryMockRecorder) GetCounter(ctx, c interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounter", reflect.TypeOf((*MockRepository)(nil).GetCounter), ctx, c)
}

// GetCounterAll mocks base method.
func (m *MockRepository) GetCounterAll(ctx context.Context) (map[string]int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCounterAll", ctx)
	ret0, _ := ret[0].(map[string]int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCounterAll indicates an expected call of GetCounterAll.
func (mr *MockRepositoryMockRecorder) GetCounterAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCounterAll", reflect.TypeOf((*MockRepository)(nil).GetCounterAll), ctx)
}

// GetGauge mocks base method.
func (m *MockRepository) GetGauge(ctx context.Context, g string) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGauge", ctx, g)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGauge indicates an expected call of GetGauge.
func (mr *MockRepositoryMockRecorder) GetGauge(ctx, g interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGauge", reflect.TypeOf((*MockRepository)(nil).GetGauge), ctx, g)
}

// GetGaugeAll mocks base method.
func (m *MockRepository) GetGaugeAll(ctx context.Context) (map[string]float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGaugeAll", ctx)
	ret0, _ := ret[0].(map[string]float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGaugeAll indicates an expected call of GetGaugeAll.
func (mr *MockRepositoryMockRecorder) GetGaugeAll(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGaugeAll", reflect.TypeOf((*MockRepository)(nil).GetGaugeAll), ctx)
}

// Ping mocks base method.
func (m *MockRepository) Ping() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping")
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockRepositoryMockRecorder) Ping() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockRepository)(nil).Ping))
}

// UpdateCounter mocks base method.
func (m *MockRepository) UpdateCounter(ctx context.Context, c string, v int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateCounter", ctx, c, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateCounter indicates an expected call of UpdateCounter.
func (mr *MockRepositoryMockRecorder) UpdateCounter(ctx, c, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateCounter", reflect.TypeOf((*MockRepository)(nil).UpdateCounter), ctx, c, v)
}

// UpdateGauge mocks base method.
func (m *MockRepository) UpdateGauge(ctx context.Context, g string, v float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGauge", ctx, g, v)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGauge indicates an expected call of UpdateGauge.
func (mr *MockRepositoryMockRecorder) UpdateGauge(ctx, g, v interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGauge", reflect.TypeOf((*MockRepository)(nil).UpdateGauge), ctx, g, v)
}
