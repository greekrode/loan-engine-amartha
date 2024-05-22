package main

import (
	"log"
	"time"

	"github.com/gin-gonic/gin"
	_borrowerHttpDelivery "github.com/greekrode/loan-engine-amartha/borrower/delivery/http"
	_borrowerRepo "github.com/greekrode/loan-engine-amartha/borrower/repository/sqlite"
	_borrowerUseCase "github.com/greekrode/loan-engine-amartha/borrower/usecase"
	"github.com/greekrode/loan-engine-amartha/db"
	_loanHttpDelivery "github.com/greekrode/loan-engine-amartha/loan/delivery/http"
	_loanRepo "github.com/greekrode/loan-engine-amartha/loan/repository/sqlite"
	_loanUsecase "github.com/greekrode/loan-engine-amartha/loan/usecase"
	_paymentHttpDelivery "github.com/greekrode/loan-engine-amartha/payment/delivery/http"
	_paymentRepo "github.com/greekrode/loan-engine-amartha/payment/repository/sqlite"
	_paymentUsecase "github.com/greekrode/loan-engine-amartha/payment/usecase"
	_paymentScheduleRepo "github.com/greekrode/loan-engine-amartha/payment_schedule/repository/sqlite"
)

func main() {
	db.InitDB()

	gin.SetMode(gin.ReleaseMode)

	router := gin.Default()
	router.SetTrustedProxies(nil)

	timeoutCtx := time.Duration(30) * time.Second

	loanRepo := _loanRepo.NewSQLiteLoanRepository(db.TrxManager)
	borrowerRepo := _borrowerRepo.NewSQLiteBorrowerRepository(db.TrxManager)
	paymentScheduleRepo := _paymentScheduleRepo.NewSQLitePaymentScheduleRepository(db.TrxManager)
	paymentRepo := _paymentRepo.NewSQLitePaymentRepository(db.TrxManager)

	loanUsecase := _loanUsecase.NewLoanUsecase(borrowerRepo, paymentScheduleRepo, loanRepo, db.TrxManager, timeoutCtx)
	borrowerUseCase := _borrowerUseCase.NewBorrowerUsecase(borrowerRepo, loanRepo, timeoutCtx)
	paymentUsecase := _paymentUsecase.NewPaymentUsecase(paymentRepo, paymentScheduleRepo, loanRepo, db.TrxManager, timeoutCtx)

	_loanHttpDelivery.NewLoanHandler(router, loanUsecase)
	_borrowerHttpDelivery.NewBorrowerHandler(router, borrowerUseCase)
	_paymentHttpDelivery.NewPaymentHandler(router, paymentUsecase)

	log.Fatal(router.Run(":8080"))
}
