package http_test

import (
	"bytes"
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
	"github.com/greekrode/loan-engine-amartha/domain/mocks"
	paymentHttp "github.com/greekrode/loan-engine-amartha/payment/delivery/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRouter(mockUCase *mocks.PaymentUsecase) *gin.Engine {
	router := gin.Default()
	router.GET("/payments/:loan_id", func(c *gin.Context) {
		handler := paymentHttp.PaymentHandler{
			PaymentUsecase: mockUCase,
		}
		handler.RequestPayment(c)
	})
	router.POST("/payments/:loan_id", func(c *gin.Context) {
		handler := paymentHttp.PaymentHandler{
			PaymentUsecase: mockUCase,
		}
		handler.MakePayment(c)
	})
	return router
}

func TestRequestPayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		loanID         string
		mockUsecase    *mocks.PaymentUsecase
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Valid Request Payment",
			loanID: "1",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("RequestPayment", mock.Anything, uint(1)).Return(&dto.RequestPaymentResponse{
					TotalDue: 1000.00,
					PaymentSchedules: []dto.GetPaymentScheduleResponse{
						{
							ID:        1,
							DueAmount: 1000.00,
							DueDate:   time.Time{},
							Paid:      false,
						},
					},
				}, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"total_due": 1000,
				"payment_schedules": [
					{
						"id": 1,
						"due_amount": 1000,
						"due_date": "0001-01-01T00:00:00Z",
						"paid": false
					}
				]
			}`,
		},
		{
			name:   "Invalid Loan ID",
			loanID: "invalid",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("RequestPayment", mock.Anything, uint(1)).Return(&dto.RequestPaymentResponse{}, errors.New("invalid loan id"))
				return mockUsecase
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid loan ID format"}`,
		},
		{
			name:   "Internal Server Error",
			loanID: "1",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("RequestPayment", mock.Anything, uint(1)).Return(&dto.RequestPaymentResponse{}, errors.New("internal server error"))
				return mockUsecase
			}(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal server error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.mockUsecase)
			req, err := http.NewRequestWithContext(context.TODO(), "GET", "/payments/"+tt.loanID, bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestMakePayment(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		loanID         string
		mockUsecase    *mocks.PaymentUsecase
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Valid Make Payment",
			loanID: "1",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("MakePayment", mock.Anything, uint(1), []uint{1}, 1000.00).Return(nil)
				return mockUsecase
			}(),
			requestBody: `{
				"amount": 1000, 
				"payment_schedules_id": [
					1
				]
			}`,
			expectedStatus: http.StatusOK,
			expectedBody:   `{"message":"payment success"}`,
		},
		{
			name:   "Invalid Loan ID",
			loanID: "invalid",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("MakePayment", mock.Anything, uint(1), []uint{1}, 1000.00).Return(errors.New("invalid loan id"))
				return mockUsecase
			}(),
			requestBody:    `{"amount": 1000, "payment_schedules_id": [1]}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid loan ID format"}`,
		},
		{
			name:   "Internal Server Error",
			loanID: "1",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("MakePayment", mock.Anything, uint(1), []uint{1}, 1000.00).Return(errors.New("internal server error"))
				return mockUsecase
			}(),
			requestBody:    `{"amount": 1000, "payment_schedules_id": [1]}`,
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal server error"}`,
		},
		{
			name:   "Invalid Request Body",
			loanID: "1",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("MakePayment", mock.Anything, uint(1), []uint{1}, 1000.00).Return(errors.New("invalid request body"))
				return mockUsecase
			}(),
			requestBody:    `{"amount": 1000, "payment_schedules_id": [1],}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid request body"}`,
		},
		{
			name:   "Payment Schedule ID Required",
			loanID: "1",
			mockUsecase: func() *mocks.PaymentUsecase {
				mockUsecase := new(mocks.PaymentUsecase)
				mockUsecase.On("MakePayment", mock.Anything, uint(1), []uint{}, 1000.00).Return(errors.New("payment schedule ID is required"))
				return mockUsecase
			}(),
			requestBody:    `{"amount": 1000, "payment_schedules_id": []}`,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"payment schedule ID is required"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.mockUsecase)
			req, err := http.NewRequestWithContext(context.TODO(), "POST", "/payments/"+tt.loanID, bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
