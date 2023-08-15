// Code generated by mockery v2.15.0. DO NOT EDIT.

package mocks

import (
	context "context"

	cnc "github.com/celestiaorg/go-cnc"

	mock "github.com/stretchr/testify/mock"
)

// CNCClientI is an autogenerated mock type for the CNCClientI type
type CNCClientI struct {
	mock.Mock
}

// NamespacedData provides a mock function with given fields: ctx, namespaceID, height
func (_m *CNCClientI) NamespacedData(ctx context.Context, namespaceID cnc.Namespace, height uint64) ([][]byte, error) {
	ret := _m.Called(ctx, namespaceID, height)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(context.Context, cnc.Namespace, uint64) [][]byte); ok {
		r0 = rf(ctx, namespaceID, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, cnc.Namespace, uint64) error); ok {
		r1 = rf(ctx, namespaceID, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NamespacedShares provides a mock function with given fields: ctx, namespaceID, height
func (_m *CNCClientI) NamespacedShares(ctx context.Context, namespaceID cnc.Namespace, height uint64) ([][]byte, error) {
	ret := _m.Called(ctx, namespaceID, height)

	var r0 [][]byte
	if rf, ok := ret.Get(0).(func(context.Context, cnc.Namespace, uint64) [][]byte); ok {
		r0 = rf(ctx, namespaceID, height)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([][]byte)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, cnc.Namespace, uint64) error); ok {
		r1 = rf(ctx, namespaceID, height)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// SubmitPFB provides a mock function with given fields: ctx, namespaceID, blob, fee, gasLimit
func (_m *CNCClientI) SubmitPFB(ctx context.Context, namespaceID cnc.Namespace, blob []byte, fee int64, gasLimit uint64) (*cnc.TxResponse, error) {
	ret := _m.Called(ctx, namespaceID, blob, fee, gasLimit)

	var r0 *cnc.TxResponse
	if rf, ok := ret.Get(0).(func(context.Context, cnc.Namespace, []byte, int64, uint64) *cnc.TxResponse); ok {
		r0 = rf(ctx, namespaceID, blob, fee, gasLimit)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*cnc.TxResponse)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, cnc.Namespace, []byte, int64, uint64) error); ok {
		r1 = rf(ctx, namespaceID, blob, fee, gasLimit)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewCNCClientI interface {
	mock.TestingT
	Cleanup(func())
}

// NewCNCClientI creates a new instance of CNCClientI. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewCNCClientI(t mockConstructorTestingTNewCNCClientI) *CNCClientI {
	mock := &CNCClientI{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}