package handler

import (
	"encoding/json"
	"nailzbydardo/internal/db/sqlc"
	"nailzbydardo/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type AppointmentHandler struct {
	appointmentService *service.AppointmentService
}

func NewAppointmentHandler(appointmentService *service.AppointmentService) *AppointmentHandler {
	return &AppointmentHandler{appointmentService: appointmentService}
}

type createAppointmentRequest struct {
	ClientID string    `json:"client_id"`
	ApptDate time.Time `json:"appt_date"`
	Notes    *string   `json:"notes"`
}
type updateAppointmentRequest struct {
	ApptDate      time.Time              `json:"appt_date"`
	ApptStatus    sqlc.AppointmentStatus `json:"appt_status"`
	LateFee       *int64                 `json:"late_fee"`
	PaymentMethod *sqlc.PaymentMethod    `json:"payment_method"`
	Notes         *string                `json:"notes"`
	ReceiptURL    *string                `json:"receipt_url"`
	LoyaltyReward bool                   `json:"loyalty_reward"`
	Tip           *int64                 `json:"tip"`
}
type addAppointmentServiceRequest struct {
	ServiceName  string `json:"service_name"`
	ServicePrice int64  `json:"service_price"`
	DesignPrice  int64  `json:"design_price"`
}
type addAppointmentDiscountRequest struct {
	DiscountName  string            `json:"discount_name"`
	DiscountType  sqlc.DiscountType `json:"discount_type"`
	DiscountValue int64             `json:"discount_value"`
}

func (h *AppointmentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req createAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	appointment, err := h.appointmentService.CreateAppointment(r.Context(), req.ClientID, req.ApptDate, req.Notes)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, appointment)
}

func (h *AppointmentHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	appointment, err := h.appointmentService.GetAppointmentDetail(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointment)
}
func (h *AppointmentHandler) List(w http.ResponseWriter, r *http.Request) {
	appointments, err := h.appointmentService.ListAppointments(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointments)
}
func (h *AppointmentHandler) ListUpcoming(w http.ResponseWriter, r *http.Request) {
	appointments, err := h.appointmentService.ListUpcomingAppointments(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointments)
}
func (h *AppointmentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req updateAppointmentRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	appointment, err := h.appointmentService.UpdateAppointment(r.Context(), id, req.ApptDate, req.ApptStatus, req.LateFee, req.PaymentMethod, req.Notes, req.ReceiptURL, req.LoyaltyReward, req.Tip)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointment)
}

func (h *AppointmentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.appointmentService.DeleteAppointment(r.Context(),  id)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *AppointmentHandler) GetTotal(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	total, err := h.appointmentService.CalculateAppointmentTotal(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, total)
}
func (h *AppointmentHandler) AddService(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req addAppointmentServiceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	service, err := h.appointmentService.AddAppointmentService(r.Context(), id, req.ServiceName, req.ServicePrice,  req.DesignPrice)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, service)
}
func (h *AppointmentHandler) RemoveService(w http.ResponseWriter, r *http.Request) {
	serviceID := chi.URLParam(r, "serviceId")

	err := h.appointmentService.RemoveAppointmentService(r.Context(), serviceID)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *AppointmentHandler) AddDiscount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req addAppointmentDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	service, err := h.appointmentService.AddAppointmentDiscount(r.Context(), id, req.DiscountName, req.DiscountType,  req.DiscountValue)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, service)
}
func (h *AppointmentHandler) RemoveDiscount(w http.ResponseWriter, r *http.Request) {
	discountID := chi.URLParam(r, "discountId")

	err := h.appointmentService.RemoveAppointmentDiscount(r.Context(), discountID)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}