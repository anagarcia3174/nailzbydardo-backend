package handler

import (
	"encoding/json"
	"nailzbydardo/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type CatalogHandler struct {
	catalogService *service.CatalogService
}

func NewCatalogHandler(catalogService *service.CatalogService) *CatalogHandler {
	return &CatalogHandler{catalogService: catalogService}
}

type catalogRequest struct {
	ServiceName  string `json:"service_name"`
	Price int64  `json:"service_price"`
}

func (h *CatalogHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req catalogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	svc, err := h.catalogService.CreateService(r.Context(), req.ServiceName, req.Price)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, svc)
}

func (h *CatalogHandler) List(w http.ResponseWriter, r *http.Request) {
	svc, err := h.catalogService.ListServices(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *CatalogHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req catalogRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	svc, err := h.catalogService.UpdateService(r.Context(), id, req.ServiceName, req.Price)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, svc)
}

func (h *CatalogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.catalogService.DeleteService(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
