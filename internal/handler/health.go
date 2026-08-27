package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// imports you'll need:
// - "context"
// - "encoding/json" — if you want the response body to be JSON (recommended,
//   since your whole API will be JSON — good to establish the convention now)
// - "net/http"
// - "time" — for a ping timeout, same reasoning as step 3
// - "github.com/jackc/pgx/v5/pgxpool" — to reference *pgxpool.Pool in your struct

type HealthHandler struct {
	pool *pgxpool.Pool
}

type HealthResponse struct {
    Status   string `json:"status"`
    Database string `json:"database"`
}

func NewHealthHandler(pool *pgxpool.Pool) *HealthHandler {
	return &HealthHandler{pool: pool}
}

func (h *HealthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctxWithTimeout, cancel := context.WithTimeout(r.Context(), 10*time.Second)
	defer cancel()

	err := h.pool.Ping(ctxWithTimeout)
	response := HealthResponse{}
	statusCode := http.StatusOK
	if err != nil {
		statusCode = http.StatusServiceUnavailable
		response.Status = "error"
		response.Database = "unreachable"
	} else {
		response.Status = "ok"
		response.Database = "connected"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}