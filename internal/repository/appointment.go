package repository

import (
	"context"
	"errors"
	"fmt"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type AppointmentRepository struct {
	q *sqlc.Queries
}

func NewAppointmentRepository(q *sqlc.Queries) *AppointmentRepository {
	return &AppointmentRepository{q: q}
}

func (r *AppointmentRepository) CreateAppointment(ctx context.Context, clientID string, appointmentDate time.Time, notes *string) (model.Appointment, error) {
	pgID, err := stringToPgUUID(clientID)
	if err != nil {
		return model.Appointment{}, ErrInvalidID
	}
	createAppointmentParams := sqlc.CreateAppointmentParams{
		ClientID: pgID,
		ApptDate: timeToPgTimestamptz(appointmentDate),
		Notes:    stringPtrToPgText(notes),
	}
	a, err := r.q.CreateAppointment(ctx, createAppointmentParams)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("error creating appointment: %w", err)
	}
	appointmentID, err := pgUUIDToString(a.ID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	clientID, err = pgUUIDToString(a.ClientID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	appointment := model.Appointment{
		ID:            appointmentID,
		ClientID:      clientID,
		ApptDate:      pgTimestamptzToTime(a.ApptDate),
		ApptStatus:    a.ApptStatus,
		LateFee:       pgNumericToCentsPtr(a.LateFee),
		PaymentMethod: nullPaymentMethodToPtr(a.PaymentMethod),
		Notes:         pgTextToStringPtr(a.Notes),
		ReceiptURL:    pgTextToStringPtr(a.ReceiptUrl),
		LoyaltyReward: a.LoyaltyReward,
		Tip:           pgNumericToCentsPtr(a.Tip),
		CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
	}
	return appointment, nil
}

func (r *AppointmentRepository) GetAppointment(ctx context.Context, id string) (model.Appointment, error) {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return model.Appointment{}, ErrInvalidID
	}
	a, err := r.q.GetAppointment(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Appointment{}, ErrNotFound
		}
		return model.Appointment{}, fmt.Errorf("error getting appointment: %w", err)
	}
	appointmentID, err := pgUUIDToString(a.ID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	clientID, err := pgUUIDToString(a.ClientID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	appointment := model.Appointment{
		ID:            appointmentID,
		ClientID:      clientID,
		ApptDate:      pgTimestamptzToTime(a.ApptDate),
		ApptStatus:    a.ApptStatus,
		LateFee:       pgNumericToCentsPtr(a.LateFee),
		PaymentMethod: nullPaymentMethodToPtr(a.PaymentMethod),
		Notes:         pgTextToStringPtr(a.Notes),
		ReceiptURL:    pgTextToStringPtr(a.ReceiptUrl),
		LoyaltyReward: a.LoyaltyReward,
		Tip:           pgNumericToCentsPtr(a.Tip),
		CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
	}
	return appointment, nil
}

func (r *AppointmentRepository) ListAppointments(ctx context.Context) ([]model.AppointmentWithClient, error) {
	appointmentRows, err := r.q.ListAppointmentsWithClient(ctx)
	if err != nil {
		return []model.AppointmentWithClient{}, fmt.Errorf("error getting appointments:  %w", err)
	}
	appointments := []model.AppointmentWithClient{}
	for _, a := range appointmentRows {
		appointmentID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.AppointmentWithClient{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		appointment := model.AppointmentWithClient{
			ID:         appointmentID,
			ApptDate:   pgTimestamptzToTime(a.ApptDate),
			ApptStatus: a.ApptStatus,
			ClientName: a.ClientName,
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *AppointmentRepository) ListAppointmentsByDateRange(ctx context.Context, dateOne time.Time, dateTwo time.Time) ([]model.Appointment, error) {
	listAppointmentsByDateRangeParams := sqlc.ListAppointmentsByDateRangeParams{
		ApptDate:   timeToPgTimestamptz(dateOne),
		ApptDate_2: timeToPgTimestamptz(dateTwo),
	}
	appointmentRows, err := r.q.ListAppointmentsByDateRange(ctx, listAppointmentsByDateRangeParams)
	if err != nil {
		return []model.Appointment{}, fmt.Errorf("error getting appointments:  %w", err)
	}
	appointments := []model.Appointment{}
	for _, a := range appointmentRows {
		appointmentID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		clientID, err := pgUUIDToString(a.ClientID)
		if err != nil {
			return []model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		appointment := model.Appointment{
			ID:            appointmentID,
			ClientID:      clientID,
			ApptDate:      pgTimestamptzToTime(a.ApptDate),
			ApptStatus:    a.ApptStatus,
			LateFee:       pgNumericToCentsPtr(a.LateFee),
			PaymentMethod: nullPaymentMethodToPtr(a.PaymentMethod),
			Notes:         pgTextToStringPtr(a.Notes),
			ReceiptURL:    pgTextToStringPtr(a.ReceiptUrl),
			LoyaltyReward: a.LoyaltyReward,
			Tip:           pgNumericToCentsPtr(a.Tip),
			CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *AppointmentRepository) ListCompleteAppointmentsForPeriod(ctx context.Context, dateOne time.Time, dateTwo time.Time) ([]model.AppointmentSummary, error) {
	listCompleteAppointmentsForPeriodParams := sqlc.ListCompleteAppointmentsForPeriodParams{
		ApptDate:   timeToPgTimestamptz(dateOne),
		ApptDate_2: timeToPgTimestamptz(dateTwo),
	}
	appointmentRows, err := r.q.ListCompleteAppointmentsForPeriod(ctx, listCompleteAppointmentsForPeriodParams)
	if err != nil {
		return []model.AppointmentSummary{}, fmt.Errorf("error getting complete appointments:  %w", err)
	}
	appointments := []model.AppointmentSummary{}
	for _, a := range appointmentRows {
		appointmentID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.AppointmentSummary{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		appointment := model.AppointmentSummary{
			ID:       appointmentID,
			ApptDate: pgTimestamptzToTime(a.ApptDate),
			Tip:      pgNumericToCentsPtr(a.Tip),
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *AppointmentRepository) ListUpcomingAppointments(ctx context.Context) ([]model.Appointment, error) {
	appointmentRows, err := r.q.ListUpcomingAppointments(ctx)
	if err != nil {
		return []model.Appointment{}, fmt.Errorf("error getting upcoming appointments: %w", err)
	}
	appointments := []model.Appointment{}
	for _, a := range appointmentRows {
		appointmentID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		clientID, err := pgUUIDToString(a.ClientID)
		if err != nil {
			return []model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		appointment := model.Appointment{
			ID:            appointmentID,
			ClientID:      clientID,
			ApptDate:      pgTimestamptzToTime(a.ApptDate),
			ApptStatus:    a.ApptStatus,
			LateFee:       pgNumericToCentsPtr(a.LateFee),
			PaymentMethod: nullPaymentMethodToPtr(a.PaymentMethod),
			Notes:         pgTextToStringPtr(a.Notes),
			ReceiptURL:    pgTextToStringPtr(a.ReceiptUrl),
			LoyaltyReward: a.LoyaltyReward,
			Tip:           pgNumericToCentsPtr(a.Tip),
			CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *AppointmentRepository) ListAppointmentsForCalendar(ctx context.Context) ([]model.CalendarAppointment, error) {
	year := time.Now().UTC().Year()
	start := time.Date(year, time.January, 1, 0, 0, 0, 0, time.UTC)
	end := time.Date(year+1, time.January, 1, 0, 0, 0, 0, time.UTC)
	params := sqlc.ListAppointmentsForCalendarParams{
		ApptDate:   timeToPgTimestamptz(start),
		ApptDate_2: timeToPgTimestamptz(end),
	}
	appointmentRows, err := r.q.ListAppointmentsForCalendar(ctx, params)
	if err != nil {
		return []model.CalendarAppointment{}, fmt.Errorf("error getting upcoming appointments: %w", err)
	}
	appointments := []model.CalendarAppointment{}
	for _, a := range appointmentRows {
		appointmentID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.CalendarAppointment{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		appointment := model.CalendarAppointment{
			ID:         appointmentID,
			ApptDate:   pgTimestamptzToTime(a.ApptDate),
			Notes:      pgTextToStringPtr(a.Notes),
			ClientName: a.ClientName,
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *AppointmentRepository) ListUpcomingAppointmentsForDashboard(ctx context.Context) ([]model.AppointmentForDashboard, error) {
	appointmentRows, err := r.q.ListUpcomingAppointmentsForDashboard(ctx)
	if err != nil {
		return []model.AppointmentForDashboard{}, fmt.Errorf("error getting upcoming appointments: %w", err)
	}
	appointments := []model.AppointmentForDashboard{}
	for _, a := range appointmentRows {
		appointmentID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.AppointmentForDashboard{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		appointment := model.AppointmentForDashboard{
			ID:         appointmentID,
			ApptDate:   pgTimestamptzToTime(a.ApptDate),
			ClientName: a.ClientName,
		}
		appointments = append(appointments, appointment)
	}
	return appointments, nil
}

func (r *AppointmentRepository) GetAppointmentCountForPeriod(ctx context.Context, dateOne time.Time, dateTwo time.Time) (int64, error) {
	params := sqlc.GetAppointmentCountForPeriodParams{
		ApptDate:   timeToPgTimestamptz(dateOne),
		ApptDate_2: timeToPgTimestamptz(dateTwo),
	}
	count, err := r.q.GetAppointmentCountForPeriod(ctx, params)
	if err != nil {
		return 0, fmt.Errorf("error getting appointment count for period: %w", err)
	}
	return count, nil
}

func (r *AppointmentRepository) UpdateAppointment(ctx context.Context, id string, apptDate time.Time, apptStatus sqlc.AppointmentStatus, lateFee *int64, paymentMethod *sqlc.PaymentMethod, notes *string, receiptURL *string, loyaltyReward bool, tip *int64) (model.Appointment, error) {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return model.Appointment{}, ErrInvalidID
	}
	params := sqlc.UpdateAppointmentParams{
		ID:            pgID,
		ApptDate:      timeToPgTimestamptz(apptDate),
		ApptStatus:    apptStatus,
		LateFee:       centsPtrToPgNumeric(lateFee),
		PaymentMethod: ptrToNullPaymentMethod(paymentMethod),
		Notes:         stringPtrToPgText(notes),
		ReceiptUrl:    stringPtrToPgText(receiptURL),
		LoyaltyReward: loyaltyReward,
		Tip:           centsPtrToPgNumeric(tip),
	}
	a, err := r.q.UpdateAppointment(ctx, params)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Appointment{}, ErrNotFound
		}
		return model.Appointment{}, fmt.Errorf("error updating appointment: %w", err)
	}
	appointmentID, err := pgUUIDToString(a.ID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	clientID, err := pgUUIDToString(a.ClientID)
	if err != nil {
		return model.Appointment{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	appointment := model.Appointment{
		ID:            appointmentID,
		ClientID:      clientID,
		ApptDate:      pgTimestamptzToTime(a.ApptDate),
		ApptStatus:    a.ApptStatus,
		LateFee:       pgNumericToCentsPtr(a.LateFee),
		PaymentMethod: nullPaymentMethodToPtr(a.PaymentMethod),
		Notes:         pgTextToStringPtr(a.Notes),
		ReceiptURL:    pgTextToStringPtr(a.ReceiptUrl),
		LoyaltyReward: a.LoyaltyReward,
		Tip:           pgNumericToCentsPtr(a.Tip),
		CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
	}
	return appointment, nil

}

func (r *AppointmentRepository) DeleteAppointment(ctx context.Context, id string) error {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return ErrInvalidID
	}
	err = r.q.DeleteAppointment(ctx, pgID)
	if err != nil {
		return fmt.Errorf("error deleting appointment: %w", err)
	}
	return nil
}

func (r *AppointmentRepository) CreateAppointmentService(ctx context.Context, appointmentID string, serviceName string, servicePrice int64, designPrice int64) (model.AppointmentService, error) {
	pgID, err := stringToPgUUID(appointmentID)
	if err != nil {
		return model.AppointmentService{}, ErrInvalidID
	}
	params := sqlc.CreateAppointmentServiceParams{
		AppointmentID: pgID,
		ServiceName:   serviceName,
		ServicePrice:  centsToPgNumeric(servicePrice),
		DesignPrice:   centsToPgNumeric(designPrice),
	}
	a, err := r.q.CreateAppointmentService(ctx, params)
	if err != nil {
		return model.AppointmentService{}, fmt.Errorf("error creating appointment service: %w", err)
	}
	serviceID, err := pgUUIDToString(a.ID)
	if err != nil {
		return model.AppointmentService{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	apptID, err := pgUUIDToString(a.AppointmentID)
	if err != nil {
		return model.AppointmentService{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	service := model.AppointmentService{
		ID:            serviceID,
		AppointmentID: apptID,
		ServiceName:   a.ServiceName,
		ServicePrice:  pgNumericToCents(a.ServicePrice),
		DesignPrice:   pgNumericToCents(a.DesignPrice),
		CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
	}
	return service, nil
}

func (r *AppointmentRepository) ListAppointmentServicesByAppointment(ctx context.Context, appointmentID string) ([]model.AppointmentServiceSummary, error) {
	pgID, err := stringToPgUUID(appointmentID)
	if err != nil {
		return []model.AppointmentServiceSummary{}, ErrInvalidID
	}
	appointmentServiceRows, err := r.q.ListAppointmentServicesByAppointment(ctx, pgID)
	if err != nil {
		return []model.AppointmentServiceSummary{}, fmt.Errorf("error getting appointment services: %w", err)
	}
	appointmentServices := []model.AppointmentServiceSummary{}
	for _, a := range appointmentServiceRows {
		serviceID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.AppointmentServiceSummary{}, fmt.Errorf("error converting uuid to id: %w", err)
		}

		service := model.AppointmentServiceSummary{
			ID:           serviceID,
			ServiceName:  a.ServiceName,
			ServicePrice: pgNumericToCents(a.ServicePrice),
			DesignPrice:  pgNumericToCents(a.DesignPrice),
		}
		appointmentServices = append(appointmentServices, service)
	}
	return appointmentServices, nil
}

func (r *AppointmentRepository) DeleteAppointmentService(ctx context.Context, id string) error {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return ErrInvalidID
	}
	err = r.q.DeleteAppointmentService(ctx, pgID)
	if err != nil {
		return fmt.Errorf("error deleting appointment service: %w", err)
	}
	return nil
}

func (r *AppointmentRepository) CreateAppointmentDiscount(ctx context.Context, appointmentID string, discountName string, discountType sqlc.DiscountType, discountValue int64) (model.AppointmentDiscount, error) {
	pgID, err := stringToPgUUID(appointmentID)
	if err != nil {
		return model.AppointmentDiscount{}, ErrInvalidID
	}
	var pgValue pgtype.Numeric
	if discountType == sqlc.DiscountTypePercent {
		pgValue = intToPgNumeric(discountValue)
	} else {
		pgValue = centsToPgNumeric(discountValue)
	}
	params := sqlc.CreateAppointmentDiscountParams{
		AppointmentID: pgID,
		DiscountName:  discountName,
		DiscountType:  discountType,
		DiscountValue: pgValue,
	}
	a, err := r.q.CreateAppointmentDiscount(ctx, params)
	if err != nil {
		return model.AppointmentDiscount{}, fmt.Errorf("error creating appointment discount: %w", err)
	}
	discountID, err := pgUUIDToString(a.ID)
	if err != nil {
		return model.AppointmentDiscount{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	apptID, err := pgUUIDToString(a.AppointmentID)
	if err != nil {
		return model.AppointmentDiscount{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	if a.DiscountType == sqlc.DiscountTypePercent {
		discountValue = pgNumericToInt(a.DiscountValue)
	} else {
		discountValue = pgNumericToCents(a.DiscountValue)
	}
	discount := model.AppointmentDiscount{
		ID:            discountID,
		AppointmentID: apptID,
		DiscountName:  a.DiscountName,
		DiscountType:  a.DiscountType,
		DiscountValue: discountValue,
		CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
	}
	return discount, nil
}

func (r *AppointmentRepository) ListAppointmentDiscountsByAppointment(ctx context.Context, appointmentID string) ([]model.AppointmentDiscountSummary, error) {
	pgID, err := stringToPgUUID(appointmentID)
	if err != nil {
		return []model.AppointmentDiscountSummary{}, ErrInvalidID
	}
	appointmentDiscountRows, err := r.q.ListAppointmentDiscountsByAppointment(ctx, pgID)
	if err != nil {
		return []model.AppointmentDiscountSummary{}, fmt.Errorf("error getting appointment discounts: %w", err)
	}
	appointmentDiscounts := []model.AppointmentDiscountSummary{}
	for _, a := range appointmentDiscountRows {
		discountID, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.AppointmentDiscountSummary{}, fmt.Errorf("error converting uuid to id: %w", err)
		}
		var value int64
		if a.DiscountType == sqlc.DiscountTypePercent {
			value = pgNumericToInt(a.DiscountValue)
		} else {
			value = pgNumericToCents(a.DiscountValue)
		}
		discount := model.AppointmentDiscountSummary{
			ID:            discountID,
			DiscountName:  a.DiscountName,
			DiscountType:  a.DiscountType,
			DiscountValue: value,
		}
		appointmentDiscounts = append(appointmentDiscounts, discount)
	}
	return appointmentDiscounts, nil
}

func (r *AppointmentRepository) DeleteAppointmentDiscount(ctx context.Context, id string) error {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return ErrInvalidID
	}
	err = r.q.DeleteAppointmentDiscount(ctx, pgID)
	if err != nil {
		return fmt.Errorf("error deleting appointment discount: %w", err)
	}
	return nil
}

func (r *AppointmentRepository) GetAppointmentDetail(ctx context.Context, id string) (model.AppointmentDetail, error) {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return model.AppointmentDetail{}, ErrInvalidID
	}

	a, err := r.q.GetAppointmentWithClient(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.AppointmentDetail{}, ErrNotFound
		}

		return model.AppointmentDetail{}, fmt.Errorf("error getting appointment detail: %w", err)
	}

	appointmentID, err := pgUUIDToString(a.ID)
	if err != nil {
		return model.AppointmentDetail{}, fmt.Errorf("error converting uuid to id: %w", err)
	}

	clientID, err := pgUUIDToString(a.ClientID)
	if err != nil {
		return model.AppointmentDetail{}, fmt.Errorf("error converting uuid to id: %w", err)
	}

	appointment := model.Appointment{
		ID:            appointmentID,
		ClientID:      clientID,
		ApptDate:      pgTimestamptzToTime(a.ApptDate),
		ApptStatus:    a.ApptStatus,
		LateFee:       pgNumericToCentsPtr(a.LateFee),
		PaymentMethod: nullPaymentMethodToPtr(a.PaymentMethod),
		Notes:         pgTextToStringPtr(a.Notes),
		ReceiptURL:    pgTextToStringPtr(a.ReceiptUrl),
		LoyaltyReward: a.LoyaltyReward,
		Tip:           pgNumericToCentsPtr(a.Tip),
		CreatedAt:     pgTimestamptzToTime(a.CreatedAt),
	}

	client := model.ClientSummary{
		ID:            clientID,
		ClientName:    a.ClientName,
		ContactMethod: pgTextToStringPtr(a.ContactMethod),
	}

	services, err := r.ListAppointmentServicesByAppointment(ctx, id)
	if err != nil {
		return model.AppointmentDetail{}, err
	}

	discounts, err := r.ListAppointmentDiscountsByAppointment(ctx, id)
	if err != nil {
		return model.AppointmentDetail{}, err
	}
	appointmentDetail := model.AppointmentDetail{
		Appointment:          appointment,
		AppointmentServices:  services,
		AppointmentDiscounts: discounts,
		Client:               client,
		AppointmentRank:      a.AppointmentRank,
	}

	return appointmentDetail, nil
}
