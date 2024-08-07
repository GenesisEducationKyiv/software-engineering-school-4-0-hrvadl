// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/service/cron (interfaces: RateGetter)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_rategetter.go -package=mocks . RateGetter
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	rate "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/rate"
	gomock "go.uber.org/mock/gomock"
)

// MockRateGetter is a mock of RateGetter interface.
type MockRateGetter struct {
	ctrl     *gomock.Controller
	recorder *MockRateGetterMockRecorder
}

// MockRateGetterMockRecorder is the mock recorder for MockRateGetter.
type MockRateGetterMockRecorder struct {
	mock *MockRateGetter
}

// NewMockRateGetter creates a new mock instance.
func NewMockRateGetter(ctrl *gomock.Controller) *MockRateGetter {
	mock := &MockRateGetter{ctrl: ctrl}
	mock.recorder = &MockRateGetterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateGetter) EXPECT() *MockRateGetterMockRecorder {
	return m.recorder
}

// Get mocks base method.
func (m *MockRateGetter) Get(arg0 context.Context) (*rate.Exchange, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0)
	ret0, _ := ret[0].(*rate.Exchange)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockRateGetterMockRecorder) Get(arg0 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockRateGetter)(nil).Get), arg0)
}
