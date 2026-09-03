package handler

import (
	"encoding/json"
	"errors"
	"log"
	"nailzbydardo/internal/service"
	"net/http"
)

type AuthHandler struct {
	authService   *service.AuthService
	cookieName    string
	secureCookies bool
}

func NewAuthHandler(authServicee *service.AuthService, cookieName string, secureCookies bool) *AuthHandler {
	return &AuthHandler{authService: authServicee, cookieName: cookieName, secureCookies: secureCookies}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{Error: "invalid request body"})
		return
	}
	session, err := h.authService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		writeError(w, err)
		return
	}
	cookie := http.Cookie{
    Name:     h.cookieName,
    Value:    session.ID,
    Expires:  session.ExpiresAt,
    HttpOnly: true,
    SameSite: http.SameSiteLaxMode,
    Path:     "/",
    Secure:   h.secureCookies,
}
	http.SetCookie(w, &cookie)
	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)
	if err == nil {
		err = h.authService.Logout(r.Context(), cookie.Value)
		if err != nil {
			log.Printf("error logging out: %v", err)
		}
	} else if !errors.Is(err, http.ErrNoCookie) {
		http.Error(w, "failed to read cookie", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:   h.cookieName,
		Value:  "",
		MaxAge: -1,
	})

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie(h.cookieName)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	_, err = h.authService.ValidateSession(r.Context(), cookie.Value)
	if err != nil {
		writeError(w, err)
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"authenticated": true})
}
