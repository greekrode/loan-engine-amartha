package sqlite_test

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/greekrode/loan-engine-amartha/borrower/repository/sqlite"
	"github.com/greekrode/loan-engine-amartha/db"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/utils"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type BorrowerRepositorySuite struct {
	suite.Suite
	tm   db.TransactionManager
	mock sqlmock.Sqlmock
}

func (s *BorrowerRepositorySuite) SetupSuite() {
	var err error
	s.tm, s.mock, err = utils.SetupMockDB(s.T())
	s.Require().NoError(err)
}

func (s *BorrowerRepositorySuite) AfterTest(_, _ string) {
	s.Require().NoError(s.mock.ExpectationsWereMet())
}

func (s *BorrowerRepositorySuite) TestCreateBorrower() {
	tests := []struct {
		name     string
		setup    func()
		borrower domain.Borrower
		wantErr  bool
	}{
		{
			name: "Success",
			setup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec("INSERT INTO `borrowers`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "John", "Doe", "john.doe@example.com").WillReturnResult(sqlmock.NewResult(1, 1))
				s.mock.ExpectCommit()
			},
			borrower: domain.Borrower{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
			wantErr:  false,
		},
		{
			name: "Failure",
			setup: func() {
				s.mock.ExpectBegin()
				s.mock.ExpectExec("INSERT INTO `borrowers`").WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), "John", "Doe", "john.doe@example.com").WillReturnError(fmt.Errorf("insert error"))
				s.mock.ExpectRollback()
			},
			borrower: domain.Borrower{FirstName: "John", LastName: "Doe", Email: "john.doe@example.com"},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()
			repo := sqlite.NewSQLiteBorrowerRepository(s.tm)
			err := repo.CreateBorrower(context.TODO(), &tt.borrower, nil)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
			}
		})
	}
}

func (s *BorrowerRepositorySuite) TestFindBorrowerByID() {
	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)
	tests := []struct {
		name       string
		borrowerID uint
		setup      func()
		borrower   *domain.Borrower
		wantErr    bool
	}{
		{
			name:       "Success",
			borrowerID: 1,
			setup: func() {
				query := "SELECT * FROM `borrowers` WHERE `borrowers`.`id` = ? AND `borrowers`.`deleted_at` IS NULL ORDER BY `borrowers`.`id` LIMIT 1"
				escapedQuery := regexp.QuoteMeta(query)
				rows := sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "first_name", "last_name", "email"}).
					AddRow(1, fixedTime, fixedTime, nil, "John", "Doe", "john.doe@example.com")
				s.mock.ExpectQuery(escapedQuery).WithArgs(1).WillReturnRows(rows)
			},
			borrower: &domain.Borrower{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: fixedTime,
					UpdatedAt: fixedTime,
					DeletedAt: gorm.DeletedAt{Valid: false},
				},
				FirstName: "John",
				LastName:  "Doe",
				Email:     "john.doe@example.com",
			},
			wantErr: false,
		},
		{
			name:       "NotFound",
			borrowerID: 99,
			setup: func() {
				query := "SELECT * FROM `borrowers` WHERE `borrowers`.`id` = ? AND `borrowers`.`deleted_at` IS NULL ORDER BY `borrowers`.`id` LIMIT 1"
				escapedQuery := regexp.QuoteMeta(query)
				s.mock.ExpectQuery(escapedQuery).WithArgs(99).WillReturnRows(sqlmock.NewRows(nil))
			},
			borrower: nil,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		s.Run(tt.name, func() {
			tt.setup()
			repo := sqlite.NewSQLiteBorrowerRepository(s.tm)
			got, err := repo.FindBorrowerByID(context.TODO(), tt.borrowerID)
			if tt.wantErr {
				s.Error(err)
			} else {
				s.NoError(err)
				s.Equal(tt.borrower.FirstName, got.FirstName)
				s.Equal(tt.borrower.LastName, got.LastName)
				s.Equal(tt.borrower.Email, got.Email)
				s.Equal(tt.borrower.ID, got.ID)
			}
		})
	}
}

func TestBorrowerRepositorySuite(t *testing.T) {
	suite.Run(t, new(BorrowerRepositorySuite))
}
