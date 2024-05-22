package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
)

type paymentUsecase struct {
	paymentRepo         domain.PaymentRepository
	paymentScheduleRepo domain.PaymentScheduleRepository
	loanRepo            domain.LoanRepository
	transactionManager  db.TransactionManager
	contextTimeout      time.Duration
}

func NewPaymentUsecase(p domain.PaymentRepository, ps domain.PaymentScheduleRepository, l domain.LoanRepository, tm db.TransactionManager, timeout time.Duration) domain.PaymentUsecase {
	return &paymentUsecase{
		paymentRepo:         p,
		paymentScheduleRepo: ps,
		loanRepo:            l,
		transactionManager:  tm,
		contextTimeout:      timeout,
	}
}

func (p *paymentUsecase) retrieveAndValidateLoanAndSchedules(ctx context.Context, loanID uint) (*domain.Loan, []domain.PaymentSchedule, error) {
	loan, err := p.loanRepo.FindLoanByID(ctx, loanID)
	if err != nil {
		return nil, nil, err
	}

	today := time.Now()
	paymentSchedules, err := p.paymentScheduleRepo.GetUnpaidPaymentSchedulesByLoanID(ctx, loanID, today)
	if err != nil {
		return nil, nil, err
	}

	return loan, paymentSchedules, nil
}

func (p *paymentUsecase) RequestPayment(ctx context.Context, loanID uint) (*dto.RequestPaymentResponse, error) {
	_, paymentSchedules, err := p.retrieveAndValidateLoanAndSchedules(ctx, loanID)
	if err != nil {
		return nil, err
	}

	totalDue := 0.0
	scheduleResponses := make([]dto.GetPaymentScheduleResponse, len(paymentSchedules))
	for i, schedule := range paymentSchedules {
		totalDue += schedule.DueAmount
		scheduleResponses[i] = dto.GetPaymentScheduleResponse{
			ID:        schedule.ID,
			DueDate:   schedule.DueDate,
			DueAmount: schedule.DueAmount,
		}
	}

	return &dto.RequestPaymentResponse{
		TotalDue:         totalDue,
		PaymentSchedules: scheduleResponses,
	}, nil
}

func (p *paymentUsecase) MakePayment(ctx context.Context, loanID uint, paymentSchedulesID []uint, amount float64) error {
	loan, paymentSchedules, err := p.retrieveAndValidateLoanAndSchedules(ctx, loanID)
	if err != nil {
		return err
	}

	totalDue := 0.0
	for _, schedule := range paymentSchedules {
		totalDue += schedule.DueAmount
	}

	if amount != totalDue {
		return errors.New("payment amount does not match the total due amount")
	}

	tx := p.transactionManager.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			p.transactionManager.Rollback(tx)
			panic(r)
		}
	}()

	payment := &domain.Payment{
		LoanID: loanID,
		Amount: amount,
	}
	if err := p.paymentRepo.CreatePayment(ctx, payment, tx); err != nil {
		p.transactionManager.Rollback(tx)
		return err
	}

	if err := p.paymentScheduleRepo.BulkPayPaymentSchedules(ctx, paymentSchedulesID, tx); err != nil {
		p.transactionManager.Rollback(tx)
		return err
	}

	loan.OutstandingAmount -= amount
	if err := p.loanRepo.UpdateLoan(ctx, loan, tx); err != nil {
		p.transactionManager.Rollback(tx)
		return err
	}

	return p.transactionManager.Commit(tx)
}
