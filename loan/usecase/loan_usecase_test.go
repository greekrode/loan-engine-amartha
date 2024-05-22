package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
	"github.com/greekrode/loan-engine-amartha/domain/mocks"
	loanUsecase "github.com/greekrode/loan-engine-amartha/loan/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type LoanUsecaseSuite struct {
	suite.Suite
	timeout time.Duration
}

func (s *LoanUsecaseSuite) SetupSuite() {
	s.timeout = 2 * time.Second
}

func (s *LoanUsecaseSuite) TestGetLoanDetails() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		loanID        uint
		setupMocks    func(*mocks.BorrowerRepository, *mocks.LoanRepository)
		expected      *dto.GetLoanDetailsResponse
		expectedError error
	}{
		{
			name:   "Success Get Loan Details",
			loanID: 1,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mlr.On("FindLoanByID", mock.Anything, uint(1)).Return(&domain.Loan{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					BorrowerID:        1,
					Principal:         100.00,
					InterestRate:      10.00,
					DurationWeeks:     52,
					OutstandingAmount: 1000.00,
					StartDate:         fixedTime,
					PaymentSchedules: []domain.PaymentSchedule{
						{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
								DeletedAt: gorm.DeletedAt{Valid: false},
							},
							DueAmount: 100.00,
							DueDate:   fixedTime,
							Paid:      false,
						},
					},
				}, nil)
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				}, nil)
			},
			expected: &dto.GetLoanDetailsResponse{
				Principal:         100.00,
				InterestRate:      10.00,
				OutstandingAmount: 1000.00,
				Duration:          52,
				StartDate:         fixedTime,
				CreatedAt:         fixedTime,
				Borrower: dto.GetBorrowerResponse{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					CreatedAt: fixedTime,
				},
				PaymentSchedule: []dto.GetPaymentScheduleResponse{
					{
						DueAmount: 100.00,
						DueDate:   fixedTime,
						Paid:      false,
					},
				},
			},
		},
		{
			name:   "Error Finding Loan",
			loanID: 2,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mlr.On("FindLoanByID", mock.Anything, uint(2)).Return(nil, errors.New("not found"))
			},
			expected:      &dto.GetLoanDetailsResponse{},
			expectedError: errors.New("not found"),
		},
		{
			name:   "Error Finding Borrower",
			loanID: 3,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository) {
				mlr.On("FindLoanByID", mock.Anything, uint(3)).Return(&domain.Loan{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					BorrowerID:        1,
					Principal:         100.00,
					InterestRate:      10.00,
					DurationWeeks:     52,
					OutstandingAmount: 1000.00,
					StartDate:         fixedTime,
					PaymentSchedules: []domain.PaymentSchedule{
						{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
								DeletedAt: gorm.DeletedAt{Valid: false},
							},
							DueAmount: 100.00,
							DueDate:   fixedTime,
							Paid:      false,
						},
					},
				}, nil)
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(nil, errors.New("not found"))
			},
			expected:      &dto.GetLoanDetailsResponse{},
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			mockBorrowerRepo := new(mocks.BorrowerRepository)
			mockPaymentScheduleRepo := new(mocks.PaymentScheduleRepository)
			mockLoanRepo := new(mocks.LoanRepository)
			mockTransactionmanager := new(mocks.TransactionManager)

			uc := loanUsecase.NewLoanUsecase(mockBorrowerRepo, mockPaymentScheduleRepo, mockLoanRepo, mockTransactionmanager, s.timeout)

			tt.setupMocks(mockBorrowerRepo, mockLoanRepo)
			result, err := uc.GetLoanDetails(context.TODO(), tt.loanID)
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

func (s *LoanUsecaseSuite) TestGetOutstandingAmount() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		loanID        uint
		setupMocks    func(*mocks.LoanRepository)
		expected      float64
		expectedError error
	}{
		{
			name:   "Success Get Outstanding Amount",
			loanID: 1,
			setupMocks: func(mlr *mocks.LoanRepository) {
				mlr.On("FindLoanByID", mock.Anything, uint(1)).Return(&domain.Loan{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					BorrowerID:        1,
					Principal:         100.00,
					InterestRate:      10.00,
					DurationWeeks:     52,
					OutstandingAmount: 1000.00,
					StartDate:         fixedTime,
					PaymentSchedules: []domain.PaymentSchedule{
						{
							Model: gorm.Model{
								ID:        1,
								CreatedAt: fixedTime,
								UpdatedAt: fixedTime,
								DeletedAt: gorm.DeletedAt{Valid: false},
							},
							DueAmount: 100,
							DueDate:   fixedTime,
							Paid:      false,
						},
					},
				}, nil)
			},
			expected:      1000.00,
			expectedError: nil,
		},
		{
			name:   "Loan Not Found",
			loanID: 2,
			setupMocks: func(mlr *mocks.LoanRepository) {
				mlr.On("FindLoanByID", mock.Anything, uint(2)).Return(nil, errors.New("not found"))
			},
			expected:      0.00,
			expectedError: errors.New("not found"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			mockBorrowerRepo := new(mocks.BorrowerRepository)
			mockPaymentScheduleRepo := new(mocks.PaymentScheduleRepository)
			mockLoanRepo := new(mocks.LoanRepository)
			mockTransactionmanager := new(mocks.TransactionManager)

			uc := loanUsecase.NewLoanUsecase(mockBorrowerRepo, mockPaymentScheduleRepo, mockLoanRepo, mockTransactionmanager, s.timeout)

			tt.setupMocks(mockLoanRepo)
			result, err := uc.GetOutstandingAmount(context.TODO(), tt.loanID)
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

func (s *LoanUsecaseSuite) TestCreateLoan() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		borrowerID    uint
		principal     float64
		interestRate  float64
		durationWeeks int32
		startDate     time.Time
		setupMocks    func(*mocks.BorrowerRepository, *mocks.LoanRepository, *mocks.PaymentScheduleRepository, *mocks.TransactionManager)
		expected      *dto.CreateLoanResponse
		expectedError error
	}{
		{
			name:          "Successful Creation",
			borrowerID:    1,
			principal:     1000.00,
			interestRate:  5.00,
			durationWeeks: 2,
			startDate:     fixedTime,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository, mtm *mocks.TransactionManager) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mlr.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil)
				mlr.On("UpdateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil)
				mpsr.On("BulkCreatePaymentSchedule", mock.Anything, mock.AnythingOfType("[]domain.PaymentSchedule"), mock.Anything).Return(nil)
			},
			expected: &dto.CreateLoanResponse{
				ID:                0,
				Principal:         1000.00,
				InterestRate:      5.00,
				Duration:          2,
				StartDate:         fixedTime,
				OutstandingAmount: 1001.92,
				PaymentSchedules: []dto.GetPaymentScheduleResponse{
					{
						ID:        0,
						DueAmount: 500.96,
						DueDate:   fixedTime.Add(7 * 24 * time.Hour),
						Paid:      false,
					},
					{
						ID:        0,
						DueAmount: 500.96,
						DueDate:   fixedTime.Add(14 * 24 * time.Hour),
						Paid:      false,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:          "Borrower Not Found",
			borrowerID:    2,
			principal:     500.00,
			interestRate:  5.00,
			durationWeeks: 52,
			startDate:     fixedTime,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository, mtm *mocks.TransactionManager) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(2)).Return(nil, errors.New("borrower not found"))
			},
			expected:      nil,
			expectedError: errors.New("borrower not found"),
		},
		{
			name:          "Error Creating Loan",
			borrowerID:    1,
			principal:     1000.00,
			interestRate:  5.00,
			durationWeeks: 2,
			startDate:     fixedTime,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository, mtm *mocks.TransactionManager) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mlr.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(errors.New("error creating loan"))
			},
			expected:      nil,
			expectedError: errors.New("error creating loan"),
		},
		{
			name:          "Error Creating Payment Schedules",
			borrowerID:    1,
			principal:     1000.00,
			interestRate:  5.00,
			durationWeeks: 2,
			startDate:     fixedTime,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository, mtm *mocks.TransactionManager) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mlr.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil)
				mlr.On("UpdateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil)
				mpsr.On("BulkCreatePaymentSchedule", mock.Anything, mock.AnythingOfType("[]domain.PaymentSchedule"), mock.Anything).Return(errors.New("error creating payment schedules"))
			},
			expected:      nil,
			expectedError: errors.New("error creating payment schedules"),
		},
		{
			name:          "Error Updating Loan",
			borrowerID:    1,
			principal:     1000.00,
			interestRate:  5.00,
			durationWeeks: 2,
			startDate:     fixedTime,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository, mtm *mocks.TransactionManager) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mlr.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil)
				mlr.On("UpdateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(errors.New("error updating loan"))
			},
			expected:      nil,
			expectedError: errors.New("error updating loan"),
		},
		{
			name:          "Error Beginning Transaction",
			borrowerID:    1,
			principal:     1000.00,
			interestRate:  5.00,
			durationWeeks: 2,
			startDate:     fixedTime,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository, mtm *mocks.TransactionManager) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{
					Error: errors.New("error beginning transaction"),
				})
			},
			expected:      nil,
			expectedError: errors.New("error beginning transaction"),
		},
		{
			name:          "Error Committing Transaction",
			borrowerID:    1,
			principal:     1000.00,
			interestRate:  5.00,
			durationWeeks: 2,
			startDate:     fixedTime,
			setupMocks: func(mbr *mocks.BorrowerRepository, mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository, mtm *mocks.TransactionManager) {
				mbr.On("FindBorrowerByID", mock.Anything, uint(1)).Return(&domain.Borrower{
					Model: gorm.Model{
						ID:        1,
						CreatedAt: fixedTime,
						UpdatedAt: fixedTime,
						DeletedAt: gorm.DeletedAt{Valid: false},
					},
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(errors.New("error committing transaction"))
				mtm.On("Rollback", mock.Anything).Return(nil)
				mlr.On("CreateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil)
				mlr.On("UpdateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil)
				mpsr.On("BulkCreatePaymentSchedule", mock.Anything, mock.AnythingOfType("[]domain.PaymentSchedule"), mock.Anything).Return(nil)
			},
			expected:      nil,
			expectedError: errors.New("error committing transaction"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			mockBorrowerRepo := new(mocks.BorrowerRepository)
			mockLoanRepo := new(mocks.LoanRepository)
			mockPaymentScheduleRepo := new(mocks.PaymentScheduleRepository)
			mockTransactionManager := new(mocks.TransactionManager)

			uc := loanUsecase.NewLoanUsecase(mockBorrowerRepo, mockPaymentScheduleRepo, mockLoanRepo, mockTransactionManager, s.timeout)

			tt.setupMocks(mockBorrowerRepo, mockLoanRepo, mockPaymentScheduleRepo, mockTransactionManager)
			result, err := uc.CreateLoan(context.TODO(), tt.borrowerID, tt.principal, tt.interestRate, tt.durationWeeks, tt.startDate)
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

func TestLoanUsecaseSuite(t *testing.T) {
	suite.Run(t, new(LoanUsecaseSuite))
}
