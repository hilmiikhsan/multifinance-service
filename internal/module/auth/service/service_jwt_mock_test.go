// Code generated by MockGen. DO NOT EDIT.
// Source: ports.go
//
// Generated by this command:
//
//	mockgen -source=ports.go -destination=../../internal/module/auth/service/service_jwt_mock_test.go -package=service
//

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	jwt_handler "github.com/hilmiikhsan/multifinance-service/pkg/jwt_handler"
	gomock "go.uber.org/mock/gomock"
)

// MockJWT is a mock of JWT interface.
type MockJWT struct {
	ctrl     *gomock.Controller
	recorder *MockJWTMockRecorder
	isgomock struct{}
}

// MockJWTMockRecorder is the mock recorder for MockJWT.
type MockJWTMockRecorder struct {
	mock *MockJWT
}

// NewMockJWT creates a new mock instance.
func NewMockJWT(ctrl *gomock.Controller) *MockJWT {
	mock := &MockJWT{ctrl: ctrl}
	mock.recorder = &MockJWTMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockJWT) EXPECT() *MockJWTMockRecorder {
	return m.recorder
}

// GenerateTokenString mocks base method.
func (m *MockJWT) GenerateTokenString(ctx context.Context, payload jwt_handler.CostumClaimsPayload) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateTokenString", ctx, payload)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GenerateTokenString indicates an expected call of GenerateTokenString.
func (mr *MockJWTMockRecorder) GenerateTokenString(ctx, payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateTokenString", reflect.TypeOf((*MockJWT)(nil).GenerateTokenString), ctx, payload)
}

// ParseTokenString mocks base method.
func (m *MockJWT) ParseTokenString(ctx context.Context, tokenString string) (*jwt_handler.CustomClaims, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ParseTokenString", ctx, tokenString)
	ret0, _ := ret[0].(*jwt_handler.CustomClaims)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ParseTokenString indicates an expected call of ParseTokenString.
func (mr *MockJWTMockRecorder) ParseTokenString(ctx, tokenString any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ParseTokenString", reflect.TypeOf((*MockJWT)(nil).ParseTokenString), ctx, tokenString)
}
