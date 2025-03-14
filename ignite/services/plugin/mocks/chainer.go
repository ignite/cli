// Code generated by mockery v2.46.0. DO NOT EDIT.

package mocks

import mock "github.com/stretchr/testify/mock"

// ChainerInterface is an autogenerated mock type for the Chainer type
type ChainerInterface struct {
	mock.Mock
}

type ChainerInterface_Expecter struct {
	mock *mock.Mock
}

func (_m *ChainerInterface) EXPECT() *ChainerInterface_Expecter {
	return &ChainerInterface_Expecter{mock: &_m.Mock}
}

// AppPath provides a mock function with given fields:
func (_m *ChainerInterface) AppPath() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for AppPath")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ChainerInterface_AppPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'AppPath'
type ChainerInterface_AppPath_Call struct {
	*mock.Call
}

// AppPath is a helper method to define mock.On call
func (_e *ChainerInterface_Expecter) AppPath() *ChainerInterface_AppPath_Call {
	return &ChainerInterface_AppPath_Call{Call: _e.mock.On("AppPath")}
}

func (_c *ChainerInterface_AppPath_Call) Run(run func()) *ChainerInterface_AppPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ChainerInterface_AppPath_Call) Return(_a0 string) *ChainerInterface_AppPath_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ChainerInterface_AppPath_Call) RunAndReturn(run func() string) *ChainerInterface_AppPath_Call {
	_c.Call.Return(run)
	return _c
}

// ConfigPath provides a mock function with given fields:
func (_m *ChainerInterface) ConfigPath() string {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ConfigPath")
	}

	var r0 string
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}

// ChainerInterface_ConfigPath_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ConfigPath'
type ChainerInterface_ConfigPath_Call struct {
	*mock.Call
}

// ConfigPath is a helper method to define mock.On call
func (_e *ChainerInterface_Expecter) ConfigPath() *ChainerInterface_ConfigPath_Call {
	return &ChainerInterface_ConfigPath_Call{Call: _e.mock.On("ConfigPath")}
}

func (_c *ChainerInterface_ConfigPath_Call) Run(run func()) *ChainerInterface_ConfigPath_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ChainerInterface_ConfigPath_Call) Return(_a0 string) *ChainerInterface_ConfigPath_Call {
	_c.Call.Return(_a0)
	return _c
}

func (_c *ChainerInterface_ConfigPath_Call) RunAndReturn(run func() string) *ChainerInterface_ConfigPath_Call {
	_c.Call.Return(run)
	return _c
}

// Home provides a mock function with given fields:
func (_m *ChainerInterface) Home() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for Home")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChainerInterface_Home_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'Home'
type ChainerInterface_Home_Call struct {
	*mock.Call
}

// Home is a helper method to define mock.On call
func (_e *ChainerInterface_Expecter) Home() *ChainerInterface_Home_Call {
	return &ChainerInterface_Home_Call{Call: _e.mock.On("Home")}
}

func (_c *ChainerInterface_Home_Call) Run(run func()) *ChainerInterface_Home_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ChainerInterface_Home_Call) Return(_a0 string, _a1 error) *ChainerInterface_Home_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChainerInterface_Home_Call) RunAndReturn(run func() (string, error)) *ChainerInterface_Home_Call {
	_c.Call.Return(run)
	return _c
}

// ID provides a mock function with given fields:
func (_m *ChainerInterface) ID() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for ID")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChainerInterface_ID_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'ID'
type ChainerInterface_ID_Call struct {
	*mock.Call
}

// ID is a helper method to define mock.On call
func (_e *ChainerInterface_Expecter) ID() *ChainerInterface_ID_Call {
	return &ChainerInterface_ID_Call{Call: _e.mock.On("ID")}
}

func (_c *ChainerInterface_ID_Call) Run(run func()) *ChainerInterface_ID_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ChainerInterface_ID_Call) Return(_a0 string, _a1 error) *ChainerInterface_ID_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChainerInterface_ID_Call) RunAndReturn(run func() (string, error)) *ChainerInterface_ID_Call {
	_c.Call.Return(run)
	return _c
}

// RPCPublicAddress provides a mock function with given fields:
func (_m *ChainerInterface) RPCPublicAddress() (string, error) {
	ret := _m.Called()

	if len(ret) == 0 {
		panic("no return value specified for RPCPublicAddress")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func() (string, error)); ok {
		return rf()
	}
	if rf, ok := ret.Get(0).(func() string); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ChainerInterface_RPCPublicAddress_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'RPCPublicAddress'
type ChainerInterface_RPCPublicAddress_Call struct {
	*mock.Call
}

// RPCPublicAddress is a helper method to define mock.On call
func (_e *ChainerInterface_Expecter) RPCPublicAddress() *ChainerInterface_RPCPublicAddress_Call {
	return &ChainerInterface_RPCPublicAddress_Call{Call: _e.mock.On("RPCPublicAddress")}
}

func (_c *ChainerInterface_RPCPublicAddress_Call) Run(run func()) *ChainerInterface_RPCPublicAddress_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run()
	})
	return _c
}

func (_c *ChainerInterface_RPCPublicAddress_Call) Return(_a0 string, _a1 error) *ChainerInterface_RPCPublicAddress_Call {
	_c.Call.Return(_a0, _a1)
	return _c
}

func (_c *ChainerInterface_RPCPublicAddress_Call) RunAndReturn(run func() (string, error)) *ChainerInterface_RPCPublicAddress_Call {
	_c.Call.Return(run)
	return _c
}

// NewChainerInterface creates a new instance of ChainerInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewChainerInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *ChainerInterface {
	mock := &ChainerInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
