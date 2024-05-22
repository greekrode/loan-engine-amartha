package dto

type RequestPaymentResponse struct {
	TotalDue         float64                      `json:"total_due"`
	PaymentSchedules []GetPaymentScheduleResponse `json:"payment_schedules"`
}
