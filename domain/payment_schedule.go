package domain

import (
	"context"
	"time"

	"gorm.io/gorm"
)

type PaymentSchedule struct {
	gorm.Model
	DueAmount float64   `gorm:"not null" json:"due_amount"`
	DueDate   time.Time `gorm:"not null" json:"due_date"`
	Paid      bool      `gorm:"not null;default:false" json:"paid"`
	LoanID    uint      `gorm:"not null" json:"loan_id"`
}

type PaymentScheduleUsecase interface {
	MakePayment(ctx context.Context, loanID uint, amount float64) error
}

type PaymentScheduleRepository interface {
	CreatePaymentSchedule(ctx context.Context, bs *PaymentSchedule, tx *gorm.DB) error
	BulkCreatePaymentSchedule(ctx context.Context, bs []PaymentSchedule, tx *gorm.DB) error

	GetPaymentSchedulesByLoanID(ctx context.Context, loanID uint) ([]PaymentSchedule, error)
	GetUnpaidPaymentSchedulesByLoanID(ctx context.Context, loanID uint, date time.Time) ([]PaymentSchedule, error)

	UpdatePaymentSchedule(ctx context.Context, bs *PaymentSchedule, tx *gorm.DB) error
	BulkPayPaymentSchedules(ctx context.Context, paymentSchedulesID []uint, tx *gorm.DB) error
}
