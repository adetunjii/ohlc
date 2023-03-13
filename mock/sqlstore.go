// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/adetunjii/ohlc/store (interfaces: Sqlstore)

// Package mock is a generated GoMock package.
package mock

import (
	reflect "reflect"

	store "github.com/adetunjii/ohlc/store"
	gomock "github.com/golang/mock/gomock"
)

// MockSqlstore is a mock of Sqlstore interface.
type MockSqlstore struct {
	ctrl     *gomock.Controller
	recorder *MockSqlstoreMockRecorder
}

// MockSqlstoreMockRecorder is the mock recorder for MockSqlstore.
type MockSqlstoreMockRecorder struct {
	mock *MockSqlstore
}

// NewMockSqlstore creates a new mock instance.
func NewMockSqlstore(ctrl *gomock.Controller) *MockSqlstore {
	mock := &MockSqlstore{ctrl: ctrl}
	mock.recorder = &MockSqlstoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSqlstore) EXPECT() *MockSqlstoreMockRecorder {
	return m.recorder
}

// PriceData mocks base method.
func (m *MockSqlstore) PriceData() store.PriceDataStore {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PriceData")
	ret0, _ := ret[0].(store.PriceDataStore)
	return ret0
}

// PriceData indicates an expected call of PriceData.
func (mr *MockSqlstoreMockRecorder) PriceData() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PriceData", reflect.TypeOf((*MockSqlstore)(nil).PriceData))
}
