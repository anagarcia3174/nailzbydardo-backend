package main

import (
	"context"
	"fmt"
	"log"
	"nailzbydardo/internal/config"
	"nailzbydardo/internal/db"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/handler"
	"nailzbydardo/internal/middleware"
	"nailzbydardo/internal/repository"
	"nailzbydardo/internal/router"
	"nailzbydardo/internal/service"
	"net/http"
	"time"
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

	q := sqlc.New(pool)
	clientRepo := repository.NewClientRepository(q)
	clientService := service.NewClientService(clientRepo)
	clientHandler := handler.NewClientHandler(clientService)
	serviceRepo := repository.NewServiceRepository(q)
	catalogService := service.NewCatalogService(serviceRepo)
	catalogHandler := handler.NewCatalogHandler(catalogService)
	healthHandler := handler.NewHealthHandler(pool)
	expenseRepo := repository.NewExpenseRepository(q)
	expenseService := service.NewExpenseService(expenseRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	appointmentRepo := repository.NewAppointmentRepository(q)
	appointmentService := service.NewAppointmentService(appointmentRepo, clientRepo)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	userRepo := repository.NewUserRepository(q)
	sessionRepo := repository.NewSessionRepository(q)
	authService := service.NewAuthService(userRepo, sessionRepo)
	authHandler := handler.NewAuthHandler(authService, cfg.SessionCookieName, cfg.IsProduction())
	authMiddleware := middleware.RequireAuth(authService, cfg.SessionCookieName)
	dashboardService := service.NewDashboardService(appointmentService, expenseService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	handlers := router.Handlers{Client: clientHandler, Health: healthHandler, Catalog: catalogHandler, Expense: expenseHandler, Appointment: appointmentHandler, Auth: authHandler, Dashboard: dashboardHandler}
	mux := router.New(handlers,  authMiddleware)

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
