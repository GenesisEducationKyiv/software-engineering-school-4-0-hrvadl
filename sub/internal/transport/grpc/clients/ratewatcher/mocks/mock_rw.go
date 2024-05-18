// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher (interfaces: RateWatcherServiceClient)
//
// Generated by this command:
//
//	mockgen -destination=./mocks/mock_rw.go -package=mocks github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher RateWatcherServiceClient
//

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	ratewatcher "github.com/hrvadl/btcratenotifier/protos/gen/go/v1/ratewatcher"
	gomock "go.uber.org/mock/gomock"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// MockRateWatcherServiceClient is a mock of RateWatcherServiceClient interface.
type MockRateWatcherServiceClient struct {
	ctrl     *gomock.Controller
	recorder *MockRateWatcherServiceClientMockRecorder
}

// MockRateWatcherServiceClientMockRecorder is the mock recorder for MockRateWatcherServiceClient.
type MockRateWatcherServiceClientMockRecorder struct {
	mock *MockRateWatcherServiceClient
}

// NewMockRateWatcherServiceClient creates a new mock instance.
func NewMockRateWatcherServiceClient(ctrl *gomock.Controller) *MockRateWatcherServiceClient {
	mock := &MockRateWatcherServiceClient{ctrl: ctrl}
	mock.recorder = &MockRateWatcherServiceClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockRateWatcherServiceClient) EXPECT() *MockRateWatcherServiceClientMockRecorder {
	return m.recorder
}

// GetRate mocks base method.
func (m *MockRateWatcherServiceClient) GetRate(arg0 context.Context, arg1 *emptypb.Empty, arg2 ...grpc.CallOption) (*ratewatcher.RateResponse, error) {
	m.ctrl.T.Helper()
	varargs := []any{arg0, arg1}
	for _, a := range arg2 {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetRate", varargs...)
	ret0, _ := ret[0].(*ratewatcher.RateResponse)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRate indicates an expected call of GetRate.
func (mr *MockRateWatcherServiceClientMockRecorder) GetRate(arg0, arg1 any, arg2 ...any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]any{arg0, arg1}, arg2...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRate", reflect.TypeOf((*MockRateWatcherServiceClient)(nil).GetRate), varargs...)
}
