// Code generated by mockery v2.42.3. DO NOT EDIT.

package mocks

import (
	context "context"

	domain "github.com/greekrode/loan-engine-amartha/domain"
	gorm "gorm.io/gorm"

	mock "github.com/stretchr/testify/mock"
)

// BorrowerRepository is an autogenerated mock type for the BorrowerRepository type
type BorrowerRepository struct {
	mock.Mock
}

// CreateBorrower provides a mock function with given fields: ctx, borrower, tx
func (_m *BorrowerRepository) CreateBorrower(ctx context.Context, borrower *domain.Borrower, tx *gorm.DB) error {
	ret := _m.Called(ctx, borrower, tx)

	if len(ret) == 0 {
		panic("no return value specified for CreateBorrower")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, *domain.Borrower, *gorm.DB) error); ok {
		r0 = rf(ctx, borrower, tx)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// FindBorrowerByID provides a mock function with given fields: ctx, borrowerID
func (_m *BorrowerRepository) FindBorrowerByID(ctx context.Context, borrowerID uint) (*domain.Borrower, error) {
	ret := _m.Called(ctx, borrowerID)

	if len(ret) == 0 {
		panic("no return value specified for FindBorrowerByID")
	}

	var r0 *domain.Borrower
	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, uint) (*domain.Borrower, error)); ok {
		return rf(ctx, borrowerID)
	}
	if rf, ok := ret.Get(0).(func(context.Context, uint) *domain.Borrower); ok {
		r0 = rf(ctx, borrowerID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*domain.Borrower)
		}
	}

	if rf, ok := ret.Get(1).(func(context.Context, uint) error); ok {
		r1 = rf(ctx, borrowerID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// NewBorrowerRepository creates a new instance of BorrowerRepository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewBorrowerRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *BorrowerRepository {
	mock := &BorrowerRepository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}