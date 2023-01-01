// Code generated by MockGen. DO NOT EDIT.
// Source: internal/domain/post/generator.go

// Package post is a generated GoMock package.
package post

import (
	reflect "reflect"

	post "github.com/becosuke/guestbook/api/internal/domain/post"
	gomock "github.com/golang/mock/gomock"
)

// MockGenerator is a mock of Generator interface.
type MockGenerator struct {
	ctrl     *gomock.Controller
	recorder *MockGeneratorMockRecorder
}

// MockGeneratorMockRecorder is the mock recorder for MockGenerator.
type MockGeneratorMockRecorder struct {
	mock *MockGenerator
}

// NewMockGenerator creates a new mock instance.
func NewMockGenerator(ctrl *gomock.Controller) *MockGenerator {
	mock := &MockGenerator{ctrl: ctrl}
	mock.recorder = &MockGeneratorMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockGenerator) EXPECT() *MockGeneratorMockRecorder {
	return m.recorder
}

// GenerateSerial mocks base method.
func (m *MockGenerator) GenerateSerial() *post.Serial {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GenerateSerial")
	ret0, _ := ret[0].(*post.Serial)
	return ret0
}

// GenerateSerial indicates an expected call of GenerateSerial.
func (mr *MockGeneratorMockRecorder) GenerateSerial() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GenerateSerial", reflect.TypeOf((*MockGenerator)(nil).GenerateSerial))
}
