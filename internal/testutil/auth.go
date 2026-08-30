package testutil

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"
)

func CreateTestUser(t *testing.T, pool *pgxpool.Pool, email, password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("error hashing password: %v", err)
	}
	_, err = pool.Exec(t.Context(), "INSERT INTO users (email, password_hash) VALUES ($1, $2)", email, hash)
	if err != nil {
		t.Fatalf("error creating test user: %v", err)
	}
}

func Login(t *testing.T, server *httptest.Server, email, password string) *http.Cookie {
	loginRequest := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{
		Email:    email,
		Password: password,
	}
	body, err := json.Marshal(loginRequest)
	if err != nil {
		t.Fatalf("error creating login request body: %v", err)
	}
	req, err := http.NewRequest(http.MethodPost, server.URL+"/auth/login", bytes.NewReader(body))
	if err != nil {
		t.Fatalf("error creating login request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("error sending login request: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status %d, got %d", http.StatusNoContent, resp.StatusCode)
	}
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session_id" {
			return cookie
		}
	}
	t.Fatalf("session cookie not found")
	return nil
}
