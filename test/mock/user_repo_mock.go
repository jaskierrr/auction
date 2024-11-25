// Code generated by MockGen. DO NOT EDIT.
// Source: ./user_repository.go

// Package mock is a generated GoMock package.
package mock

import (
	context "context"
	entities "main/internal/entities"
	grpc "main/pkg/grpc"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	pgx "github.com/jackc/pgx/v5"
)

// MockUserRepo is a mock of UserRepo interface.
type MockUserRepo struct {
	ctrl     *gomock.Controller
	recorder *MockUserRepoMockRecorder
}

// MockUserRepoMockRecorder is the mock recorder for MockUserRepo.
type MockUserRepoMockRecorder struct {
	mock *MockUserRepo
}

// NewMockUserRepo creates a new mock instance.
func NewMockUserRepo(ctrl *gomock.Controller) *MockUserRepo {
	mock := &MockUserRepo{ctrl: ctrl}
	mock.recorder = &MockUserRepoMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUserRepo) EXPECT() *MockUserRepoMockRecorder {
	return m.recorder
}

// CreateUser mocks base method.
func (m *MockUserRepo) CreateUser(ctx context.Context, in *grpc.CreateUserRequest) (entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateUser", ctx, in)
	ret0, _ := ret[0].(entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateUser indicates an expected call of CreateUser.
func (mr *MockUserRepoMockRecorder) CreateUser(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateUser", reflect.TypeOf((*MockUserRepo)(nil).CreateUser), ctx, in)
}

// GetUser mocks base method.
func (m *MockUserRepo) GetUser(ctx context.Context, in *grpc.GetUserRequest) (entities.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUser", ctx, in)
	ret0, _ := ret[0].(entities.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUser indicates an expected call of GetUser.
func (mr *MockUserRepoMockRecorder) GetUser(ctx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUser", reflect.TypeOf((*MockUserRepo)(nil).GetUser), ctx, in)
}

// PlaceBidWriteTransaction mocks base method.
func (m *MockUserRepo) PlaceBidWriteTransaction(ctx context.Context, tx pgx.Tx, in *grpc.DepositBalanceRequest) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "PlaceBidWriteTransaction", ctx, tx, in)
	ret0, _ := ret[0].(error)
	return ret0
}

// PlaceBidWriteTransaction indicates an expected call of PlaceBidWriteTransaction.
func (mr *MockUserRepoMockRecorder) PlaceBidWriteTransaction(ctx, tx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "PlaceBidWriteTransaction", reflect.TypeOf((*MockUserRepo)(nil).PlaceBidWriteTransaction), ctx, tx, in)
}

// StartTx mocks base method.
func (m *MockUserRepo) StartTx(ctx context.Context) (pgx.Tx, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "StartTx", ctx)
	ret0, _ := ret[0].(pgx.Tx)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// StartTx indicates an expected call of StartTx.
func (mr *MockUserRepoMockRecorder) StartTx(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "StartTx", reflect.TypeOf((*MockUserRepo)(nil).StartTx), ctx)
}

// UpdateBalance mocks base method.
func (m *MockUserRepo) UpdateBalance(ctx context.Context, tx pgx.Tx, in *grpc.DepositBalanceRequest) (string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateBalance", ctx, tx, in)
	ret0, _ := ret[0].(string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateBalance indicates an expected call of UpdateBalance.
func (mr *MockUserRepoMockRecorder) UpdateBalance(ctx, tx, in interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateBalance", reflect.TypeOf((*MockUserRepo)(nil).UpdateBalance), ctx, tx, in)
}
