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
	loanHttp "github.com/greekrode/loan-engine-amartha/loan/delivery/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRouter(mockUCase *mocks.LoanUsecase) *gin.Engine {
	router := gin.Default()
	router.GET("/loans/:loan_id", func(c *gin.Context) {
		handler := loanHttp.LoanHandler{
			LoanUsecase: mockUCase,
		}
		handler.GetLoanDetails(c)
	})
	router.GET("/loans/:loan_id/outstanding", func(c *gin.Context) {
		handler := loanHttp.LoanHandler{
			LoanUsecase: mockUCase,
		}
		handler.GetOutstanding(c)
	})
	router.POST("/loans", func(c *gin.Context) {
		handler := loanHttp.LoanHandler{
			LoanUsecase: mockUCase,
		}
		handler.CreateLoan(c)
	})
	return router
}

func TestGetLoanDetails(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		loanID         string
		mockUsecase    *mocks.LoanUsecase
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Valid Loan Details",
			loanID: "1",
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("GetLoanDetails", mock.Anything, uint(1)).Return(&dto.GetLoanDetailsResponse{
					Principal:         100,
					InterestRate:      10,
					OutstandingAmount: 1000,
					Duration:          52,
					StartDate:         time.Time{},
					CreatedAt:         time.Time{},
					Borrower: dto.GetBorrowerResponse{
						ID:        1,
						FirstName: "John",
						LastName:  "Doe",
						Email:     "john.doe@example.com",
						CreatedAt: time.Time{},
					},
					PaymentSchedule: []dto.GetPaymentScheduleResponse{
						{
							DueAmount: 10,
							DueDate:   time.Time{},
							Paid:      false,
						},
					},
				},
					nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"principal": 100,
				"interest_rate": 10,
				"outstanding_amount": 1000,
				"duration": 52,
				"start_date": "0001-01-01T00:00:00Z",
				"created_at": "0001-01-01T00:00:00Z",
				"borrower": {
					"id": 1,
					"first_name": "John",
					"last_name": "Doe",
					"email": "john.doe@example.com",
					"created_at": "0001-01-01T00:00:00Z"
				},
				"payment_schedules": [
					{
						"due_amount": 10,
						"due_date": "0001-01-01T00:00:00Z",
						"paid": false
					}
				]
			}`,
		},
		{
			name:   "Invalid Loan ID",
			loanID: "abc",
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("GetLoanDetails", mock.Anything, uint(1)).Return(&dto.GetLoanDetailsResponse{}, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid loan ID format"}`,
		},
		{
			name:   "Loan Usecase Error",
			loanID: "1",
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("GetLoanDetails", mock.Anything, uint(1)).Return(&dto.GetLoanDetailsResponse{}, errors.New("internal error"))
				return mockUsecase
			}(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.mockUsecase)
			req, err := http.NewRequestWithContext(context.TODO(), "GET", "/loans/"+tt.loanID, bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestGetOutstanding(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		loanID         string
		mockUsecase    *mocks.LoanUsecase
		expectedStatus int
		expectedBody   string
	}{
		{
			name:   "Valid Outstanding",
			loanID: "1",
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("GetOutstandingAmount", mock.Anything, uint(1)).Return(1000.00, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"outstanding_amount": 1000
			}`,
		},
		{
			name:   "Invalid Loan ID",
			loanID: "abc",
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("GetOutstandingAmount", mock.Anything, uint(1)).Return(0.00, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid loan ID format"}`,
		},
		{
			name:   "Loan Usecase Error",
			loanID: "1",
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("GetOutstandingAmount", mock.Anything, uint(1)).Return(0.00, errors.New("internal error"))
				return mockUsecase
			}(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.mockUsecase)
			req, err := http.NewRequestWithContext(context.TODO(), "GET", "/loans/"+tt.loanID+"/outstanding", bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestCreateLoan(t *testing.T) {
	gin.SetMode(gin.TestMode)

	fixedTime := time.Date(2023, time.January, 1, 0, 0, 0, 0, time.UTC)

	tests := []struct {
		name           string
		mockUsecase    *mocks.LoanUsecase
		requestBody    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Valid Create Loan",
			requestBody: `{
				"borrower_id": 1,
				"principal": 100.00,
				"interest_rate": 10.00,
				"duration": 52,
				"start_date": "2023-01-01"
			}`,
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("CreateLoan", mock.Anything, uint(1), 100.00, 10.00, int32(52), fixedTime).Return(&dto.CreateLoanResponse{
					ID:                1,
					Principal:         100.00,
					InterestRate:      10.00,
					Duration:          52,
					StartDate:         fixedTime,
					OutstandingAmount: 1000,
					PaymentSchedules: []dto.GetPaymentScheduleResponse{
						{
							ID:        1,
							DueAmount: 100,
							DueDate:   fixedTime,
						},
					},
				}, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusOK,
			expectedBody: `{
				"id": 1,
				"principal": 100.00,
				"interest_rate": 10.00,
				"duration": 52,
				"outstanding_amount": 1000,
				"start_date": "2023-01-01T00:00:00Z",
				"payment_schedules": [
					{
						"id": 1,
						"due_amount": 100.00,
						"due_date": "2023-01-01T00:00:00Z",
						"paid": false
					}
				]
			}`,
		},
		{
			name: "Invalid JSON Request",
			requestBody: `{
				"borrower_id": 1,
			}`,
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("CreateLoan", mock.Anything, uint(1), 100.00, 10.00, int32(52), fixedTime).Return(&dto.CreateLoanResponse{}, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid request body"}`,
		},
		{
			name: "Invalid Start Date",
			requestBody: `{
				"borrower_id" : 1,
				"principal": 100.00,
				"interest_rate": 10.00,
				"duration": 52,
				"start_date": "23-01-01"
			}`,
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("CreateLoan", mock.Anything, uint(1), 100.00, 10.00, int32(52), fixedTime).Return(&dto.CreateLoanResponse{}, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid date format, should be YYYY-MM-DD"}`,
		},
		{
			name: "Loan Usecase Error",
			requestBody: `{
				"borrower_id": 1,
				"principal": 100.00,
				"interest_rate": 10.00,
				"duration": 52,
				"start_date": "2023-01-01"
			}`,
			mockUsecase: func() *mocks.LoanUsecase {
				mockUsecase := new(mocks.LoanUsecase)
				mockUsecase.On("CreateLoan", mock.Anything, uint(1), 100.00, 10.00, int32(52), fixedTime).Return(nil, errors.New("internal error"))
				return mockUsecase
			}(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.mockUsecase)
			req, err := http.NewRequestWithContext(context.TODO(), "POST", "/loans", bytes.NewBufferString(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
