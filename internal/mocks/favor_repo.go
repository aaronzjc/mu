// Code generated by mockery v2.14.1. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/aaronzjc/mu/internal/application/dto"
	mock "github.com/stretchr/testify/mock"

	model "github.com/aaronzjc/mu/internal/domain/model"
)

// FavorRepo is an autogenerated mock type for the FavorRepo type
type FavorRepo struct {
	mock.Mock
}

type FavorRepo_Expecter struct {
	mock *mock.Mock
}

func (_m *FavorRepo) EXPECT() *FavorRepo_Expecter {
	return &FavorRepo_Expecter{mock: &_m.Mock}
}

// CreateFavor provides a mock function with given fields: _a0, _a1
func (_m *FavorRepo) CreateFavor(_a0 context.Context, _a1 model.Favor) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Favor) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FavorRepo_CreateFavor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CreateFavor'
type FavorRepo_CreateFavor_Call struct {
	*mock.Call
}

// CreateFavor is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 model.Favor
func (_e *FavorRepo_Expecter) CreateFavor(_a0 interface{}, _a1 interface{}) *FavorRepo_CreateFavor_Call {
	return &FavorRepo_CreateFavor_Call{Call: _e.mock.On("CreateFavor", _a0, _a1)}
}

func (_c *FavorRepo_CreateFavor_Call) Run(run func(_a0 context.Context, _a1 model.Favor)) *FavorRepo_CreateFavor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.Favor))
	})
	return _c
}

func (_c *FavorRepo_CreateFavor_Call) Return(_a0 error) *FavorRepo_CreateFavor_Call {
	_c.Call.Return(_a0)
	return _c
}

// Del provides a mock function with given fields: _a0, _a1
func (_m *FavorRepo) Del(_a0 context.Context, _a1 model.Favor) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, model.Favor) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FavorRepo_Del_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Del'
type FavorRepo_Del_Call struct {
	*mock.Call
}

// Del is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 model.Favor
func (_e *FavorRepo_Expecter) Del(_a0 interface{}, _a1 interface{}) *FavorRepo_Del_Call {
	return &FavorRepo_Del_Call{Call: _e.mock.On("Del", _a0, _a1)}
}

func (_c *FavorRepo_Del_Call) Run(run func(_a0 context.Context, _a1 model.Favor)) *FavorRepo_Del_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(model.Favor))
	})
	return _c
}

func (_c *FavorRepo_Del_Call) Return(_a0 error) *FavorRepo_Del_Call {
	_c.Call.Return(_a0)
	return _c
}

// GetFavor provides a mock function with given fields: _a0, _a1
func (_m *FavorRepo) GetFavor(_a0 context.Context, _a1 *dto.Query) (model.Favor, error) {
	ret := _m.Called(_a0, _a1)

	var r0 model.Favor
	if rf, ok := ret.Get(0).(func(context.Context, *dto.Query) model.Favor); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(model.Favor)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dto.Query) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FavorRepo_GetFavor_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetFavor'
type FavorRepo_GetFavor_Call struct {
	*mock.Call
}

// GetFavor is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *dto.Query
func (_e *FavorRepo_Expecter) GetFavor(_a0 interface{}, _a1 interface{}) *FavorRepo_GetFavor_Call {
	return &FavorRepo_GetFavor_Call{Call: _e.mock.On("GetFavor", _a0, _a1)}
}

func (_c *FavorRepo_GetFavor_Call) Run(run func(_a0 context.Context, _a1 *dto.Query)) *FavorRepo_GetFavor_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.Query))
	})
	return _c
}

func (_c *FavorRepo_GetFavor_Call) Return(_a0 model.Favor, _a1 error) *FavorRepo_GetFavor_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// GetList provides a mock function with given fields: _a0, _a1
func (_m *FavorRepo) GetList(_a0 context.Context, _a1 *dto.Query) ([]model.Favor, error) {
	ret := _m.Called(_a0, _a1)

	var r0 []model.Favor
	if rf, ok := ret.Get(0).(func(context.Context, *dto.Query) []model.Favor); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]model.Favor)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *dto.Query) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// FavorRepo_GetList_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'GetList'
type FavorRepo_GetList_Call struct {
	*mock.Call
}

// GetList is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *dto.Query
func (_e *FavorRepo_Expecter) GetList(_a0 interface{}, _a1 interface{}) *FavorRepo_GetList_Call {
	return &FavorRepo_GetList_Call{Call: _e.mock.On("GetList", _a0, _a1)}
}

func (_c *FavorRepo_GetList_Call) Run(run func(_a0 context.Context, _a1 *dto.Query)) *FavorRepo_GetList_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.Query))
	})
	return _c
}

func (_c *FavorRepo_GetList_Call) Return(_a0 []model.Favor, _a1 error) *FavorRepo_GetList_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

// Sites provides a mock function with given fields: _a0, _a1
func (_m *FavorRepo) Sites(_a0 context.Context, _a1 *dto.Query) []string {
	ret := _m.Called(_a0, _a1)

	var r0 []string
	if rf, ok := ret.Get(0).(func(context.Context, *dto.Query) []string); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	return r0
}

// FavorRepo_Sites_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Sites'
type FavorRepo_Sites_Call struct {
	*mock.Call
}

// Sites is a helper method to define mock.On call
//   - _a0 context.Context
//   - _a1 *dto.Query
func (_e *FavorRepo_Expecter) Sites(_a0 interface{}, _a1 interface{}) *FavorRepo_Sites_Call {
	return &FavorRepo_Sites_Call{Call: _e.mock.On("Sites", _a0, _a1)}
}

func (_c *FavorRepo_Sites_Call) Run(run func(_a0 context.Context, _a1 *dto.Query)) *FavorRepo_Sites_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(*dto.Query))
	})
	return _c
}

func (_c *FavorRepo_Sites_Call) Return(_a0 []string) *FavorRepo_Sites_Call {
	_c.Call.Return(_a0)
	return _c
}

type mockConstructorTestingTNewFavorRepo interface {
	mock.TestingT
	Cleanup(func())
}

// NewFavorRepo creates a new instance of FavorRepo. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewFavorRepo(t mockConstructorTestingTNewFavorRepo) *FavorRepo {
	mock := &FavorRepo{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}