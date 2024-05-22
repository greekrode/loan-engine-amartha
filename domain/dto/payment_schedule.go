package dto

import "time"

type GetPaymentScheduleResponse struct {
	ID        uint      `json:"id,omitempty"`
	DueAmount float64   `json:"due_amount"`
	DueDate   time.Time `json:"due_date"`
	Paid      bool      `json:"paid"`
}

type MakePaymentRequest struct {
	PaymentSchedulesID []uint  `json:"payment_schedules_id"`
	Amount             float64 `json:"amount"`
}
