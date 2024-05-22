package sqlite_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/loan/repository/sqlite"
	"github.com/greekrode/loan-engine-amartha/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type LoanRepositorySuite struct {
	suite.Suite
	tm   db.TransactionManager
	mock sqlmock.Sqlmock
}

func (s *LoanRepositorySuite) SetupSuite() {
	var err error
	s.tm, s.mock, err = utils.SetupMockDB(s.T())
	s.Require().NoError(err)
}

func (s *LoanRepositorySuite) AfterTest(_, _ string) {
	s.Require().NoError(s.mock.ExpectationsWereMet())
}

func (s *LoanRepositorySuite) TestCreateLoan() {
	tests := []struct {
		name    string
		setup   func()
		loan    domain.Loan
		wantErr bool
	}{
		{
			name: "Success",
			setup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec("INSERT INTO `loans`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 100.00, 10.00, 52, 1000.00, sqlmock.AnyArg()).WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			},
			loan: domain.Loan{
				BorrowerID:        1,
				Principal:         100,
				InterestRate:      10,
				DurationWeeks:     52,
				OutstandingAmount: 1000,
			},
			wantErr: false,
		},
		{
			name: "Failure",
			setup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec("INSERT INTO `loans`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), 1, 100.00, 10.00, 52, 1000.00, sqlmock.AnyArg()).WillReturnError(fmt.Errorf("insert error"))
				s.mock.ExpectRollback()
			},
			loan: domain.Loan{
				BorrowerID:        1,
				Principal:         100,
				InterestRate:      10,
				DurationWeeks:     52,
				OutstandingAmount: 1000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()
			repo := sqlite.NewSQLiteLoanRepository(s.tm)
			err := repo.CreateLoan(context.TODO(), &tt.loan, nil)
			if tt.wantErr {
				s.Require().Error(err)
			} else {
				s.Require().NoError(err)
			}
		})
	}
}

func (s *LoanRepositorySuite) TestFindLoanByID() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name    string
		loanID  uint
		setup   func()
		loan    *domain.Loan
		wantErr bool
	}{
		{
			name:   "Success",
			loanID: 1,
			setup: func() {
				loanQuery := "SELECT * FROM `loans` WHERE `loans`.`id` = ? AND `loans`.`deleted_at` IS NULL ORDER BY `loans`.`id` LIMIT 1"
				loanEscapedQuery := regexp.QuoteMeta(loanQuery)
				loanRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "borrower_id", "principal", "interest_rate", "duration_weeks", "outstanding_amount", "start_date"}).
					AddRow(1, fixedTime, fixedTime, nil, 1, 100.00, 10.00, 52, 1000.00, fixedTime)
				s.mock.ExpectQuery(loanEscapedQuery).WithArgs(1).WillReturnRows(loanRows)

				paymentSchedulesQuery := "SELECT * FROM `payment_schedules` WHERE `payment_schedules`.`loan_id` = ? AND `payment_schedules`.`deleted_at` IS NULL"
				paymentSchedulesEscapedQuery := regexp.QuoteMeta(paymentSchedulesQuery)
				paymentSchedulesRows := sqlmock.NewRows([]string{"id", "loan_id", "due_date", "amount_due", "status"}).
					AddRow(1, 1, fixedTime, 500.00, "pending")
				s.mock.ExpectQuery(paymentSchedulesEscapedQuery).WithArgs(1).WillReturnRows(paymentSchedulesRows)
			},
			loan: &domain.Loan{
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
						LoanID:    1,
						DueDate:   fixedTime,
						DueAmount: 500.00,
						Paid:      false,
					},
				},
			},
			wantErr: false,
		},
		{
			name:   "NotFound",
			loanID: 1,
			setup: func() {
				loanQuery := "SELECT * FROM `loans` WHERE `loans`.`id` = ? AND `loans`.`deleted_at` IS NULL ORDER BY `loans`.`id` LIMIT 1"
				loanEscapedQuery := regexp.QuoteMeta(loanQuery)
				loanRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "borrower_id", "principal", "interest_rate", "duration_weeks", "outstanding_amount", "start_date"})
				s.mock.ExpectQuery(loanEscapedQuery).WithArgs(1).WillReturnRows(loanRows)
			},
			loan:    nil,
			wantErr: true,
		},
		{
			name:   "DatabaseError",
			loanID: 1,
			setup: func() {
				loanQuery := "SELECT * FROM `loans` WHERE `loans`.`id` = ? AND `loans`.`deleted_at` IS NULL ORDER BY `loans`.`id` LIMIT 1"
				loanEscapedQuery := regexp.QuoteMeta(loanQuery)
				s.mock.ExpectQuery(loanEscapedQuery).WithArgs(1).WillReturnError(fmt.Errorf("database error"))
			},
			loan:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()
			repo := sqlite.NewSQLiteLoanRepository(s.tm)
			got, err := repo.FindLoanByID(context.TODO(), tt.loanID)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.loan.BorrowerID, got.BorrowerID)
				s.Equal(tt.loan.Principal, got.Principal)
				s.Equal(tt.loan.InterestRate, got.InterestRate)
				s.Equal(tt.loan.DurationWeeks, got.DurationWeeks)
				s.Equal(tt.loan.OutstandingAmount, got.OutstandingAmount)
				s.Equal(tt.loan.StartDate, got.StartDate)
				s.Len(got.PaymentSchedules, 1)
			}
		})
	}
}

func (s *LoanRepositorySuite) TestGetLoansByBorrowerID() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		borrowerID uint
		setup      func()
		loan       []domain.Loan
		wantErr    bool
	}{
		{
			name:       "Success",
			borrowerID: 1,
			setup: func() {
				loanQuery := "SELECT * FROM `loans` WHERE borrower_id = ? AND `loans`.`deleted_at` IS NULL"
				loanEscapedQuery := regexp.QuoteMeta(loanQuery)
				loanRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "borrower_id", "principal", "interest_rate", "duration_weeks", "outstanding_amount", "start_date"}).
					AddRow(1, fixedTime, fixedTime, nil, 1, 100.00, 10.00, 52, 1000.00, fixedTime)
				s.mock.ExpectQuery(loanEscapedQuery).WithArgs(1).WillReturnRows(loanRows)

				paymentSchedulesQuery := "SELECT * FROM `payment_schedules` WHERE `payment_schedules`.`loan_id` = ? AND `payment_schedules`.`deleted_at` IS NULL"
				paymentSchedulesEscapedQuery := regexp.QuoteMeta(paymentSchedulesQuery)
				paymentSchedulesRows := sqlmock.NewRows([]string{"id", "loan_id", "due_date", "amount_due", "status"}).
					AddRow(1, 1, fixedTime, 500.00, "pending")
				s.mock.ExpectQuery(paymentSchedulesEscapedQuery).WithArgs(1).WillReturnRows(paymentSchedulesRows)
			},
			loan: []domain.Loan{
				{
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
							LoanID:    1,
							DueDate:   fixedTime,
							DueAmount: 500.00,
							Paid:      false,
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name:       "NotFound",
			borrowerID: 1,
			setup: func() {
				loanQuery := "SELECT * FROM `loans` WHERE borrower_id = ? AND `loans`.`deleted_at` IS NULL"
				loanEscapedQuery := regexp.QuoteMeta(loanQuery)
				loanRows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "borrower_id", "principal", "interest_rate", "duration_weeks", "outstanding_amount", "start_date"})
				s.mock.ExpectQuery(loanEscapedQuery).WithArgs(1).WillReturnRows(loanRows)
			},
			loan:    nil,
			wantErr: true,
		},
		{
			name:       "DatabaseError",
			borrowerID: 1,
			setup: func() {
				loanQuery := "SELECT * FROM `loans` WHERE borrower_id = ? AND `loans`.`deleted_at` IS NULL"
				loanEscapedQuery := regexp.QuoteMeta(loanQuery)
				s.mock.ExpectQuery(loanEscapedQuery).WithArgs(1).WillReturnError(fmt.Errorf("database error"))
			},
			loan:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()
			repo := sqlite.NewSQLiteLoanRepository(s.tm)
			got, err := repo.GetLoansByBorrowerID(context.TODO(), tt.borrowerID)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.loan[0].BorrowerID, got[0].BorrowerID)
				s.Equal(tt.loan[0].Principal, got[0].Principal)
				s.Equal(tt.loan[0].InterestRate, got[0].InterestRate)
				s.Equal(tt.loan[0].DurationWeeks, got[0].DurationWeeks)
				s.Equal(tt.loan[0].OutstandingAmount, got[0].OutstandingAmount)
				s.Equal(tt.loan[0].StartDate, got[0].StartDate)
				s.Len(got[0].PaymentSchedules, 1)
			}
		})
	}
}

func TestLoanRepositorySuite(t *testing.T) {
	suite.Run(t, new(LoanRepositorySuite))
}
