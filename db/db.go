package db

import (
	"log"

	"github.com/greekrode/loan-engine-amartha/domain"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var TrxManager TransactionManager

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("data.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	DB.AutoMigrate(&domain.Borrower{}, &domain.Loan{}, &domain.PaymentSchedule{}, &domain.Payment{})

	TrxManager = NewGormTransactionManager(DB)
}
