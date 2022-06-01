// Code generated by mockery v2.12.2. DO NOT EDIT.

package mocks

import (
	context "context"

	cosmosclient "github.com/ignite-hq/cli/ignite/pkg/cosmosclient"

	mock "github.com/stretchr/testify/mock"

	testing "testing"
)

// TXsCollecter is an autogenerated mock type for the TXsCollecter type
type TXsCollecter struct {
	mock.Mock
}

type TXsCollecter_Expecter struct {
	mock *mock.Mock
}

func (_m *TXsCollecter) EXPECT() *TXsCollecter_Expecter {
	return &TXsCollecter_Expecter{mock: &_m.Mock}
}

// CollectTXs provides a mock function with given fields: ctx, fromHeight, tc
func (_m *TXsCollecter) CollectTXs(ctx context.Context, fromHeight int64, tc chan<- []cosmosclient.TX) error {
	ret := _m.Called(ctx, fromHeight, tc)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, chan<- []cosmosclient.TX) error); ok {
		r0 = rf(ctx, fromHeight, tc)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// TXsCollecter_CollectTXs_Call is a *mock.Call that shadows Run/Return methods with type explicit version for method 'CollectTXs'
type TXsCollecter_CollectTXs_Call struct {
	*mock.Call
}

// CollectTXs is a helper method to define mock.On call
//  - ctx context.Context
//  - fromHeight int64
//  - tc chan<- []cosmosclient.TX
func (_e *TXsCollecter_Expecter) CollectTXs(ctx interface{}, fromHeight interface{}, tc interface{}) *TXsCollecter_CollectTXs_Call {
	return &TXsCollecter_CollectTXs_Call{Call: _e.mock.On("CollectTXs", ctx, fromHeight, tc)}
}

func (_c *TXsCollecter_CollectTXs_Call) Run(run func(ctx context.Context, fromHeight int64, tc chan<- []cosmosclient.TX)) *TXsCollecter_CollectTXs_Call {
	_c.Call.Run(func(args mock.Arguments) {
		run(args[0].(context.Context), args[1].(int64), args[2].(chan<- []cosmosclient.TX))
	})
	return _c
}

func (_c *TXsCollecter_CollectTXs_Call) Return(_a0 error) *TXsCollecter_CollectTXs_Call {
	_c.Call.Return(_a0)
	return _c
}

// NewTXsCollecter creates a new instance of TXsCollecter. It also registers the testing.TB interface on the mock and a cleanup function to assert the mocks expectations.
func NewTXsCollecter(t testing.TB) *TXsCollecter {
	mock := &TXsCollecter{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
