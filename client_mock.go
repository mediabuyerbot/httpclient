// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mediabuyerbot/httpclient (interfaces: Doer)

// Package httpclient is a generated GoMock package.
package httpclient

import (
	gomock "github.com/golang/mock/gomock"
	http "net/http"
	reflect "reflect"
)

// MockDoer is a mock of Doer interface
type MockDoer struct {
	ctrl     *gomock.Controller
	recorder *MockDoerMockRecorder
}

// MockDoerMockRecorder is the mock recorder for MockDoer
type MockDoerMockRecorder struct {
	mock *MockDoer
}

// NewMockDoer creates a new mock instance
func NewMockDoer(ctrl *gomock.Controller) *MockDoer {
	mock := &MockDoer{ctrl: ctrl}
	mock.recorder = &MockDoerMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDoer) EXPECT() *MockDoerMockRecorder {
	return m.recorder
}

// Do mocks base method
func (m *MockDoer) Do(arg0 *http.Request) (*http.Response, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Do", arg0)
	ret0, _ := ret[0].(*http.Response)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Do indicates an expected call of Do
func (mr *MockDoerMockRecorder) Do(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Do", reflect.TypeOf((*MockDoer)(nil).Do), arg0)
}
