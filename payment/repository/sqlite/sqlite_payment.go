package sqlite

import (
	"context"

	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"gorm.io/gorm"
)

type sqlitePaymentRepository struct {
	TransactionManager db.TransactionManager
}

func NewSQLitePaymentRepository(tm db.TransactionManager) *sqlitePaymentRepository {
	return &sqlitePaymentRepository{TransactionManager: tm}
}

func (s *sqlitePaymentRepository) CreatePayment(ctx context.Context, payment *domain.Payment, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}

	return tx.WithContext(ctx).Create(&payment).Error
}
