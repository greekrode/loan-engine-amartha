package dto

import "time"

type CreateLoanRequest struct {
	BorrowerID   uint    `json:"borrower_id"`
	Principal    float64 `json:"principal"`
	InterestRate float64 `json:"interest_rate"`
	Duration     int     `json:"duration"`
	StartDate    string  `json:"start_date"`
}

type CreateLoanResponse struct {
	ID                uint                         `json:"id"`
	Principal         float64                      `json:"principal"`
	InterestRate      float64                      `json:"interest_rate"`
	Duration          int                          `json:"duration"`
	StartDate         time.Time                    `json:"start_date"`
	OutstandingAmount float64                      `json:"outstanding_amount"`
	PaymentSchedules  []GetPaymentScheduleResponse `json:"payment_schedules"`
}

type GetLoanDetailsResponse struct {
	Principal         float64                      `json:"principal"`
	InterestRate      float64                      `json:"interest_rate"`
	OutstandingAmount float64                      `json:"outstanding_amount"`
	Duration          int                          `json:"duration"`
	StartDate         time.Time                    `json:"start_date"`
	CreatedAt         time.Time                    `json:"created_at"`
	Borrower          GetBorrowerResponse          `json:"borrower"`
	PaymentSchedule   []GetPaymentScheduleResponse `json:"payment_schedules"`
}

type GetOutstandingResponse struct {
	OutstandingAmount float64 `json:"outstanding_amount"`
}
