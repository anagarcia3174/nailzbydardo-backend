// internal/handler/errors.go
package handler

import (
	"encoding/json"
	"errors"
	"log"
	"nailzbydardo/internal/repository"
	"nailzbydardo/internal/service"
	"net/http"
)


type errorResponse struct {
	Error string `json:"error"`
}

func writeError(w http.ResponseWriter, err error) {
	var status int
	switch {
	case errors.Is(err, repository.ErrNotFound):
		status = http.StatusNotFound // 404
	case errors.Is(err, service.ErrInvalidInput):
		status = http.StatusBadRequest // 400
	case errors.Is(err, service.ErrNonexistentClient):
		status = http.StatusBadRequest // 400
	case errors.Is(err, service.ErrInvalidCredentials):
		status = http.StatusUnauthorized // 401
	case errors.Is(err, service.ErrUnauthorized):
		status = http.StatusUnauthorized // 401
	case errors.Is(err, repository.ErrInvalidID):
    	status = http.StatusBadRequest // 400
	default:
		status = http.StatusInternalServerError // 500
	}
	message := err.Error()
	if status == http.StatusInternalServerError {
		log.Printf("internal error: %v", err)
		message = "internal server error"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
