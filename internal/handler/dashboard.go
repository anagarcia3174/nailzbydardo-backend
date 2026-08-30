package handler

import (
	"nailzbydardo/internal/service"
	"net/http"
	"time"
)

type DashboardHandler struct {
	dashboardService *service.DashboardService
}

func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{dashboardService: dashboardService}
}

func (h *DashboardHandler) Get(w http.ResponseWriter, r *http.Request) {
	dashboardSummary, err := h.dashboardService.GetDashboard(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, dashboardSummary)
}

func (h *DashboardHandler) GetFinancials(w http.ResponseWriter, r *http.Request) {
	start := r.URL.Query().Get("start")
	end := r.URL.Query().Get("end")
	startDate, err := time.Parse(time.RFC3339, start)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid start date"})
		return
	}
	endDate, err := time.Parse(time.RFC3339, end)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid end date"})
		return
	}
	financials, err := h.dashboardService.GetFinancials(r.Context(), startDate, endDate)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, financials)
}
