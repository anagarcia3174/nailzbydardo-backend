package testutil

import (
	"context"
	"nailzbydardo/internal/app"
	"nailzbydardo/internal/config"
	"nailzbydardo/internal/db"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
)

func NewTestServer(t *testing.T) (*httptest.Server, *pgxpool.Pool) {
	cfg, err := config.Load("../../.env.test")

	if err != nil {
		t.Fatalf("error loading config: %v", err)
	}

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		t.Fatalf("error connecting to database: %v", err)
	}

	_, err = pool.Exec(ctx, "TRUNCATE TABLE users, sessions, clients, services, expenses, appointments, appointment_services, appointment_discounts RESTART IDENTITY CASCADE")
	if err != nil {
		t.Fatalf("error truncating test database: %v", err)
	}
	mux := app.BuildHandlers(pool, cfg)
	server := httptest.NewServer(mux)

	t.Cleanup(func() {
		server.Close()
		pool.Close()
	})

	return server, pool
}

