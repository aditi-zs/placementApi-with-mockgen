// Code generated by MockGen. DO NOT EDIT.
// Source: interface.go

// Package service is a generated GoMock package.
package service

import (
	context "context"
	reflect "reflect"

	entities "github.com/aditi-zs/Placement-API/entities"
	gomock "github.com/golang/mock/gomock"
	uuid "github.com/google/uuid"
)

// MockStudentSvc is a mock of StudentSvc interface.
type MockStudentSvc struct {
	ctrl     *gomock.Controller
	recorder *MockStudentSvcMockRecorder
}

// MockStudentSvcMockRecorder is the mock recorder for MockStudentSvc.
type MockStudentSvcMockRecorder struct {
	mock *MockStudentSvc
}

// NewMockStudentSvc creates a new mock instance.
func NewMockStudentSvc(ctrl *gomock.Controller) *MockStudentSvc {
	mock := &MockStudentSvc{ctrl: ctrl}
	mock.recorder = &MockStudentSvcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStudentSvc) EXPECT() *MockStudentSvcMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockStudentSvc) Create(ctx context.Context, stu *entities.Student) (entities.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, stu)
	ret0, _ := ret[0].(entities.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockStudentSvcMockRecorder) Create(ctx, stu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStudentSvc)(nil).Create), ctx, stu)
}

// Delete mocks base method.
func (m *MockStudentSvc) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockStudentSvcMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStudentSvc)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockStudentSvc) Get(ctx context.Context, name, branch, includeCompany string) ([]entities.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx, name, branch, includeCompany)
	ret0, _ := ret[0].([]entities.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockStudentSvcMockRecorder) Get(ctx, name, branch, includeCompany interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStudentSvc)(nil).Get), ctx, name, branch, includeCompany)
}

// GetByID mocks base method.
func (m *MockStudentSvc) GetByID(ctx context.Context, id uuid.UUID) (entities.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(entities.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockStudentSvcMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockStudentSvc)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *MockStudentSvc) Update(ctx context.Context, id uuid.UUID, stu *entities.Student) (entities.Student, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, stu)
	ret0, _ := ret[0].(entities.Student)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockStudentSvcMockRecorder) Update(ctx, id, stu interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockStudentSvc)(nil).Update), ctx, id, stu)
}

// MockCompanySvc is a mock of CompanySvc interface.
type MockCompanySvc struct {
	ctrl     *gomock.Controller
	recorder *MockCompanySvcMockRecorder
}

// MockCompanySvcMockRecorder is the mock recorder for MockCompanySvc.
type MockCompanySvcMockRecorder struct {
	mock *MockCompanySvc
}

// NewMockCompanySvc creates a new mock instance.
func NewMockCompanySvc(ctrl *gomock.Controller) *MockCompanySvc {
	mock := &MockCompanySvc{ctrl: ctrl}
	mock.recorder = &MockCompanySvcMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCompanySvc) EXPECT() *MockCompanySvcMockRecorder {
	return m.recorder
}

// Create mocks base method.
func (m *MockCompanySvc) Create(ctx context.Context, cmp entities.Company) (entities.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", ctx, cmp)
	ret0, _ := ret[0].(entities.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Create indicates an expected call of Create.
func (mr *MockCompanySvcMockRecorder) Create(ctx, cmp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockCompanySvc)(nil).Create), ctx, cmp)
}

// Delete mocks base method.
func (m *MockCompanySvc) Delete(ctx context.Context, id uuid.UUID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", ctx, id)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete.
func (mr *MockCompanySvcMockRecorder) Delete(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockCompanySvc)(nil).Delete), ctx, id)
}

// Get mocks base method.
func (m *MockCompanySvc) Get(ctx context.Context) ([]entities.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", ctx)
	ret0, _ := ret[0].([]entities.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get.
func (mr *MockCompanySvcMockRecorder) Get(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockCompanySvc)(nil).Get), ctx)
}

// GetByID mocks base method.
func (m *MockCompanySvc) GetByID(ctx context.Context, id uuid.UUID) (entities.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByID", ctx, id)
	ret0, _ := ret[0].(entities.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByID indicates an expected call of GetByID.
func (mr *MockCompanySvcMockRecorder) GetByID(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByID", reflect.TypeOf((*MockCompanySvc)(nil).GetByID), ctx, id)
}

// Update mocks base method.
func (m *MockCompanySvc) Update(ctx context.Context, id uuid.UUID, cmp entities.Company) (entities.Company, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Update", ctx, id, cmp)
	ret0, _ := ret[0].(entities.Company)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Update indicates an expected call of Update.
func (mr *MockCompanySvcMockRecorder) Update(ctx, id, cmp interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Update", reflect.TypeOf((*MockCompanySvc)(nil).Update), ctx, id, cmp)
}
