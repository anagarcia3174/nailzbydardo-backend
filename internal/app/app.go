package app

import (
	"nailzbydardo/internal/config"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/handler"
	"nailzbydardo/internal/middleware"
	"nailzbydardo/internal/repository"
	"nailzbydardo/internal/router"
	"nailzbydardo/internal/service"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func BuildHandlers(pool *pgxpool.Pool, cfg *config.Config) *chi.Mux {
	q := sqlc.New(pool)
	clientRepo := repository.NewClientRepository(q)
	appointmentRepo := repository.NewAppointmentRepository(q)
	appointmentService := service.NewAppointmentService(appointmentRepo, clientRepo)
	appointmentHandler := handler.NewAppointmentHandler(appointmentService)
	clientService := service.NewClientService(clientRepo, appointmentService)
	clientHandler := handler.NewClientHandler(clientService)
	serviceRepo := repository.NewServiceRepository(q)
	catalogService := service.NewCatalogService(serviceRepo)
	catalogHandler := handler.NewCatalogHandler(catalogService)
	healthHandler := handler.NewHealthHandler(pool)
	expenseRepo := repository.NewExpenseRepository(q)
	expenseService := service.NewExpenseService(expenseRepo)
	expenseHandler := handler.NewExpenseHandler(expenseService)
	userRepo := repository.NewUserRepository(q)
	sessionRepo := repository.NewSessionRepository(q)
	authService := service.NewAuthService(userRepo, sessionRepo)
	authHandler := handler.NewAuthHandler(authService, cfg.SessionCookieName, cfg.IsProduction())
	authMiddleware := middleware.RequireAuth(authService, cfg.SessionCookieName)
	dashboardService := service.NewDashboardService(appointmentService, expenseService)
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	calendarHandler := handler.NewCalendarHandler(appointmentService, cfg.CalendarSecret, cfg.BaseURL)
	handlers := router.Handlers{Client: clientHandler, Health: healthHandler, Catalog: catalogHandler, Expense: expenseHandler, Appointment: appointmentHandler, Auth: authHandler, Dashboard: dashboardHandler, Calendar: calendarHandler}
	return router.New(handlers, authMiddleware, cfg.FrontendURL)
}
