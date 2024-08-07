// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/sub/internal/service/sub (interfaces: RecipientDeleter)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_deleter.go -package=mocks . RecipientDeleter
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "go.uber.org/mock/gomock"
)

// MockRecipientDeleter is a mock of RecipientDeleter interface.
type MockRecipientDeleter struct {
	ctrl     *gomock.Controller
	recorder *MockRecipientDeleterMockRecorder
}

// MockRecipientDeleterMockRecorder is the mock recorder for MockRecipientDeleter.
type MockRecipientDeleterMockRecorder struct {
	mock *MockRecipientDeleter
}

// NewMockRecipientDeleter creates a new mock instance.
func NewMockRecipientDeleter(ctrl *gomock.Controller) *MockRecipientDeleter {
	mock := &MockRecipientDeleter{ctrl: ctrl}
	mock.recorder = &MockRecipientDeleterMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRecipientDeleter) EXPECT() *MockRecipientDeleterMockRecorder {
	return m.recorder
}

// DeleteByEmail mocks base method.
func (m *MockRecipientDeleter) DeleteByEmail(arg0 context.Context, arg1 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteByEmail", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteByEmail indicates an expected call of DeleteByEmail.
func (mr *MockRecipientDeleterMockRecorder) DeleteByEmail(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteByEmail", reflect.TypeOf((*MockRecipientDeleter)(nil).DeleteByEmail), arg0, arg1)
}
