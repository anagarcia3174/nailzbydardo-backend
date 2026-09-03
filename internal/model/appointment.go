package model

import (
	"nailzbydardo/internal/db/sqlc"
	"time"
)

type Appointment struct {
	ID            string                 `json:"id"`
	ClientID      string                 `json:"client_id"`
	ApptDate      time.Time              `json:"appt_date"`
	ApptStatus    sqlc.AppointmentStatus `json:"appt_status"`
	LateFee       *int64                 `json:"late_fee"`
	PaymentMethod *sqlc.PaymentMethod    `json:"payment_method"`
	Notes         *string                `json:"notes"`
	ReceiptURL    *string                `json:"receipt_url"`
	LoyaltyReward bool                   `json:"loyalty_reward"`
	Tip           *int64                 `json:"tip"`
	CreatedAt     time.Time              `json:"created_at"`
}
type AppointmentWithClient struct {
	ID         string                 `json:"id"`
	ApptDate   time.Time              `json:"appt_date"`
	ApptStatus sqlc.AppointmentStatus `json:"appt_status"`
	ClientName string                 `json:"client_name"`
}
type AppointmentService struct {
	ID            string    `json:"id"`
	AppointmentID string    `json:"appointment_id"`
	ServiceName   string    `json:"service_name"`
	ServicePrice  int64     `json:"service_price"`
	DesignPrice   int64     `json:"design_price"`
	CreatedAt     time.Time `json:"created_at"`
}
type AppointmentDiscount struct {
	ID            string            `json:"id"`
	AppointmentID string            `json:"appointment_id"`
	DiscountName  string            `json:"discount_name"`
	DiscountType  sqlc.DiscountType `json:"discount_type"`
	DiscountValue int64             `json:"discount_value"`
	CreatedAt     time.Time         `json:"created_at"`
}
type AppointmentDetail struct {
	Appointment          Appointment                  `json:"appointment"`
	AppointmentServices  []AppointmentServiceSummary  `json:"appointment_services"`
	AppointmentDiscounts []AppointmentDiscountSummary `json:"appointment_discounts"`
	Client               ClientSummary                `json:"client_summary"`
	AppointmentRank      int32                        `json:"appointment_rank"`
}

type AppointmentSummary struct {
	ID       string    `json:"id"`
	ApptDate time.Time `json:"appt_date"`
	Tip      *int64    `json:"tip"`
}

type AppointmentServiceSummary struct {
	ID           string `json:"id"`
	ServiceName  string `json:"service_name"`
	ServicePrice int64  `json:"service_price"`
	DesignPrice  int64  `json:"design_price"`
}

type AppointmentDiscountSummary struct {
	ID            string            `json:"id"`
	DiscountName  string            `json:"discount_name"`
	DiscountType  sqlc.DiscountType `json:"discount_type"`
	DiscountValue int64             `json:"discount_value"`
}

type AppointmentForDashboard struct {
	ID         string    `json:"id"`
	ApptDate   time.Time `json:"appt_date"`
	ClientName string    `json:"client_name"`
}

type CalendarAppointment struct {
	ID         string    `json:"id"`
	ApptDate   time.Time `json:"appt_date"`
	Notes      *string   `json:"notes"`
	ClientName string    `json:"client_name"`
}
