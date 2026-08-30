package middleware

import (
	"nailzbydardo/internal/service"
	"net/http"
)


func RequireAuth(authService *service.AuthService, cookieName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(cookieName)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			_, err = authService.ValidateSession(r.Context(), cookie.Value)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
