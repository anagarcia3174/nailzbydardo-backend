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

type SessionRepository struct {
	q *sqlc.Queries
}

func NewSessionRepository(q *sqlc.Queries) *SessionRepository {
	return &SessionRepository{q: q}
}

func (r *SessionRepository) CreateSession(ctx context.Context, userID string, expiresAt time.Time) (model.Session, error) {
	pgID, err := stringToPgUUID(userID)
	if err != nil {
		return model.Session{}, fmt.Errorf("error converting string to uuid: %w", err)
	}
	timestamptz := timeToPgTimestamptz(expiresAt)

	createSessionParams := sqlc.CreateSessionParams{
		UserID:    pgID,
		ExpiresAt: timestamptz,
	}

	s, err := r.q.CreateSession(ctx, createSessionParams)
	if err != nil {
		return model.Session{}, fmt.Errorf("error creating session: %w", err)
	}
	id, err := pgUUIDToString(s.ID)
	if err != nil {
		return model.Session{}, fmt.Errorf("error converting uuid to string: %w", err)
	}
	stringUserID, err := pgUUIDToString(s.UserID)
	if err != nil {
		return model.Session{}, fmt.Errorf("error converting uuid to string: %w", err)
	}
	expiresTime := pgTimestamptzToTime(s.ExpiresAt)
	createdAt := pgTimestamptzToTime(s.CreatedAt)
	session := model.Session{
		ID:        id,
		UserID:    stringUserID,
		ExpiresAt: expiresTime,
		CreatedAt: createdAt,
	}
	return session, nil
}

func (r *SessionRepository) GetSessionByID(ctx context.Context, id string) (model.Session, error) {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return model.Session{}, fmt.Errorf("error converting string to uuid: %w", err)
	}

	s, err := r.q.GetSessionByID(ctx, pgID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Session{}, ErrNotFound
		}
		return model.Session{}, fmt.Errorf("error getting session by id: %w", err)
	}
	sessionID, err := pgUUIDToString(s.ID)
	if err != nil {
		return model.Session{}, fmt.Errorf("error converting uuid to string: %w", err)
	}
	stringUserID, err := pgUUIDToString(s.UserID)
	if err != nil {
		return model.Session{}, fmt.Errorf("error converting uuid to string: %w", err)
	}
	expiresTime := pgTimestamptzToTime(s.ExpiresAt)
	createdAt := pgTimestamptzToTime(s.CreatedAt)
	session := model.Session{
		ID:        sessionID,
		UserID:    stringUserID,
		ExpiresAt: expiresTime,
		CreatedAt: createdAt,
	}
	return session, nil
}

func (r *SessionRepository) DeleteSession(ctx context.Context, id string) error {
	pgID, err := stringToPgUUID(id)
	if err != nil {
		return fmt.Errorf("error converting string to uuid: %w", err)
	}

	err = r.q.DeleteSession(ctx, pgID)
	if err != nil {
		return fmt.Errorf("error deleting session: %w", err)
	}
	return nil
}

func (r *SessionRepository) DeleteExpiredSessions(ctx context.Context) (int64, error) {
	return r.q.DeleteExpiredSessions(ctx)
}
