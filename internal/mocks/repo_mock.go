// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ozonmp/est-water-api/internal/app/repo (interfaces: EventRepo)

// Package mocks is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/ozonmp/est-water-api/internal/model"
)

// MockEventRepo is a mock of EventRepo interface.
type MockEventRepo struct {
	ctrl     *gomock.Controller
	recorder *MockEventRepoMockRecorder
}

// MockEventRepoMockRecorder is the mock recorder for MockEventRepo.
type MockEventRepoMockRecorder struct {
	mock *MockEventRepo
}

// NewMockEventRepo creates a new mock instance.
func NewMockEventRepo(ctrl *gomock.Controller) *MockEventRepo {
	mock := &MockEventRepo{ctrl: ctrl}
	mock.recorder = &MockEventRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockEventRepo) EXPECT() *MockEventRepoMockRecorder {
	return m.recorder
}

// Add mocks base method.
func (m *MockEventRepo) Add(arg0 []model.WaterEvent) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Add", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Add indicates an expected call of Add.
func (mr *MockEventRepoMockRecorder) Add(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Add", reflect.TypeOf((*MockEventRepo)(nil).Add), arg0)
}

// Lock mocks base method.
func (m *MockEventRepo) Lock(arg0 uint64) ([]model.WaterEvent, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Lock", arg0)
	ret0, _ := ret[0].([]model.WaterEvent)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Lock indicates an expected call of Lock.
func (mr *MockEventRepoMockRecorder) Lock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Lock", reflect.TypeOf((*MockEventRepo)(nil).Lock), arg0)
}

// Remove mocks base method.
func (m *MockEventRepo) Remove(arg0 []uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Remove", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Remove indicates an expected call of Remove.
func (mr *MockEventRepoMockRecorder) Remove(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Remove", reflect.TypeOf((*MockEventRepo)(nil).Remove), arg0)
}

// Unlock mocks base method.
func (m *MockEventRepo) Unlock(arg0 []uint64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Unlock", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// Unlock indicates an expected call of Unlock.
func (mr *MockEventRepoMockRecorder) Unlock(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Unlock", reflect.TypeOf((*MockEventRepo)(nil).Unlock), arg0)
}
