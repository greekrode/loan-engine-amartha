package domain

import (
	"context"
	"time"

	"github.com/greekrode/loan-engine-amartha/domain/dto"
	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	BorrowerID        uint              `gorm:"not null" json:"borrower_id"`
	Principal         float64           `gorm:"not null" json:"principal"`
	InterestRate      float64           `gorm:"not null" json:"interest_rate"`
	DurationWeeks     int               `gorm:"not null" json:"duration_weeks"`
	OutstandingAmount float64           `gorm:"not null" json:"outstanding_amount"`
	StartDate         time.Time         `gorm:"not null" json:"start_date"`
	PaymentSchedules  []PaymentSchedule `gorm:"foreignKey:LoanID"`
}

type LoanUsecase interface {
	CreateLoan(ctx context.Context, borrowerID uint, principal, interestRate float64, durationWeeks int32, startDate time.Time) (*dto.CreateLoanResponse, error)
	GetLoanDetails(ctx context.Context, loanID uint) (*dto.GetLoanDetailsResponse, error)
	GetOutstandingAmount(ctx context.Context, loanID uint) (float64, error)
}

type LoanRepository interface {
	CreateLoan(ctx context.Context, loan *Loan, tx *gorm.DB) error

	FindLoanByID(ctx context.Context, loanID uint) (*Loan, error)
	GetLoansByBorrowerID(ctx context.Context, borrowerID uint) ([]Loan, error)

	UpdateLoan(ctx context.Context, loan *Loan, tx *gorm.DB) error
}
