package handler

import (
	"encoding/json"
	"nailzbydardo/internal/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type ExpenseHandler struct {
	expenseService *service.ExpenseService
}

func NewExpenseHandler(expenseService *service.ExpenseService) *ExpenseHandler {
	return &ExpenseHandler{expenseService: expenseService}
}

type expenseRequest struct {
	ExpenseName   string    `json:"expense_name"`
	Price         int64     `json:"price"`
	DatePurchased time.Time `json:"date_purchased"`
	ReceiptURL    *string   `json:"receipt_url"`
}

func (h *ExpenseHandler) List(w http.ResponseWriter, r *http.Request) {
	expenses, err := h.expenseService.ListExpenses(r.Context())
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, expenses)
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req expenseRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	expense, err := h.expenseService.CreateExpense(r.Context(), req.ExpenseName, req.Price, req.DatePurchased, req.ReceiptURL)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusCreated, expense)
}

func (h *ExpenseHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	err := h.expenseService.DeleteExpense(r.Context(), id)
	if err != nil {
		writeError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
