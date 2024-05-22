package db

import (
	"gorm.io/gorm"
)

type TransactionManager interface {
	Begin() *gorm.DB
	Commit(tx *gorm.DB) error
	Rollback(tx *gorm.DB) error
	GetDB() *gorm.DB
}

type GormTransactionmanager struct {
	DB *gorm.DB
}

func NewGormTransactionManager(db *gorm.DB) *GormTransactionmanager {
	return &GormTransactionmanager{DB: db}
}

func (tm *GormTransactionmanager) Begin() *gorm.DB {
	return tm.DB.Begin()
}

func (tm *GormTransactionmanager) Commit(tx *gorm.DB) error {
	return tx.Commit().Error
}

func (tm *GormTransactionmanager) Rollback(tx *gorm.DB) error {
	return tx.Rollback().Error
}

func (tm *GormTransactionmanager) GetDB() *gorm.DB {
	return tm.DB
}
