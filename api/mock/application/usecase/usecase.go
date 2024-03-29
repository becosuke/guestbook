// Code generated by MockGen. DO NOT EDIT.
// Source: internal/application/usecase/usecase.go

// Package mock_usecase is a generated GoMock package.
package mock_usecase

import (
	context "context"
	reflect "reflect"

	post "github.com/becosuke/guestbook/api/internal/domain/post"
	gomock "github.com/golang/mock/gomock"
)

// MockUsecase is a mock of Usecase interface.
type MockUsecase struct {
	ctrl     *gomock.Controller
	recorder *MockUsecaseMockRecorder
}

// MockUsecaseMockRecorder is the mock recorder for MockUsecase.
type MockUsecaseMockRecorder struct {
	mock *MockUsecase
}

// NewMockUsecase creates a new mock instance.
func NewMockUsecase(ctrl *gomock.Controller) *MockUsecase {
	mock := &MockUsecase{ctrl: ctrl}
	mock.recorder = &MockUsecaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUsecase) EXPECT() *MockUsecaseMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockUsecase) Create(arg0 context.Context, arg1 *post.Post) (*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", arg0, arg1)
	ret0, _ := ret[0].(*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockUsecaseMockRecorder) Create(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockUsecase)(nil).Create), arg0, arg1)
}

// Delete mocks base method.
func (m *MockUsecase) Delete(arg0 context.Context, arg1 *post.Serial) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockUsecaseMockRecorder) Delete(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockUsecase)(nil).Delete), arg0, arg1)
}

// Get mocks base method.
func (m *MockUsecase) Get(arg0 context.Context, arg1 *post.Serial) (*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", arg0, arg1)
	ret0, _ := ret[0].(*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockUsecaseMockRecorder) Get(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockUsecase)(nil).Get), arg0, arg1)
}

// Range mocks base method.
func (m *MockUsecase) Range(arg0 context.Context, arg1 *post.PageOption) ([]*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Range", arg0, arg1)
	ret0, _ := ret[0].([]*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Range indicates an expected call of Range.
func (mr *MockUsecaseMockRecorder) Range(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Range", reflect.TypeOf((*MockUsecase)(nil).Range), arg0, arg1)
}

// Update mocks base method.
func (m *MockUsecase) Update(arg0 context.Context, arg1 *post.Post) (*post.Post, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", arg0, arg1)
	ret0, _ := ret[0].(*post.Post)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockUsecaseMockRecorder) Update(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockUsecase)(nil).Update), arg0, arg1)
}
