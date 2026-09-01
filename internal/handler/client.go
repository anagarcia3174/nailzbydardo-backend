package handler

import (
	"encoding/json"
	"nailzbydardo/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ClientHandler struct {
	clientService *service.ClientService
}

func NewClientHandler(clientService *service.ClientService) *ClientHandler {
	return &ClientHandler{clientService: clientService}
}

type clientRequest struct {
	ClientName    string     `json:"client_name"`
	ContactMethod *string    `json:"contact_method"`
	Notes         *string    `json:"notes"`
	Birthday      *time.Time `json:"birthday"`
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	client, err := h.clientService.CreateClient(r.Context(), req.ClientName, req.ContactMethod, req.Notes, req.Birthday)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, client)
}

func (h *ClientHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	client, err := h.clientService.GetClient(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) List(w http.ResponseWriter, r *http.Request) {
	client, err := h.clientService.ListClients(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) GetAppointments(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	appointments, err := h.clientService.GetClientAppointments(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, appointments)
}
func (h *ClientHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req clientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	client, err := h.clientService.UpdateClient(r.Context(), id, req.ClientName, req.ContactMethod, req.Notes, req.Birthday)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.clientService.SoftDeleteClient(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *ClientHandler) GetTotalSpent(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	spent, err := h.clientService.GetClientTotalSpent(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, spent)
}
