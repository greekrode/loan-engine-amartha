package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
)

type LoanHandler struct {
	LoanUsecase domain.LoanUsecase
}

func NewLoanHandler(g *gin.Engine, l domain.LoanUsecase) {
	handler := &LoanHandler{LoanUsecase: l}
	g.POST("/loans", handler.CreateLoan)
	g.GET("/loans/:loan_id", handler.GetLoanDetails)
	g.GET("/loans/:loan_id/outstanding", handler.GetOutstanding)
}

func (l *LoanHandler) CreateLoan(c *gin.Context) {
	var req dto.CreateLoanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "invalid request body"})
		return
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "invalid date format, should be YYYY-MM-DD"})
		return
	}

	ctx := c.Request.Context()
	loanResponse, err := l.LoanUsecase.CreateLoan(ctx, req.BorrowerID, req.Principal, req.InterestRate, int32(req.Duration), startDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CommonResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, loanResponse)
}

func (l *LoanHandler) GetLoanDetails(c *gin.Context) {
	loanID := c.Param("loan_id")
	parsedLoanID, err := strconv.ParseUint(loanID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "invalid loan ID format"})
		return
	}

	ctx := c.Request.Context()
	loanResponse, err := l.LoanUsecase.GetLoanDetails(ctx, uint(parsedLoanID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CommonResponse{Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, loanResponse)
}

func (l *LoanHandler) GetOutstanding(c *gin.Context) {
	loanID := c.Param("loan_id")
	parsedLoanID, err := strconv.ParseUint(loanID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "invalid loan ID format"})
		return
	}

	ctx := c.Request.Context()
	outstanding, err := l.LoanUsecase.GetOutstandingAmount(ctx, uint(parsedLoanID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, dto.CommonResponse{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, dto.GetOutstandingResponse{OutstandingAmount: outstanding})
}
