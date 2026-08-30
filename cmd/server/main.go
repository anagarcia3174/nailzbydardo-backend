package main

import (
	"context"
	"fmt"
	"log"
	"nailzbydardo/internal/app"
	"nailzbydardo/internal/config"
	"nailzbydardo/internal/db"
	"net/http"
	"time"
)

func main() {
	cfg, err := config.Load(".env")

	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	ctx := context.Background()

	pool, err := db.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}
	defer pool.Close()

	mux := app.BuildHandlers(pool, cfg)

	addr := fmt.Sprintf(":%s", cfg.Port)
	server := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on port %s...", cfg.Port)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
