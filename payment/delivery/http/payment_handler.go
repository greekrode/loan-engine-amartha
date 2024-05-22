package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
)

type PaymentHandler struct {
	PaymentUsecase domain.PaymentUsecase
}

func NewPaymentHandler(g *gin.Engine, p domain.PaymentUsecase) {
	handler := &PaymentHandler{PaymentUsecase: p}

	g.GET("/payments/:loan_id", handler.RequestPayment)
	g.POST("/payments/:loan_id", handler.MakePayment)
}

func (p *PaymentHandler) RequestPayment(c *gin.Context) {
	loanID := c.Param("loan_id")
	parsedLoanID, err := strconv.ParseUint(loanID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "invalid loan ID format"})
		return
	}

	ctx := c.Request.Context()
	paymentResponse, err := p.PaymentUsecase.RequestPayment(ctx, uint(parsedLoanID))
	if err != nil {
		c.JSON(500, dto.CommonResponse{Message: err.Error()})
		return
	}

	c.JSON(200, paymentResponse)
}

func (p *PaymentHandler) MakePayment(c *gin.Context) {
	loanID := c.Param("loan_id")
	parsedLoanID, err := strconv.ParseUint(loanID, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "invalid loan ID format"})
		return
	}

	ctx := c.Request.Context()
	var paymentRequest dto.MakePaymentRequest
	if err := c.ShouldBindJSON(&paymentRequest); err != nil {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "invalid request body"})
		return
	}

	if len(paymentRequest.PaymentSchedulesID) == 0 {
		c.JSON(http.StatusBadRequest, dto.CommonResponse{Message: "payment schedule ID is required"})
		return
	}

	err = p.PaymentUsecase.MakePayment(ctx, uint(parsedLoanID), paymentRequest.PaymentSchedulesID, paymentRequest.Amount)
	if err != nil {
		c.JSON(500, dto.CommonResponse{Message: err.Error()})
		return
	}

	c.JSON(200, dto.CommonResponse{Message: "payment success"})

}
