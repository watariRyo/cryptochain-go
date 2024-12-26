// Code generated by MockGen. DO NOT EDIT.
// Source: ./web/usecase/usecase.go
//
// Generated by this command:
//
//	mockgen -source ./web/usecase/usecase.go -destination ./web/usecase/mock/usecase.go -package usecase
//

// Package usecase is a generated GoMock package.
package usecase

import (
	reflect "reflect"

	uuid "github.com/google/uuid"
	model "github.com/watariRyo/cryptochain-go/web/domain/model"
	gomock "go.uber.org/mock/gomock"
)

// MockUseCaseInterface is a mock of UseCaseInterface interface.
type MockUseCaseInterface struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseInterfaceMockRecorder
	isgomock struct{}
}

// MockUseCaseInterfaceMockRecorder is the mock recorder for MockUseCaseInterface.
type MockUseCaseInterfaceMockRecorder struct {
	mock *MockUseCaseInterface
}

// NewMockUseCaseInterface creates a new mock instance.
func NewMockUseCaseInterface(ctrl *gomock.Controller) *MockUseCaseInterface {
	mock := &MockUseCaseInterface{ctrl: ctrl}
	mock.recorder = &MockUseCaseInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCaseInterface) EXPECT() *MockUseCaseInterfaceMockRecorder {
	return m.recorder
}

// GetBlock mocks base method.
func (m *MockUseCaseInterface) GetBlock() []*model.Block {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetBlock")
	ret0, _ := ret[0].([]*model.Block)
	return ret0
}

// GetBlock indicates an expected call of GetBlock.
func (mr *MockUseCaseInterfaceMockRecorder) GetBlock() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetBlock", reflect.TypeOf((*MockUseCaseInterface)(nil).GetBlock))
}

// GetTransactionPool mocks base method.
func (m *MockUseCaseInterface) GetTransactionPool() map[uuid.UUID]*model.Transaction {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetTransactionPool")
	ret0, _ := ret[0].(map[uuid.UUID]*model.Transaction)
	return ret0
}

// GetTransactionPool indicates an expected call of GetTransactionPool.
func (mr *MockUseCaseInterfaceMockRecorder) GetTransactionPool() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetTransactionPool", reflect.TypeOf((*MockUseCaseInterface)(nil).GetTransactionPool))
}

// GetWalletInfo mocks base method.
func (m *MockUseCaseInterface) GetWalletInfo() (*model.WalletInfo, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWalletInfo")
	ret0, _ := ret[0].(*model.WalletInfo)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWalletInfo indicates an expected call of GetWalletInfo.
func (mr *MockUseCaseInterfaceMockRecorder) GetWalletInfo() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWalletInfo", reflect.TypeOf((*MockUseCaseInterface)(nil).GetWalletInfo))
}

// Mine mocks base method.
func (m *MockUseCaseInterface) Mine(payload string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Mine", payload)
	ret0, _ := ret[0].(error)
	return ret0
}

// Mine indicates an expected call of Mine.
func (mr *MockUseCaseInterfaceMockRecorder) Mine(payload any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Mine", reflect.TypeOf((*MockUseCaseInterface)(nil).Mine), payload)
}

// MineTransactions mocks base method.
func (m *MockUseCaseInterface) MineTransactions() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "MineTransactions")
	ret0, _ := ret[0].(error)
	return ret0
}

// MineTransactions indicates an expected call of MineTransactions.
func (mr *MockUseCaseInterfaceMockRecorder) MineTransactions() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "MineTransactions", reflect.TypeOf((*MockUseCaseInterface)(nil).MineTransactions))
}

// SyncWithRootState mocks base method.
func (m *MockUseCaseInterface) SyncWithRootState() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SyncWithRootState")
	ret0, _ := ret[0].(error)
	return ret0
}

// SyncWithRootState indicates an expected call of SyncWithRootState.
func (mr *MockUseCaseInterfaceMockRecorder) SyncWithRootState() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SyncWithRootState", reflect.TypeOf((*MockUseCaseInterface)(nil).SyncWithRootState))
}

// Transact mocks base method.
func (m *MockUseCaseInterface) Transact(req *model.Transact) (map[uuid.UUID]*model.Transaction, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Transact", req)
	ret0, _ := ret[0].(map[uuid.UUID]*model.Transaction)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Transact indicates an expected call of Transact.
func (mr *MockUseCaseInterfaceMockRecorder) Transact(req any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Transact", reflect.TypeOf((*MockUseCaseInterface)(nil).Transact), req)
}
