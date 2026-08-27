package main

import (
	"context"
	"fmt"
	"log"
	"nailzbydardo/internal/config"
	"nailzbydardo/internal/db"
	"nailzbydardo/internal/handler"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer pool.Close()

	r := chi.NewRouter()
	healthHandler := handler.NewHealthHandler(pool)

	r.Get("/health", healthHandler.ServeHTTP)

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr: addr,
		Handler: r,
		ReadTimeout:    10 * time.Second,  // Max duration for reading the entire request
		WriteTimeout:   10 * time.Second,
	}

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}