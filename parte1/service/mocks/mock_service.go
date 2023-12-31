// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package mock_service is a generated GoMock package.
package mock_service

import (
	service "TestAmalga/parte1/service"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// CalcularResumen mocks base method.
func (m *MockService) CalcularResumen(fecha, dias string) (*service.Resumen, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CalcularResumen", fecha, dias)
	ret0, _ := ret[0].(*service.Resumen)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CalcularResumen indicates an expected call of CalcularResumen.
func (mr *MockServiceMockRecorder) CalcularResumen(fecha, dias interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CalcularResumen", reflect.TypeOf((*MockService)(nil).CalcularResumen), fecha, dias)
}
