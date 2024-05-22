package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
	"github.com/greekrode/loan-engine-amartha/domain/mocks"
	paymentUsecase "github.com/greekrode/loan-engine-amartha/payment/usecase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type PaymentUsecaseSuite struct {
	suite.Suite
	timeout time.Duration
}

func (s *PaymentUsecaseSuite) SetupSuite() {
	s.timeout = 2 * time.Second
}

func (s *PaymentUsecaseSuite) TestRequestPayment() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		loanID        uint
		setupMocks    func(*mocks.LoanRepository, *mocks.PaymentScheduleRepository)
		expected      *dto.RequestPaymentResponse
		expectedError error
	}{
		{
			name:   "Success Request Payment",
			loanID: 1,
			setupMocks: func(mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository) {
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
				mpsr.On("GetUnpaidPaymentSchedulesByLoanID", mock.Anything, uint(1), mock.Anything).Return([]domain.PaymentSchedule{
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
				}, nil)
			},
			expected: &dto.RequestPaymentResponse{
				TotalDue: 100.00,
				PaymentSchedules: []dto.GetPaymentScheduleResponse{
					{
						ID:        1,
						DueAmount: 100.00,
						DueDate:   fixedTime,
					},
				},
			},
			expectedError: nil,
		},
		{
			name:   "Loan Not Found",
			loanID: 1,
			setupMocks: func(mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository) {
				mlr.On("FindLoanByID", mock.Anything, uint(1)).Return(nil, errors.New("loan not found"))
			},
			expected:      nil,
			expectedError: errors.New("loan not found"),
		},
		{
			name:   "Error Finding Unpaid Payment Schedules",
			loanID: 1,
			setupMocks: func(mlr *mocks.LoanRepository, mpsr *mocks.PaymentScheduleRepository) {
				mlr.On("FindLoanByID", mock.Anything, uint(1)).Return(&domain.Loan{}, nil)
				mpsr.On("GetUnpaidPaymentSchedulesByLoanID", mock.Anything, uint(1), mock.Anything).Return(nil, errors.New("error finding unpaid payment schedules"))
			},
			expected:      nil,
			expectedError: errors.New("error finding unpaid payment schedules"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			mockPaymentRepo := new(mocks.PaymentRepository)
			mockPaymentScheduleRepo := new(mocks.PaymentScheduleRepository)
			mockLoanRepo := new(mocks.LoanRepository)
			mockTransactionmanager := new(mocks.TransactionManager)

			uc := paymentUsecase.NewPaymentUsecase(mockPaymentRepo, mockPaymentScheduleRepo, mockLoanRepo, mockTransactionmanager, s.timeout)

			tt.setupMocks(mockLoanRepo, mockPaymentScheduleRepo)
			result, err := uc.RequestPayment(context.TODO(), tt.loanID)
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

func (s *PaymentUsecaseSuite) TestMakePayment() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name          string
		loanID        uint
		paymentIDs    []uint
		amount        float64
		setupMocks    func(*mocks.PaymentRepository, *mocks.PaymentScheduleRepository, *mocks.LoanRepository, *mocks.TransactionManager)
		expectedError error
	}{
		{
			name:       "Successful Payment",
			loanID:     1,
			paymentIDs: []uint{1},
			amount:     100.00,
			setupMocks: func(mpr *mocks.PaymentRepository, mpsr *mocks.PaymentScheduleRepository, mlr *mocks.LoanRepository, mtm *mocks.TransactionManager) {
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
				mpsr.On("GetUnpaidPaymentSchedulesByLoanID", mock.Anything, uint(1), mock.Anything).Return([]domain.PaymentSchedule{
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
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mpr.On("CreatePayment", mock.Anything, mock.AnythingOfType(
					"*domain.Payment"), mock.Anything).Return(nil)
				mpsr.On("BulkPayPaymentSchedules", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mlr.On("UpdateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(nil).Return(nil).Return(nil).Return(nil)
			},
			expectedError: nil,
		},
		{
			name:       "Error Beginning Transaction",
			loanID:     1,
			paymentIDs: []uint{1},
			amount:     100.00,
			setupMocks: func(mpr *mocks.PaymentRepository, mpsr *mocks.PaymentScheduleRepository, mlr *mocks.LoanRepository, mtm *mocks.TransactionManager) {
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
				mpsr.On("GetUnpaidPaymentSchedulesByLoanID", mock.Anything, uint(1), mock.Anything).Return([]domain.PaymentSchedule{
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
				}, nil)

				mtm.On("Begin").Return(&gorm.DB{
					Error: errors.New("error beginning transaction"),
				})
			},
			expectedError: errors.New("error beginning transaction"),
		},
		{
			name:       "Error Creating Payment",
			loanID:     1,
			paymentIDs: []uint{1},
			amount:     100.00,
			setupMocks: func(mpr *mocks.PaymentRepository, mpsr *mocks.PaymentScheduleRepository, mlr *mocks.LoanRepository, mtm *mocks.TransactionManager) {
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
				mpsr.On("GetUnpaidPaymentSchedulesByLoanID", mock.Anything, uint(1), mock.Anything).Return([]domain.PaymentSchedule{
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
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mpr.On("CreatePayment", mock.Anything, mock.AnythingOfType("*domain.Payment"), mock.Anything).Return(errors.New("error creating payment"))
			},
			expectedError: errors.New("error creating payment"),
		},
		{
			name:       "Error Bulk Paying Payment Schedules",
			loanID:     1,
			paymentIDs: []uint{1},
			amount:     100.00,
			setupMocks: func(mpr *mocks.PaymentRepository, mpsr *mocks.PaymentScheduleRepository, mlr *mocks.LoanRepository, mtm *mocks.TransactionManager) {
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
				mpsr.On("GetUnpaidPaymentSchedulesByLoanID", mock.Anything, uint(1), mock.Anything).Return([]domain.PaymentSchedule{
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
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mpr.On("CreatePayment", mock.Anything, mock.AnythingOfType("*domain.Payment"), mock.Anything).Return(nil)
				mpsr.On("BulkPayPaymentSchedules", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error bulk paying payment schedules"))
			},
			expectedError: errors.New("error bulk paying payment schedules"),
		},
		{
			name:       "Error Updating Loan",
			loanID:     1,
			paymentIDs: []uint{1},
			amount:     100.00,
			setupMocks: func(mpr *mocks.PaymentRepository, mpsr *mocks.PaymentScheduleRepository, mlr *mocks.LoanRepository, mtm *mocks.TransactionManager) {
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
				mpsr.On("GetUnpaidPaymentSchedulesByLoanID", mock.Anything, uint(1), mock.Anything).Return([]domain.PaymentSchedule{
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
				}, nil)
				mtm.On("Begin").Return(&gorm.DB{})
				mtm.On("Commit", mock.Anything).Return(nil)
				mtm.On("Rollback", mock.Anything).Return(nil)
				mpr.On("CreatePayment", mock.Anything, mock.AnythingOfType("*domain.Payment"), mock.Anything).Return(nil)
				mpsr.On("BulkPayPaymentSchedules", mock.Anything, mock.Anything, mock.Anything).Return(nil)
				mlr.On("UpdateLoan", mock.Anything, mock.AnythingOfType("*domain.Loan"), mock.Anything).Return(errors.New("error updating loan"))
			},
			expectedError: errors.New("error updating loan"),
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			mockPaymentRepo := new(mocks.PaymentRepository)
			mockLoanRepo := new(mocks.LoanRepository)
			mockPaymentScheduleRepo := new(mocks.PaymentScheduleRepository)
			mockTransactionManager := new(mocks.TransactionManager)

			uc := paymentUsecase.NewPaymentUsecase(mockPaymentRepo, mockPaymentScheduleRepo, mockLoanRepo, mockTransactionManager, s.timeout)

			tt.setupMocks(mockPaymentRepo, mockPaymentScheduleRepo, mockLoanRepo, mockTransactionManager)
			err := uc.MakePayment(context.TODO(), tt.loanID, tt.paymentIDs, tt.amount)
			if tt.expectedError != nil {
				assert.Error(s.T(), err)
				assert.Equal(s.T(), tt.expectedError.Error(), err.Error())
			} else {
				assert.NoError(s.T(), err)
			}
		})
	}
}

func TestPaymentUsecaseSuite(t *testing.T) {
	suite.Run(t, new(PaymentUsecaseSuite))
}
