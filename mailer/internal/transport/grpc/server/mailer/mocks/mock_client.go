// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/transport/grpc/server/mailer (interfaces: Client)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_client.go -package=mocks . Client
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	mail "github.com/GenesisEducationKyiv/software-engineering-school-4-0-hrvadl/mailer/internal/storage/mail"
	gomock "go.uber.org/mock/gomock"
)

// MockClient is a mock of Client interface.
type MockClient struct {
	ctrl     *gomock.Controller
	recorder *MockClientMockRecorder
}

// MockClientMockRecorder is the mock recorder for MockClient.
type MockClientMockRecorder struct {
	mock *MockClient
}

// NewMockClient creates a new mock instance.
func NewMockClient(ctrl *gomock.Controller) *MockClient {
	mock := &MockClient{ctrl: ctrl}
	mock.recorder = &MockClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockClient) EXPECT() *MockClientMockRecorder {
	return m.recorder
}

// Send mocks base method.
func (m *MockClient) Send(arg0 context.Context, arg1 mail.Mail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Send", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Send indicates an expected call of Send.
func (mr *MockClientMockRecorder) Send(arg0, arg1 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Send", reflect.TypeOf((*MockClient)(nil).Send), arg0, arg1)
}
