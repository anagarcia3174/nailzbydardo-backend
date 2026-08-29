package service

import (
	"context"
	"errors"
	"fmt"
	"nailzbydardo/internal/model"
	"nailzbydardo/internal/repository"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// AuthService struct, holding:
// - *repository.UserRepository
// - *repository.SessionRepository
// Constructor: NewAuthService(userRepo, sessionRepo) *AuthService

type AuthService struct {
	userRepo    *repository.UserRepository
	sessionRepo *repository.SessionRepository
}

func NewAuthService(userRepo *repository.UserRepository, sessionRepo *repository.SessionRepository) *AuthService {
	return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo}
}

func (s *AuthService) Login(ctx context.Context, email string, password string) (model.Session, error) {
	user, err := s.userRepo.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return model.Session{}, ErrInvalidCredentials
		}
		return model.Session{}, fmt.Errorf("error getting user during login: %w", err)
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return model.Session{}, ErrInvalidCredentials
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)

	session, err := s.sessionRepo.CreateSession(ctx, user.ID, expiresAt)
	if err != nil {
		return model.Session{}, fmt.Errorf("error creating session: %w", err)
	}
	return session, nil
}

func (s *AuthService) Logout(ctx context.Context, id string) error {
	return s.sessionRepo.DeleteSession(ctx, id)
}

func (s *AuthService) ValidateSession(ctx context.Context, id string) (model.User, error) {
	session, err := s.sessionRepo.GetSessionByID(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound){
			return model.User{}, ErrUnauthorized
		}
		return model.User{}, fmt.Errorf("error getting session by id during validation: %w", err)
	}
	if session.ExpiresAt.Before(time.Now()) {
		return model.User{}, ErrUnauthorized
	}
	user, err := s.userRepo.GetUserByID(ctx, session.UserID)
	if err != nil {
		return model.User{}, fmt.Errorf("error getting user during session validation: %w", err)
	}
	return user, nil
}

func (s *AuthService) CleanupExpiredSessions(ctx context.Context) (int64, error) {
    return s.sessionRepo.DeleteExpiredSessions(ctx)
}