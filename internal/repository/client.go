package repository

import (
	"context"
	"errors"
	"fmt"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/model"
	"time"

	"github.com/jackc/pgx/v5"
)

type ClientRepository struct {
	q *sqlc.Queries
}

func NewClientRepository(q *sqlc.Queries) *ClientRepository {
	return &ClientRepository{q: q}
}

func (r *ClientRepository) CreateClient(ctx context.Context, name string, contactMethod *string, notes *string, birthday *time.Time) (model.Client, error) {
	pgContactMethod := stringPtrToPgText(contactMethod)
	pgNotes := stringPtrToPgText(notes)
	pgBirthday := timePtrToPgDate(birthday)

	createClientParams := sqlc.CreateClientParams{
		ClientName:    name,
		ContactMethod: pgContactMethod,
		Notes:         pgNotes,
		Birthday:      pgBirthday,
	}
	c, err := r.q.CreateClient(ctx, createClientParams)
	if err != nil {
		return model.Client{}, fmt.Errorf("error creating client: %w", err)
	}

	id, err := pgUUIDToString(c.ID)
	if err != nil {
		return model.Client{}, fmt.Errorf("error converting uuid to string: %w", err)
	}

	return model.Client{
		ID:            id,
		ClientName:    c.ClientName,
		ContactMethod: pgTextToStringPtr(c.ContactMethod),
		Notes:         pgTextToStringPtr(c.Notes),
		Birthday:      pgDateToTimePtr(c.Birthday),
		CreatedAt:     pgTimestamptzToTime(c.CreatedAt),
		DeletedAt:     pgTimestamptzToTimePtr(c.DeletedAt),
	}, nil
}

func (r *ClientRepository) GetClient(ctx context.Context, id string) (model.Client, error) {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return model.Client{}, ErrInvalidID
	}
	c, err := r.q.GetClient(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows){
			return model.Client{}, ErrNotFound
		}
		return model.Client{}, fmt.Errorf("error getting client: %w", err)
	}

	userID, err := pgUUIDToString(c.ID)
	if err != nil {
		return model.Client{}, fmt.Errorf("error converting uuid to string: %w", err)
	}

	return model.Client{
		ID:            userID,
		ClientName:    c.ClientName,
		ContactMethod: pgTextToStringPtr(c.ContactMethod),
		Notes:         pgTextToStringPtr(c.Notes),
		Birthday:      pgDateToTimePtr(c.Birthday),
		CreatedAt:     pgTimestamptzToTime(c.CreatedAt),
		DeletedAt:     pgTimestamptzToTimePtr(c.DeletedAt),
	}, nil
}

func (r *ClientRepository) ListClients(ctx context.Context) ([]model.ClientSummary, error) {
	clientRows, err := r.q.ListClients(ctx)
	if err != nil {
		return []model.ClientSummary{}, fmt.Errorf("error getting list of clients: %w", err)
	}
	clients := []model.ClientSummary{}
	for _, c := range clientRows {
		id, err := pgUUIDToString(c.ID)
		if err != nil {
			return []model.ClientSummary{}, fmt.Errorf("error converting uuid to string: %w", err)
		}
		client := model.ClientSummary{
			ID: id,
			ClientName: c.ClientName,
			ContactMethod: pgTextToStringPtr(c.ContactMethod),
		}
		clients = append(clients, client)
	}
	return clients, nil
}

func (r *ClientRepository) UpdateClient(ctx context.Context, id string, name string, contactMethod *string, notes *string, birthday *time.Time) (model.Client, error) {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return model.Client{}, ErrInvalidID
	}
	updateClientParams := sqlc.UpdateClientParams{
		ID: pgID,
		ClientName: name,
		ContactMethod: stringPtrToPgText(contactMethod),
		Notes: stringPtrToPgText(notes),
		Birthday: timePtrToPgDate(birthday),
	}
	c, err := r.q.UpdateClient(ctx, updateClientParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Client{}, ErrNotFound
		}
		return model.Client{}, fmt.Errorf("error updating client: %w", err)
	}
	userID, err := pgUUIDToString(c.ID)
	if err != nil {
		return model.Client{}, fmt.Errorf("error converting uuid to string: %w", err)
	}

	return model.Client{
		ID:            userID,
		ClientName:    c.ClientName,
		ContactMethod: pgTextToStringPtr(c.ContactMethod),
		Notes:         pgTextToStringPtr(c.Notes),
		Birthday:      pgDateToTimePtr(c.Birthday),
		CreatedAt:     pgTimestamptzToTime(c.CreatedAt),
		DeletedAt:     pgTimestamptzToTimePtr(c.DeletedAt),
	}, nil
}

func (r *ClientRepository) SoftDeleteClient(ctx context.Context, id string) error {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return ErrInvalidID
	}
	err = r.q.SoftDeleteClient(ctx, pgID)
	if err != nil {
		return fmt.Errorf("error soft deleting client: %w", err)
	}
	return nil
}

func (r *ClientRepository) ListClientAppointments(ctx context.Context, id string) ([]model.Appointment, error){
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return []model.Appointment{}, ErrInvalidID
	}
	clientAppointmentsRows, err := r.q.ListClientAppointments(ctx, pgID)
	if err != nil {
		return []model.Appointment{}, fmt.Errorf("error getting client's appointments: %w", err)
	}
	clientAppointments := []model.Appointment{}
	for _, a := range clientAppointmentsRows {
		id, err := pgUUIDToString(a.ID)
		if err != nil {
			return []model.Appointment{}, fmt.Errorf("error converting uuid to string: %w", err)
		}

		appt := model.Appointment{
			ID: id,
			ApptDate: pgTimestamptzToTime(a.ApptDate),
			ApptStatus: a.ApptStatus,
			LateFee: pgNumericToCentsPtr(a.LateFee),
			PaymentMethod: nullPaymentMethodToPtr(a.PaymentMethod),
			Notes: pgTextToStringPtr(a.Notes),
			LoyaltyReward: a.LoyaltyReward,
			Tip: pgNumericToCentsPtr(a.Tip),
			CreatedAt: pgTimestamptzToTime(a.CreatedAt),
		}
		clientAppointments = append(clientAppointments, appt)
	}
	return clientAppointments, nil
}