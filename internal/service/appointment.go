package service

import (
	"context"
	"errors"
	"fmt"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/model"
	"nailzbydardo/internal/repository"
	"time"
)

type AppointmentTotal struct {
	Subtotal      int64
	DiscountTotal int64
	ServiceTotal  int64
	Tip           int64
	GrandTotal    int64
}

type AppointmentService struct {
	appointmentRepo *repository.AppointmentRepository
	clientRepo *repository.ClientRepository
}

func NewAppointmentService(appointmentRepo *repository.AppointmentRepository, clientRepo *repository.ClientRepository) *AppointmentService{
	return &AppointmentService{appointmentRepo: appointmentRepo, clientRepo: clientRepo}
}

func (s *AppointmentService) CreateAppointment(ctx context.Context, clientID string, apptDate time.Time, notes *string) (model.Appointment, error) {
	_, err := s.clientRepo.GetClient(ctx, clientID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Appointment{}, ErrNonexistentClient
		}
		return model.Appointment{}, fmt.Errorf("error getting client during appointment creation: %w", err)
	}
	return s.appointmentRepo.CreateAppointment(ctx, clientID, apptDate, notes)
}

func (s *AppointmentService) GetAppointment(ctx context.Context, id string) (model.Appointment, error) {
	return s.appointmentRepo.GetAppointment(ctx, id)
}
func (s *AppointmentService) ListAppointments(ctx context.Context) ([]model.Appointment, error) {
	return s.appointmentRepo.ListAppointments(ctx)
}
func (s *AppointmentService) ListAppointmentsByDateRange(ctx context.Context, dateOne time.Time, dateTwo time.Time) ([]model.Appointment, error) {
	return s.appointmentRepo.ListAppointmentsByDateRange(ctx, dateOne, dateTwo)
}
func (s *AppointmentService) ListUpcomingAppointments(ctx context.Context) ([]model.Appointment, error) {
	return s.appointmentRepo.ListUpcomingAppointments(ctx)
}
func (s *AppointmentService) GetAppointmentDetail(ctx context.Context, id string) (model.AppointmentDetail, error) {
	return s.appointmentRepo.GetAppointmentDetail(ctx, id)
}

func (s *AppointmentService) UpdateAppointment(ctx context.Context, id string, apptDate time.Time, apptStatus sqlc.AppointmentStatus, lateFee *int64, paymentMethod *sqlc.PaymentMethod, notes *string, receiptURL *string, loyaltyReward bool, tip *int64) (model.Appointment, error) {
	return s.appointmentRepo.UpdateAppointment(ctx, id, apptDate, apptStatus, lateFee, paymentMethod, notes, receiptURL, loyaltyReward, tip)
}

func (s *AppointmentService) DeleteAppointment(ctx context.Context, id string) error {
	return s.appointmentRepo.DeleteAppointment(ctx, id)
}


func (s *AppointmentService) AddAppointmentService(ctx context.Context, appointmentID string, serviceName string, servicePrice int64, designPrice int64) (model.AppointmentService, error) {
	return s.appointmentRepo.CreateAppointmentService(ctx, appointmentID, serviceName, servicePrice, designPrice)
}
//no need for validation as we are passing in the service's id not the appointments
func (s *AppointmentService) RemoveAppointmentService(ctx context.Context, id string) error {
	return s.appointmentRepo.DeleteAppointmentService(ctx, id)
}

func (s *AppointmentService) AddAppointmentDiscount(ctx context.Context, appointmentID string, discountName string, discountType sqlc.DiscountType, discountValue int64) (model.AppointmentDiscount, error) {
	if discountType == sqlc.DiscountTypePercent {
		if discountValue > 100 {
			return model.AppointmentDiscount{}, ErrInvalidInput
		}
	}
	return s.appointmentRepo.CreateAppointmentDiscount(ctx, appointmentID, discountName, discountType, discountValue)
}
func (s *AppointmentService) RemoveAppointmentDiscount(ctx context.Context, id string) error {
	return s.appointmentRepo.DeleteAppointmentDiscount(ctx, id)
}

// ---- CalculateAppointmentTotal ----
// NOT a passthrough — this is where the actual money math lives, the
// payoff for keeping revenue/discount calculations out of SQL.
// Takes an assembled model.AppointmentDetail (or fetches one internally
// via GetAppointmentDetail).
// Steps:
// 1. Sum all AppointmentServices' ServicePrice + DesignPrice -> subtotal.
// 2. Apply each AppointmentDiscount:
//    - type "amount": subtract discount_value directly (in cents).
//    - type "percent": subtract (subtotal_at_this_point * discount_value / 100).
//      Decide: do percent discounts apply to the running total after
//      previous discounts, or always against the original subtotal?
//      This is a real business decision if you ever allow multiple
//      discounts on one appointment — worth deciding explicitly.
// 3. Add tip separately if your definition of "total" includes it
//    (or keep tip reported separately, since tip isn't really part of
//    "service revenue" in most accounting contexts — worth deciding).
// 4. Return the computed total (int64 cents), possibly alongside a
//    breakdown (subtotal, discount amount, final total) if the
//    frontend wants to display each line.
func (s *AppointmentService) CalculateAppointmentTotal(ctx context.Context, appointmentID string) (AppointmentTotal, error) {
	detail, err := s.appointmentRepo.GetAppointmentDetail(ctx, appointmentID)
	if err != nil {
		return AppointmentTotal{}, err
	}

	var subtotal int64
	for _, svc := range detail.AppointmentServices {
		subtotal += svc.ServicePrice + svc.DesignPrice
	}

	var discountTotal int64
	for _, d := range detail.AppointmentDiscounts {
		switch d.DiscountType {
		case sqlc.DiscountTypeAmount:
			discountTotal += d.DiscountValue
		case sqlc.DiscountTypePercent:
			discountTotal += subtotal * d.DiscountValue / 100
		}
	}

	serviceTotal := subtotal - discountTotal
	if serviceTotal < 0 {
		serviceTotal = 0
	}

	var tip int64
	if detail.Appointment.Tip != nil {
		tip = *detail.Appointment.Tip
	}

	return AppointmentTotal{
		Subtotal:      subtotal,
		DiscountTotal: discountTotal,
		ServiceTotal:  serviceTotal,
		Tip:           tip,
		GrandTotal:    serviceTotal + tip,
	}, nil
}