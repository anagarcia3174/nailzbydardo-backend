package service

import (
	"context"
	"nailzbydardo/internal/model"
	"nailzbydardo/internal/repository"
	"strings"
	"time"
)

type ClientService struct {
	clientRepo *repository.ClientRepository
}

func NewClientService(clientRepo *repository.ClientRepository) *ClientService {
	return &ClientService{clientRepo: clientRepo}
}

func (s *ClientService) CreateClient(ctx context.Context, name string, contactMethod *string, notes *string, birthday *time.Time) (model.Client, error) {
	if strings.TrimSpace(name) == "" {
		return model.Client{}, ErrInvalidInput
	}
	return s.clientRepo.CreateClient(ctx, name, contactMethod, notes, birthday)
}

func (s *ClientService) GetClient(ctx context.Context, id string) (model.Client, error) {
	return s.clientRepo.GetClient(ctx, id)
}

func (s *ClientService) ListClients(ctx context.Context) ([]model.ClientSummary, error) {
	return s.clientRepo.ListClients(ctx)
}

func (s *ClientService) UpdateClient(ctx context.Context, id string, name string, contactMethod *string, notes *string, birthday *time.Time) (model.Client, error) {
	if strings.TrimSpace(name) == "" {
		return model.Client{}, ErrInvalidInput
	}
	return s.clientRepo.UpdateClient(ctx, id, name, contactMethod, notes, birthday)
}

func (s *ClientService) SoftDeleteClient(ctx context.Context, id string) error {
	return s.clientRepo.SoftDeleteClient(ctx, id)
}

func (s *ClientService) GetClientAppointments(ctx context.Context, id string) ([]model.Appointment, error) {
	return s.clientRepo.ListClientAppointments(ctx, id)
}