package domain

import (
	"context"

	"github.com/greekrode/loan-engine-amartha/domain/dto"
	"gorm.io/gorm"
)

type Borrower struct {
	gorm.Model
	FirstName string `gorm:"not null" faker:"first_name"`
	LastName  string `gorm:"not null" faker:"last_name"`
	Email     string `gorm:"not null" faker:"email"`
}

type BorrowerUsecase interface {
	IsDelinquent(ctx context.Context, borrowerID uint) (bool, error)
	CreateBorrower(ctx context.Context) (*dto.CreateBorrowerResponse, error)
}

type BorrowerRepository interface {
	CreateBorrower(ctx context.Context, borrower *Borrower, tx *gorm.DB) error
	FindBorrowerByID(ctx context.Context, borrowerID uint) (*Borrower, error)
}
