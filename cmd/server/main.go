package main

import (
	"context"
	"log"
	"nailzbydardo/internal/config"
	"nailzbydardo/internal/db"
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

	log.Print("Successfully connected to database")
}