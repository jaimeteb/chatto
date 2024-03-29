// Code generated by MockGen. DO NOT EDIT.
// Source: sql.go

// Package mocksql is a generated GoMock package.
package mocksql

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	gorm "gorm.io/gorm"
)

// MockDBClient is a mock of DBClient interface.
type MockDBClient struct {
	ctrl     *gomock.Controller
	recorder *MockDBClientMockRecorder
}

// MockDBClientMockRecorder is the mock recorder for MockDBClient.
type MockDBClientMockRecorder struct {
	mock *MockDBClient
}

// NewMockDBClient creates a new mock instance.
func NewMockDBClient(ctrl *gomock.Controller) *MockDBClient {
	mock := &MockDBClient{ctrl: ctrl}
	mock.recorder = &MockDBClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBClient) EXPECT() *MockDBClientMockRecorder {
	return m.recorder
}

// First mocks base method.
func (m *MockDBClient) First(arg0 interface{}, arg1 ...interface{}) *gorm.DB {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "First", varargs...)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// First indicates an expected call of First.
func (mr *MockDBClientMockRecorder) First(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "First", reflect.TypeOf((*MockDBClient)(nil).First), varargs...)
}

// Save mocks base method.
func (m *MockDBClient) Save(arg0 interface{}) *gorm.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Save", arg0)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Save indicates an expected call of Save.
func (mr *MockDBClientMockRecorder) Save(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Save", reflect.TypeOf((*MockDBClient)(nil).Save), arg0)
}

// Where mocks base method.
func (m *MockDBClient) Where(arg0 interface{}, arg1 ...interface{}) *gorm.DB {
	m.ctrl.T.Helper()
	varargs := []interface{}{arg0}
	for _, a := range arg1 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "Where", varargs...)
	ret0, _ := ret[0].(*gorm.DB)
	return ret0
}

// Where indicates an expected call of Where.
func (mr *MockDBClientMockRecorder) Where(arg0 interface{}, arg1 ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{arg0}, arg1...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Where", reflect.TypeOf((*MockDBClient)(nil).Where), varargs...)
}
