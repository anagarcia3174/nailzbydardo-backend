package service

import (
	"context"
	"nailzbydardo/internal/model"
	"nailzbydardo/internal/repository"
	"strings"
)

type CatalogService struct {
	serviceRepo *repository.ServiceRepository
}

func NewCatalogService(serviceRepo *repository.ServiceRepository) *CatalogService {
	return &CatalogService{serviceRepo: serviceRepo}
}

func (s *CatalogService) CreateService(ctx context.Context, name string, price int64) (model.Service, error) {
	if strings.TrimSpace(name) == "" {
		return model.Service{}, ErrInvalidInput
	}
	return s.serviceRepo.CreateService(ctx, name, price)
}

func (s *CatalogService) ListServices(ctx context.Context) ([]model.Service, error) {
	return s.serviceRepo.ListServices(ctx)
}

func (s *CatalogService) UpdateService(ctx context.Context, id string, name string, price int64) (model.Service, error) {
	if strings.TrimSpace(name) == "" {
		return model.Service{}, ErrInvalidInput
	}
	return s.serviceRepo.UpdateService(ctx, id, name, price)
}

func (s *CatalogService) DeleteService(ctx context.Context, id string) error {
	return s.serviceRepo.DeleteService(ctx, id)
}
