package usecase

import (
	"context"
	"math"
	"time"

	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
)

type loanUsecase struct {
	borrowerRepo        domain.BorrowerRepository
	paymentScheduleRepo domain.PaymentScheduleRepository
	loanRepo            domain.LoanRepository
	transactionManager  db.TransactionManager
	contextTimeout      time.Duration
}

func NewLoanUsecase(b domain.BorrowerRepository, p domain.PaymentScheduleRepository, l domain.LoanRepository, tm db.TransactionManager, timeout time.Duration) domain.LoanUsecase {
	return &loanUsecase{
		borrowerRepo:        b,
		paymentScheduleRepo: p,
		loanRepo:            l,
		transactionManager:  tm,
		contextTimeout:      timeout,
	}
}

func (l *loanUsecase) CreateLoan(ctx context.Context, borrowerID uint, principal, interestRate float64, durationWeeks int32, startDate time.Time) (*dto.CreateLoanResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	_, err := l.borrowerRepo.FindBorrowerByID(ctx, borrowerID)
	if err != nil {
		return nil, err
	}

	tx := l.transactionManager.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			l.transactionManager.Rollback(tx)
			panic(r)
		}
	}()

	weeklyInterestRate := (interestRate / 100) / 52
	var totalOutstandingAmount float64 = 0
	var paymentSchedules []domain.PaymentSchedule

	loan := domain.Loan{
		BorrowerID:    borrowerID,
		Principal:     principal,
		InterestRate:  interestRate,
		DurationWeeks: int(durationWeeks),
		StartDate:     startDate,
	}

	err = l.loanRepo.CreateLoan(ctx, &loan, tx)
	if err != nil {
		return nil, err
	}

	for i := 0; i < int(durationWeeks); i++ {
		weeklyInterest := principal * weeklyInterestRate
		weeklyPayment := principal/float64(durationWeeks) + weeklyInterest
		weeklyPayment = math.Round(weeklyPayment*100) / 100
		totalOutstandingAmount += weeklyPayment

		paymentSchedules = append(paymentSchedules, domain.PaymentSchedule{
			LoanID:    loan.ID,
			DueDate:   startDate.Add(time.Duration(i+1) * 7 * 24 * time.Hour),
			DueAmount: weeklyPayment,
		})
	}

	loan.OutstandingAmount = totalOutstandingAmount
	err = l.loanRepo.UpdateLoan(ctx, &loan, tx)
	if err != nil {
		l.transactionManager.Rollback(tx)
		return nil, err
	}

	err = l.paymentScheduleRepo.BulkCreatePaymentSchedule(ctx, paymentSchedules, tx)
	if err != nil {
		l.transactionManager.Rollback(tx)
		return nil, err
	}

	err = l.transactionManager.Commit(tx)
	if err != nil {
		l.transactionManager.Rollback(tx)
		return nil, err
	}

	return assembleCreateLoanResponse(&loan, paymentSchedules), nil
}

func assembleCreateLoanResponse(loan *domain.Loan, paymentSchedules []domain.PaymentSchedule) *dto.CreateLoanResponse {
	paymentScheduleResponses := make([]dto.GetPaymentScheduleResponse, len(paymentSchedules))
	for i, ps := range paymentSchedules {
		paymentScheduleResponses[i] = dto.GetPaymentScheduleResponse{
			DueAmount: ps.DueAmount,
			DueDate:   ps.DueDate,
			Paid:      ps.Paid,
		}
	}

	loanResponse := dto.CreateLoanResponse{
		ID:                loan.ID,
		Principal:         loan.Principal,
		InterestRate:      loan.InterestRate,
		Duration:          loan.DurationWeeks,
		StartDate:         loan.StartDate,
		OutstandingAmount: loan.OutstandingAmount,
		PaymentSchedules:  paymentScheduleResponses,
	}

	return &loanResponse
}

func (l *loanUsecase) GetLoanDetails(ctx context.Context, loanID uint) (*dto.GetLoanDetailsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	loan, err := l.loanRepo.FindLoanByID(ctx, loanID)
	if err != nil {
		return nil, err
	}

	borrower, err := l.borrowerRepo.FindBorrowerByID(ctx, loan.BorrowerID)
	if err != nil {
		return nil, err
	}

	return assembleLoanDetailsResponse(loan, borrower), nil
}

func assembleLoanDetailsResponse(loan *domain.Loan, borrower *domain.Borrower) *dto.GetLoanDetailsResponse {
	paymentScheduleResponses := make([]dto.GetPaymentScheduleResponse, len(loan.PaymentSchedules))
	for i, ps := range loan.PaymentSchedules {
		paymentScheduleResponses[i] = dto.GetPaymentScheduleResponse{
			DueAmount: ps.DueAmount,
			DueDate:   ps.DueDate,
			Paid:      ps.Paid,
		}
	}

	borrowerResponse := dto.GetBorrowerResponse{
		ID:        borrower.ID,
		FirstName: borrower.FirstName,
		LastName:  borrower.LastName,
		Email:     borrower.Email,
		CreatedAt: borrower.CreatedAt,
	}

	loanResponse := dto.GetLoanDetailsResponse{
		Principal:         loan.Principal,
		InterestRate:      loan.InterestRate,
		Duration:          loan.DurationWeeks,
		OutstandingAmount: loan.OutstandingAmount,
		StartDate:         loan.StartDate,
		CreatedAt:         loan.CreatedAt,
		Borrower:          borrowerResponse,
		PaymentSchedule:   paymentScheduleResponses,
	}

	return &loanResponse
}

func (l *loanUsecase) GetOutstandingAmount(ctx context.Context, loanID uint) (float64, error) {
	ctx, cancel := context.WithTimeout(ctx, l.contextTimeout)
	defer cancel()

	loan, err := l.loanRepo.FindLoanByID(ctx, loanID)
	if err != nil {
		return 0, err
	}

	return loan.OutstandingAmount, nil
}
