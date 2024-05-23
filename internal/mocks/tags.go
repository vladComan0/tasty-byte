// Code generated by MockGen. DO NOT EDIT.
// Source: ./internal/models/tags.go

// Package mock_models is a generated GoMock package.
package mock_models

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/vladComan0/tasty-byte/internal/models"
	transactions "github.com/vladComan0/tasty-byte/pkg/transactions"
)

// MockTagModelInterface is a mock of TagModelInterface interface.
type MockTagModelInterface struct {
	ctrl     *gomock.Controller
	recorder *MockTagModelInterfaceMockRecorder
}

// MockTagModelInterfaceMockRecorder is the mock recorder for MockTagModelInterface.
type MockTagModelInterfaceMockRecorder struct {
	mock *MockTagModelInterface
}

// NewMockTagModelInterface creates a new mock instance.
func NewMockTagModelInterface(ctrl *gomock.Controller) *MockTagModelInterface {
	mock := &MockTagModelInterface{ctrl: ctrl}
	mock.recorder = &MockTagModelInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockTagModelInterface) EXPECT() *MockTagModelInterfaceMockRecorder {
	return m.recorder
}

// GetByRecipeID mocks base method.
func (m *MockTagModelInterface) GetByRecipeID(tx transactions.Transaction, recipeID int) ([]*models.Tag, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRecipeID", tx, recipeID)
	ret0, _ := ret[0].([]*models.Tag)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRecipeID indicates an expected call of GetByRecipeID.
func (mr *MockTagModelInterfaceMockRecorder) GetByRecipeID(tx, recipeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRecipeID", reflect.TypeOf((*MockTagModelInterface)(nil).GetByRecipeID), tx, recipeID)
}

// InsertIfNotExists mocks base method.
func (m *MockTagModelInterface) InsertIfNotExists(tx transactions.Transaction, name string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertIfNotExists", tx, name)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertIfNotExists indicates an expected call of InsertIfNotExists.
func (mr *MockTagModelInterfaceMockRecorder) InsertIfNotExists(tx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertIfNotExists", reflect.TypeOf((*MockTagModelInterface)(nil).InsertIfNotExists), tx, name)
}