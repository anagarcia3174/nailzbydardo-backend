package model

import "time"

type PaymentMethod string
const (
    PaymentMethodCash   PaymentMethod = "cash"
    PaymentMethodZelle   PaymentMethod = "zelle"
    PaymentMethodCashApp PaymentMethod = "cash_app"
    PaymentMethodOther PaymentMethod = "other"
)

type AppointmentStatus string
const (
	AppointmentStatusBooked    AppointmentStatus = "booked"
	AppointmentStatusComplete  AppointmentStatus = "complete"
	AppointmentStatusNoShow    AppointmentStatus = "no_show"
	AppointmentStatusCancelled AppointmentStatus = "cancelled"
)

type DiscountType string
const (
	DiscountTypeAmount  DiscountType = "amount"
	DiscountTypePercent DiscountType = "percent"
)

type Appointment struct {
	ID string `json:"id"`
	ClientID string `json:"client_id"`
	ApptDate time.Time `json:"appt_date"`
	ApptStatus    AppointmentStatus  `json:"appt_status"`
	LateFee *int64 `json:"late_fee"`
	PaymentMethod *PaymentMethod  `json:"payment_method"`
	Notes *string `json:"notes"`
	ReceiptURL *string `json:"receipt_url"`
	LoyaltyReward bool `json:"loyalty_reward"`
	Tip *int64 `json:"tip"`
	CreatedAt time.Time `json:"created_at"`
}
type AppointmentService struct {
	ID            string        `json:"id"`
	AppointmentID string        `json:"appointment_id"`
	ServiceName   string             `json:"service_name"`
	ServicePrice  int64     `json:"service_price"`
	DesignPrice   int64     `json:"design_price"`
	CreatedAt     time.Time `json:"created_at"`
}
type AppointmentDiscount struct {
	ID            string        `json:"id"`
	AppointmentID string        `json:"appointment_id"`
	DiscountName  string             `json:"discount_name"`
	DiscountType  DiscountType       `json:"discount_type"`
	DiscountValue int64     `json:"discount_value"`
	CreatedAt     time.Time `json:"created_at"`
}
type AppointmentDetail struct {
	Appointment Appointment `json:"appointment"`
	AppointmentServices []AppointmentService `json:"appointment_services"`
	AppointmentDiscounts []AppointmentDiscount `json:"appointment_discounts"`
}