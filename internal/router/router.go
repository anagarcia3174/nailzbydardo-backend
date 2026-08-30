package router

import (
	"nailzbydardo/internal/handler"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	Health      *handler.HealthHandler
	Client      *handler.ClientHandler
	Catalog     *handler.CatalogHandler
	Expense     *handler.ExpenseHandler
	Appointment *handler.AppointmentHandler
	Auth        *handler.AuthHandler
	Dashboard   *handler.DashboardHandler
}

func New(h Handlers, authMiddleware func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/health", h.Health.ServeHTTP)
	r.Post("/auth/login", h.Auth.Login)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)

		r.Post("/auth/logout", h.Auth.Logout)
		r.Get("/auth/me", h.Auth.Me)

		r.Route("/clients", func(r chi.Router) {
			r.Get("/", h.Client.List)
			r.Post("/", h.Client.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Client.Get)
				r.Get("/appointments", h.Client.GetAppointments)
				r.Patch("/", h.Client.Update)
				r.Delete("/", h.Client.Delete)
			})
		})
		r.Route("/services", func(r chi.Router) {
			r.Get("/", h.Catalog.List)
			r.Post("/", h.Catalog.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.Patch("/", h.Catalog.Update)
				r.Delete("/", h.Catalog.Delete)
			})
		})
		r.Route("/expenses", func(r chi.Router) {
			r.Get("/", h.Expense.List)
			r.Post("/", h.Expense.Create)

			r.Route("/{id}", func(r chi.Router) {
				r.Delete("/", h.Expense.Delete)
			})
		})
		r.Route("/appointments", func(r chi.Router) {
			r.Get("/", h.Appointment.List)
			r.Post("/", h.Appointment.Create)
			r.Get("/upcoming", h.Appointment.ListUpcoming)

			r.Route("/{id}", func(r chi.Router) {
				r.Get("/", h.Appointment.Get)
				r.Patch("/", h.Appointment.Update)
				r.Delete("/", h.Appointment.Delete)
				r.Get("/total", h.Appointment.GetTotal)

				r.Post("/services", h.Appointment.AddService)
				r.Delete("/services/{serviceId}", h.Appointment.RemoveService)

				r.Post("/discounts", h.Appointment.AddDiscount)
				r.Delete("/discounts/{discountId}", h.Appointment.RemoveDiscount)
			})
		})
		r.Get("/dashboard", h.Dashboard.Get)
		r.Get("/financials", h.Dashboard.GetFinancials)

	})
	return r
}
