package sqlite

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"gorm.io/gorm"
)

type sqlitePaymentScheduleRepository struct {
	TransactionManager db.TransactionManager
}

func NewSQLitePaymentScheduleRepository(tm db.TransactionManager) *sqlitePaymentScheduleRepository {
	return &sqlitePaymentScheduleRepository{TransactionManager: tm}
}

func (s *sqlitePaymentScheduleRepository) CreatePaymentSchedule(ctx context.Context, paymentSchedule *domain.PaymentSchedule, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}

	return tx.WithContext(ctx).Create(&paymentSchedule).Error
}

func (s *sqlitePaymentScheduleRepository) BulkCreatePaymentSchedule(ctx context.Context, paymentSchedules []domain.PaymentSchedule, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}

	return tx.WithContext(ctx).Create(&paymentSchedules).Error
}

func (s *sqlitePaymentScheduleRepository) GetPaymentSchedulesByLoanID(ctx context.Context, loanID uint) ([]domain.PaymentSchedule, error) {
	var paymentSchedules []domain.PaymentSchedule

	err := s.TransactionManager.GetDB().WithContext(ctx).Where(&domain.PaymentSchedule{LoanID: loanID}).Find(&paymentSchedules).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Payment schedule not found")
		}
		return nil, err
	}

	return paymentSchedules, nil
}

func (s *sqlitePaymentScheduleRepository) GetUnpaidPaymentSchedulesByLoanID(ctx context.Context, loanID uint, date time.Time) ([]domain.PaymentSchedule, error) {
	var paymentSchedules []domain.PaymentSchedule

	err := s.TransactionManager.GetDB().WithContext(ctx).Where("loan_id = ? AND due_date <= ? AND paid = ?", loanID, date, false).Find(&paymentSchedules).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Payment schedule not found")
		}
		return nil, err
	}

	return paymentSchedules, nil
}

func (s *sqlitePaymentScheduleRepository) UpdatePaymentSchedule(ctx context.Context, paymentSchedule *domain.PaymentSchedule, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}

	return tx.WithContext(ctx).Save(&paymentSchedule).Error
}

func (s *sqlitePaymentScheduleRepository) BulkPayPaymentSchedules(ctx context.Context, paymentSchedulesID []uint, tx *gorm.DB) error {
	if tx == nil {
		tx = s.TransactionManager.GetDB()
	}

	result := tx.WithContext(ctx).Model(domain.PaymentSchedule{}).Where("id IN ?", paymentSchedulesID).Update("paid", true)
	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("no rows were updated")
	}

	return nil
}
