package domain

import (
	"context"

	"github.com/greekrode/loan-engine-amartha/domain/dto"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	LoanID uint    `gorm:"not null" json:"loan_id"`
	Amount float64 `gorm:"not null" json:"amount"`
}

type PaymentUsecase interface {
	RequestPayment(ctx context.Context, loanID uint) (*dto.RequestPaymentResponse, error)
	MakePayment(ctx context.Context, loanID uint, paymentSchedulesID []uint, amount float64) error
}

type PaymentRepository interface {
	CreatePayment(ctx context.Context, payment *Payment, tx *gorm.DB) error
}
