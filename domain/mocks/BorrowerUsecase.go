// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	context "context"

	dto "github.com/greekrode/loan-engine-amartha/domain/dto"

	mock "github.com/stretchr/testify/mock"
)

// BorrowerUsecase is an autogenerated mock type for the BorrowerUsecase type
type BorrowerUsecase struct {
	mock.Mock
}

// CreateBorrower provides a mock function with given fields: ctx
func (_m *BorrowerUsecase) CreateBorrower(ctx context.Context) (*dto.CreateBorrowerResponse, error) {
	ret := _m.Called(ctx)

	if len(ret) == 0 {
		panic("no return value specified for CreateBorrower")
	}

	var r0 *dto.CreateBorrowerResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context) (*dto.CreateBorrowerResponse, error)); ok {
		return rf(ctx)
	}
	if rf, ok := ret.Get(0).(func(context.Context) *dto.CreateBorrowerResponse); ok {
		r0 = rf(ctx)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*dto.CreateBorrowerResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context) error); ok {
		r1 = rf(ctx)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// IsDelinquent provides a mock function with given fields: ctx, borrowerID
func (_m *BorrowerUsecase) IsDelinquent(ctx context.Context, borrowerID uint) (bool, error) {
	ret := _m.Called(ctx, borrowerID)

	if len(ret) == 0 {
		panic("no return value specified for IsDelinquent")
	}

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (bool, error)); ok {
		return rf(ctx, borrowerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) bool); ok {
		r0 = rf(ctx, borrowerID)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, borrowerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBorrowerUsecase creates a new instance of BorrowerUsecase. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBorrowerUsecase(t interface {
	mock.TestingT
	Cleanup(func())
}) *BorrowerUsecase {
	mock := &BorrowerUsecase{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
