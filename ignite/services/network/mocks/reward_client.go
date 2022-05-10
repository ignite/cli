// Code generated by mockery v2.10.0. DO NOT EDIT.

package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"
	"github.com/tendermint/spn/x/reward/types"
	"google.golang.org/grpc"
)

// RewardClient is an autogenerated mock type for the RewardClient type
type RewardClient struct {
	mock.Mock
}

// Params provides a mock function with given fields: ctx, in, opts
func (_m *RewardClient) Params(ctx context.Context, in *types.QueryParamsRequest, opts ...grpc.CallOption) (*types.QueryParamsResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *types.QueryParamsResponse
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) *types.QueryParamsResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryParamsResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryParamsRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RewardPool provides a mock function with given fields: ctx, in, opts
func (_m *RewardClient) RewardPool(ctx context.Context, in *types.QueryGetRewardPoolRequest, opts ...grpc.CallOption) (*types.QueryGetRewardPoolResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *types.QueryGetRewardPoolResponse
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryGetRewardPoolRequest, ...grpc.CallOption) *types.QueryGetRewardPoolResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryGetRewardPoolResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryGetRewardPoolRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// RewardPoolAll provides a mock function with given fields: ctx, in, opts
func (_m *RewardClient) RewardPoolAll(ctx context.Context, in *types.QueryAllRewardPoolRequest, opts ...grpc.CallOption) (*types.QueryAllRewardPoolResponse, error) {
	_va := make([]interface{}, len(opts))
	for _i := range opts {
		_va[_i] = opts[_i]
	}
	var _ca []interface{}
	_ca = append(_ca, ctx, in)
	_ca = append(_ca, _va...)
	ret := _m.Called(_ca...)

	var r0 *types.QueryAllRewardPoolResponse
	if rf, ok := ret.Get(0).(func(context.Context, *types.QueryAllRewardPoolRequest, ...grpc.CallOption) *types.QueryAllRewardPoolResponse); ok {
		r0 = rf(ctx, in, opts...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*types.QueryAllRewardPoolResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, *types.QueryAllRewardPoolRequest, ...grpc.CallOption) error); ok {
		r1 = rf(ctx, in, opts...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
