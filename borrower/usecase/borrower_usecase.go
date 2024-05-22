package usecase

import (
	"context"
	"time"

	"github.com/go-faker/faker/v4"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
)

type borrowerUsecase struct {
	borrowerRepo   domain.BorrowerRepository
	loanRepo       domain.LoanRepository
	contextTimeout time.Duration
}

func NewBorrowerUsecase(b domain.BorrowerRepository, l domain.LoanRepository, timeout time.Duration) domain.BorrowerUsecase {
	return &borrowerUsecase{
		borrowerRepo:   b,
		loanRepo:       l,
		contextTimeout: timeout,
	}
}

func (b *borrowerUsecase) IsDelinquent(ctx context.Context, borrowerID uint) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	_, err := b.borrowerRepo.FindBorrowerByID(ctx, borrowerID)
	if err != nil {
		return false, err
	}

	loans, err := b.loanRepo.GetLoansByBorrowerID(ctx, borrowerID)
	if err != nil {
		return false, err
	}

	today := time.Now()
	delinquentCount := 0

	for _, loan := range loans {
		for _, schedule := range loan.PaymentSchedules {
			if schedule.DueDate.Before(today) && !schedule.Paid {
				delinquentCount++
			}
		}
	}

	return delinquentCount >= 2, nil
}

func (b *borrowerUsecase) CreateBorrower(ctx context.Context) (*dto.CreateBorrowerResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, b.contextTimeout)
	defer cancel()

	var borrower domain.Borrower
	err := faker.FakeData(&borrower)
	if err != nil {
		return nil, err
	}

	err = b.borrowerRepo.CreateBorrower(ctx, &borrower, nil)
	if err != nil {
		return nil, err
	}

	return &dto.CreateBorrowerResponse{
		ID:        borrower.ID,
		FirstName: borrower.FirstName,
		LastName:  borrower.LastName,
		Email:     borrower.Email,
		CreatedAt: borrower.CreatedAt,
	}, nil
}
