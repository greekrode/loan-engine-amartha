package http

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/greekrode/loan-engine-amartha/domain"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
)

type BorrowerHandler struct {
	BorrowerUsecase domain.BorrowerUsecase
}

func NewBorrowerHandler(g *gin.Engine, b domain.BorrowerUsecase) {
	handler := &BorrowerHandler{BorrowerUsecase: b}

	g.GET("/borrowers/:borrower_id/status", handler.CheckDelinquent)
	g.POST("/borrowers", handler.CreateBorrower)
}

func (b *BorrowerHandler) CheckDelinquent(c *gin.Context) {
	borrowerID := c.Param("borrower_id")
	parsedBorrowerID, err := strconv.ParseUint(borrowerID, 10, 32)
	if err != nil {
		c.JSON(400, dto.CommonResponse{Message: "invalid borrower ID format"})
		return
	}

	ctx := c.Request.Context()
	isDelinquent, err := b.BorrowerUsecase.IsDelinquent(ctx, uint(parsedBorrowerID))
	if err != nil {
		c.JSON(500, dto.CommonResponse{Message: err.Error()})
		return
	}

	c.JSON(200, dto.CheckDeliquentResponse{IsDelinquent: isDelinquent})
}

func (b *BorrowerHandler) CreateBorrower(c *gin.Context) {
	ctx := c.Request.Context()
	borrower, err := b.BorrowerUsecase.CreateBorrower(ctx)
	if err != nil {
		c.JSON(500, dto.CommonResponse{Message: err.Error()})
		return
	}

	c.JSON(201, borrower)
}
