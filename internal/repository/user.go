package repository

import (
	"context"
	"errors"
	"fmt"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/model"

	"github.com/jackc/pgx/v5"
)

type UserRepository struct {
	q *sqlc.Queries
}

func NewUserRepository(q *sqlc.Queries) *UserRepository {
	return &UserRepository{q: q}
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	sqlcUser, err := r.q.GetUserByEmail(ctx, email)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, fmt.Errorf("error getting user by email: %w", err)
	}
	id, err := pgUUIDToString(sqlcUser.ID)
	if err != nil {
		return model.User{}, fmt.Errorf("error converting uuid to string: %w", err)
	}
	createdAt := pgTimestamptzToTime(sqlcUser.CreatedAt)

	user := model.User{
		ID:           id,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    createdAt,
	}
	return user, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (model.User, error) {
	uuid, err := stringToPgUUID(id)
	if err != nil {
		return model.User{}, fmt.Errorf("error converting string to uuid: %w", err)
	}

	sqlcUser, err := r.q.GetUserByID(ctx, uuid)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.User{}, ErrNotFound
		}
		return model.User{}, fmt.Errorf("error getting user by id: %w", err)
	}

	createdAt := pgTimestamptzToTime(sqlcUser.CreatedAt)

	user := model.User{
		ID:           id,
		Email:        sqlcUser.Email,
		PasswordHash: sqlcUser.PasswordHash,
		CreatedAt:    createdAt,
	}
	return user, nil
}
