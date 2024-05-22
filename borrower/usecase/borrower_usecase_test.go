package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	borrowerUsecase "github.com/greekrode/loan-engine-amartha/borrower/usecase"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

type BorrowerUsecaseSuite struct {
	suite.Suite
	timeout time.Duration
}

func (s *BorrowerUsecaseSuite) SetupSuite() {
	s.timeout = 2 * time.Second
}

func (s *BorrowerUsecaseSuite) TestIsDelinquent() {
	tests := []struct {
		name          string
		borrowerID    uint
		setupMocks    func(*mocks.BorrowerRepository, *mocks.LoanRepository)
		expected      bool
		expectedError error
	}{
		{
			name:       "Delinquent Borrower",
			borrowerID: 1,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{}, nil)
				mlr.On("GetLoansByBorrowerID", mock.Anything, uint(1)).Return([]domain.Loan{
					{
						PaymentSchedules: []domain.PaymentSchedule{
							{DueDate: time.Now().Add(-24 * time.Hour), Paid: false},
							{DueDate: time.Now().Add(-48 * time.Hour), Paid: false},
						},
					},
				}, nil)
			},
			expected:      true,
			expectedError: nil,
		},
		{
			name:       "Non-Delinquent Borrower",
			borrowerID: 2,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(2)).Return(&domain.Borrower{}, nil)
				mlr.On("GetLoansByBorrowerID", mock.Anything, uint(2)).Return([]domain.Loan{
					{
						PaymentSchedules: []domain.PaymentSchedule{
							{DueDate: time.Now().Add(-24 * time.Hour), Paid: true},
						},
					},
				}, nil)
			},
			expected:      false,
			expectedError: nil,
		},
		{
			name:       "Error Finding Borrower",
			borrowerID: 3,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(3)).Return(nil, errors.New("not found"))
			},
			expected:      false,
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			mockBorrowerRepo := new(mocks.BorrowerRepository)
			mockLoanRepo := new(mocks.LoanRepository)
			uc := borrowerUsecase.NewBorrowerUsecase(mockBorrowerRepo, mockLoanRepo, s.timeout)

			tt.setupMocks(mockBorrowerRepo, mockLoanRepo)
			result, err := uc.IsDelinquent(context.TODO(), tt.borrowerID)
			if tt.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), tt.expected, result)
			}
		})
	}
}

func (s *BorrowerUsecaseSuite) TestCreateBorrower() {
	tests := []struct {
		name          string
		setupMocks    func(*mocks.BorrowerRepository, *mocks.LoanRepository)
		expectedError error
	}{
		{
			name: "Successful Creation",
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mbr.On("CreateBorrower", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Borrower"), mock.AnythingOfType("*gorm.DB")).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "Failed Creation",
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mbr.On("CreateBorrower", mock.AnythingOfType("*context.timerCtx"), mock.AnythingOfType("*domain.Borrower"), mock.AnythingOfType("*gorm.DB")).Return(errors.New("creation failed"))
			},
			expectedError: errors.New("creation failed"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			mockBorrowerRepo := new(mocks.BorrowerRepository)
			mockLoanRepo := new(mocks.LoanRepository)
			uc := borrowerUsecase.NewBorrowerUsecase(mockBorrowerRepo, mockLoanRepo, s.timeout)

			tt.setupMocks(mockBorrowerRepo, mockLoanRepo)
			borrower, err := uc.CreateBorrower(context.Background())
			if tt.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tt.expectedError.Error(), err.Error())
				assert.Nil(s.T(), borrower)
			} else {
				assert.NoError(s.T(), err)
				assert.NotNil(s.T(), borrower)
			}
		})
	}
}

func TestBorrowerUsecaseSuite(t *testing.T) {
	suite.Run(t, new(BorrowerUsecaseSuite))
}
