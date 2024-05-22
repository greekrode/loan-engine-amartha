package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"gorm.io/gorm"
)

type sqliteLoanRepository struct {
	TransactionManager db.TransactionManager
}

func NewSQLiteLoanRepository(tm db.TransactionManager) *sqliteLoanRepository {
	return &sqliteLoanRepository{TransactionManager: tm}
}

func (s *sqliteLoanRepository) CreateLoan(ctx context.Context, loan *domain.Loan, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}

	return tx.WithContext(ctx).Create(&loan).Error
}

func (s *sqliteLoanRepository) FindLoanByID(ctx context.Context, loanID uint) (*domain.Loan, error) {
	var loan domain.Loan

	err := s.TransactionManager.GetDB().WithContext(ctx).Preload("PaymentSchedules").First(&loan, loanID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Loan not found")
		}
		return nil, err
	}

	return &loan, nil
}

func (s *sqliteLoanRepository) GetLoansByBorrowerID(ctx context.Context, borrowerID uint) ([]domain.Loan, error) {
	var loans []domain.Loan
	err := s.TransactionManager.GetDB().WithContext(ctx).Where("borrower_id = ?", borrowerID).Preload("PaymentSchedules").Find(&loans).Error
	if err != nil {
		return nil, err
	}

	if len(loans) == 0 {
		return nil, fmt.Errorf("Loan not found")
	}

	return loans, nil
}

func (s *sqliteLoanRepository) UpdateLoan(ctx context.Context, loan *domain.Loan, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}
	return tx.WithContext(ctx).Save(loan).Error
}
