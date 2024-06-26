// Code generated by MockGen. DO NOT EDIT.
// Source: internal/models/ingredients.go

// Package mock_models is a generated GoMock package.
package mocks

import (
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	models "github.com/vladComan0/tasty-byte/internal/models"
	transactions "github.com/vladComan0/tasty-byte/pkg/transactions"
)

// MockIngredientModelInterface is a mock of IngredientModelInterface interface.
type MockIngredientModelInterface struct {
	ctrl     *gomock.Controller
	recorder *MockIngredientModelInterfaceMockRecorder
}

// MockIngredientModelInterfaceMockRecorder is the mock recorder for MockIngredientModelInterface.
type MockIngredientModelInterfaceMockRecorder struct {
	mock *MockIngredientModelInterface
}

// NewMockIngredientModelInterface creates a new mock instance.
func NewMockIngredientModelInterface(ctrl *gomock.Controller) *MockIngredientModelInterface {
	mock := &MockIngredientModelInterface{ctrl: ctrl}
	mock.recorder = &MockIngredientModelInterfaceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockIngredientModelInterface) EXPECT() *MockIngredientModelInterfaceMockRecorder {
	return m.recorder
}

// GetByRecipeID mocks base method.
func (m *MockIngredientModelInterface) GetByRecipeID(tx transactions.Transaction, recipeID int) ([]*models.FullIngredient, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetByRecipeID", tx, recipeID)
	ret0, _ := ret[0].([]*models.FullIngredient)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetByRecipeID indicates an expected call of GetByRecipeID.
func (mr *MockIngredientModelInterfaceMockRecorder) GetByRecipeID(tx, recipeID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetByRecipeID", reflect.TypeOf((*MockIngredientModelInterface)(nil).GetByRecipeID), tx, recipeID)
}

// InsertIfNotExists mocks base method.
func (m *MockIngredientModelInterface) InsertIfNotExists(tx transactions.Transaction, name string) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertIfNotExists", tx, name)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// InsertIfNotExists indicates an expected call of InsertIfNotExists.
func (mr *MockIngredientModelInterfaceMockRecorder) InsertIfNotExists(tx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertIfNotExists", reflect.TypeOf((*MockIngredientModelInterface)(nil).InsertIfNotExists), tx, name)
}
