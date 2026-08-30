package repository

import (
	"context"
	"errors"
	"fmt"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/model"

	"github.com/jackc/pgx/v5"
)

type ServiceRepository struct {
	q *sqlc.Queries
}

func NewServiceRepository(q *sqlc.Queries) *ServiceRepository {
	return &ServiceRepository{q: q}
}

func (r *ServiceRepository) CreateService(ctx context.Context, name string, price int64) (model.Service, error) {
	createServiceParams := sqlc.CreateServiceParams{
		ServiceName: name,
		Price: centsToPgNumeric(price),
	}

	s, err := r.q.CreateService(ctx, createServiceParams)

	if err != nil {
		return model.Service{}, fmt.Errorf("error creating service: %w", err)
	}
	id, err := pgUUIDToString(s.ID)
	if err != nil {
		return model.Service{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	service := model.Service{
		ID: id,
		ServiceName: s.ServiceName,
		Price: pgNumericToCents(s.Price),
		CreatedAt: pgTimestamptzToTime(s.CreatedAt),
	}
	return service, nil
}

func (r *ServiceRepository) ListServices(ctx context.Context) ([]model.Service, error) {
	serviceRows, err := r.q.ListServices(ctx)

	if err != nil {
		return []model.Service{}, fmt.Errorf("error getting service rows: %w", err)
	}
	services := []model.Service{}
	for _, s := range serviceRows {
		id,  err := pgUUIDToString(s.ID)
		if err != nil {
			return []model.Service{}, fmt.Errorf("error converting uuid to string: %w", err)
		}
		service := model.Service{
			ID: id,
			ServiceName: s.ServiceName,
			Price: pgNumericToCents(s.Price),
			CreatedAt: pgTimestamptzToTime(s.CreatedAt),
		}
		services = append(services, service)
	}
	return services, nil
}

func (r *ServiceRepository) UpdateService(ctx context.Context, id string, name string, price int64) (model.Service, error) {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return model.Service{}, ErrInvalidID
	}
	updateServiceParams := sqlc.UpdateServiceParams{
		ID: pgID,
		ServiceName: name,
		Price: centsToPgNumeric(price),
	}
	s, err := r.q.UpdateService(ctx, updateServiceParams)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Service{}, ErrNotFound
		}
		return model.Service{}, fmt.Errorf("error updating service:  %w", err)
	}
	serviceID, err := pgUUIDToString(s.ID)
	if err != nil {
		return model.Service{}, fmt.Errorf("error converting uuid to id: %w", err)
	}
	service := model.Service{
		ID: serviceID,
		ServiceName: s.ServiceName,
		Price: pgNumericToCents(s.Price),
		CreatedAt: pgTimestamptzToTime(s.CreatedAt),
	}
	return service, nil
}

func (r *ServiceRepository) DeleteService(ctx context.Context, id string) error {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return ErrInvalidID
	}
	err = r.q.DeleteService(ctx, pgID)
	if err != nil {
		return fmt.Errorf("error deleting service: %w",  err)
	}
	return nil
}
