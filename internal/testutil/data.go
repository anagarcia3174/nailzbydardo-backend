package testutil

import (
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

func CreateTestClient(t *testing.T, pool *pgxpool.Pool, name string, contactMethod *string, notes *string, birthday *time.Time) string {
	var id string

	err := pool.QueryRow(
		t.Context(),
		`INSERT INTO clients (client_name, contact_method, notes, birthday)
         VALUES ($1, $2, $3, $4)
         RETURNING id`,
		name, contactMethod, notes, birthday,
	).Scan(&id)

	if err != nil {
		t.Fatalf("error creating test client: %v", err)
	}

	return id
}