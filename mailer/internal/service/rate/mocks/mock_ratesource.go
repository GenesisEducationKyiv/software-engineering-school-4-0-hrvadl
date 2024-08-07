// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/rate (interfaces: RateSource)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_ratesource.go -package=mocks . RateSource
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	rate "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
	gomock "go.uber.org/mock/gomock"
)

// MockRateSource is a mock of RateSource interface.
type MockRateSource struct {
	ctrl     *gomock.Controller
	recorder *MockRateSourceMockRecorder
}

// MockRateSourceMockRecorder is the mock recorder for MockRateSource.
type MockRateSourceMockRecorder struct {
	mock *MockRateSource
}

// NewMockRateSource creates a new mock instance.
func NewMockRateSource(ctrl *gomock.Controller) *MockRateSource {
	mock := &MockRateSource{ctrl: ctrl}
	mock.recorder = &MockRateSourceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateSource) EXPECT() *MockRateSourceMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRateSource) Get(arg0 context.Context) (*rate.Exchange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*rate.Exchange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRateSourceMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRateSource)(nil).Get), arg0)
}

// Replace mocks base method.
func (m *MockRateSource) Replace(arg0 context.Context, arg1 rate.Exchange) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Replace", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Replace indicates an expected call of Replace.
func (mr *MockRateSourceMockRecorder) Replace(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Replace", reflect.TypeOf((*MockRateSource)(nil).Replace), arg0, arg1)
}
