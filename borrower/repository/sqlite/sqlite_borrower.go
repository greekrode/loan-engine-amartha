package sqlite

import (
	"context"
	"errors"
	"fmt"

	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"gorm.io/gorm"
)

type sqliteBorrowerRepository struct {
	TransactionManager db.TransactionManager
}

func NewSQLiteBorrowerRepository(tm db.TransactionManager) *sqliteBorrowerRepository {
	return &sqliteBorrowerRepository{TransactionManager: tm}
}

func (s *sqliteBorrowerRepository) CreateBorrower(ctx context.Context, borrower *domain.Borrower, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}

	return tx.WithContext(ctx).Create(&borrower).Error
}

func (s *sqliteBorrowerRepository) FindBorrowerByID(ctx context.Context, borrowerID uint) (*domain.Borrower, error) {
	var borrower domain.Borrower

	err := s.TransactionManager.GetDB().WithContext(ctx).First(&borrower, borrowerID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Borrower not found")
		}
		return nil, err
	}

	return &borrower, nil
}
