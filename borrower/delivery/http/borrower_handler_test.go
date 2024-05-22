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
	borrowerHttp "github.com/greekrode/loan-engine-amartha/borrower/delivery/http"
	"github.com/greekrode/loan-engine-amartha/domain/dto"
	"github.com/greekrode/loan-engine-amartha/domain/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupRouter(mockUCase *mocks.BorrowerUsecase) *gin.Engine {
	router := gin.Default()

	router.GET("/borrowers/:borrower_id/status", func(c *gin.Context) {
		handler := borrowerHttp.BorrowerHandler{
			BorrowerUsecase: mockUCase,
		}
		handler.CheckDelinquent(c)
	})
	router.POST("/borrowers", func(c *gin.Context) {
		handler := borrowerHttp.BorrowerHandler{
			BorrowerUsecase: mockUCase,
		}
		handler.CreateBorrower(c)
	})
	return router
}

func TestCheckDelinquent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	tests := []struct {
		name           string
		borrowerID     string
		mockUsecase    *mocks.BorrowerUsecase
		expectedStatus int
		expectedBody   string
	}{
		{
			name:       "Valid Delinquent Check",
			borrowerID: "1",
			mockUsecase: func() *mocks.BorrowerUsecase {
				mockUsecase := new(mocks.BorrowerUsecase)
				mockUsecase.On("IsDelinquent", mock.Anything, uint(1)).Return(true, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"is_delinquent":true}`,
		},
		{
			name:       "Invalid Borrower ID",
			borrowerID: "abc",
			mockUsecase: func() *mocks.BorrowerUsecase {
				mockUsecase := new(mocks.BorrowerUsecase)
				mockUsecase.On("IsDelinquent", mock.Anything, uint(1)).Return(true, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"message":"invalid borrower ID format"}`,
		},
		{
			name:       "Borrower Usecase Error",
			borrowerID: "2",
			mockUsecase: func() *mocks.BorrowerUsecase {
				mockUsecase := new(mocks.BorrowerUsecase)
				mockUsecase.On("IsDelinquent", mock.Anything, uint(2)).Return(false, errors.New("internal error"))
				return mockUsecase
			}(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"internal error"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.mockUsecase)

			req, err := http.NewRequestWithContext(context.Background(), "GET", "/borrowers/"+tt.borrowerID+"/status", nil)
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}

func TestCreateBorrower(t *testing.T) {
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name           string
		mockUsecase    *mocks.BorrowerUsecase
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "Successful Creation",
			mockUsecase: func() *mocks.BorrowerUsecase {
				mockUsecase := new(mocks.BorrowerUsecase)
				mockUsecase.On("CreateBorrower", mock.Anything).Return(&dto.CreateBorrowerResponse{
					ID:        1,
					FirstName: "John",
					LastName:  "Doe",
					Email:     "john.doe@example.com",
					CreatedAt: time.Time{},
				}, nil)
				return mockUsecase
			}(),
			expectedStatus: http.StatusCreated,
			expectedBody: `{
				"id": 1,
				"first_name": "John",
				"last_name": "Doe",
				"email": "john.doe@example.com",
				"created_at": "0001-01-01T00:00:00Z"
			}`,
		},
		{
			name: "Failed Creation",
			mockUsecase: func() *mocks.BorrowerUsecase {
				mockUsecase := new(mocks.BorrowerUsecase)
				mockUsecase.On("CreateBorrower", mock.Anything).Return(nil, errors.New("creation failed"))
				return mockUsecase
			}(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"message":"creation failed"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := setupRouter(tt.mockUsecase)
			req, err := http.NewRequestWithContext(context.TODO(), "POST", "/borrowers", bytes.NewBufferString("{}"))
			req.Header.Set("Content-Type", "application/json")

			require.NoError(t, err)

			rec := httptest.NewRecorder()
			router.ServeHTTP(rec, req)

			assert.Equal(t, tt.expectedStatus, rec.Code)
			assert.JSONEq(t, tt.expectedBody, rec.Body.String())
		})
	}
}
