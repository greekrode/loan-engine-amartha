package utils

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/greekrode/loan-engine-amartha/db"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func SetupMockDB(t *testing.T) (db.TransactionManager, sqlmock.Sqlmock, error) {
	sqlDB, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	mock.ExpectQuery("select sqlite_version()").WillReturnRows(sqlmock.NewRows([]string{"version"}).AddRow("3.31.1"))

	gormDB, err := gorm.Open(sqlite.Dialector{Conn: sqlDB}, &gorm.Config{})
	if err != nil {
		t.Fatalf("an error '%s' was not expected when setting up gorm with sqlite", err)
	}

	transactionManager := db.NewGormTransactionManager(gormDB)

	return transactionManager, mock, err
}
