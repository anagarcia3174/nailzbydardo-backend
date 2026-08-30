package handler_test

import (
	"encoding/json"
	"fmt"
	"nailzbydardo/internal/model"
	"nailzbydardo/internal/testutil"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestCreateClient(t *testing.T) {
	server, pool := testutil.NewTestServer(t)
	testutil.CreateTestUser(t, pool, "test@example.com", "testpassword123")
	cookie := testutil.Login(t, server, "test@example.com", "testpassword123")

	tests := []struct {
		name         string
		body         string
		wantStatus   int
		wantName     string
		wantMethod   *string
		wantNotes    *string
		wantBirthday *time.Time
	}{
		{
			name:       "valid client",
			body:       `{"client_name": "Test Client"}`,
			wantStatus: http.StatusCreated,
			wantName:   "Test Client",
		},
		{
			name:         "valid client with all fields",
			body:         `{"client_name": "Test Client", "contact_method": "Instagram", "notes": "allergic to gel", "birthday": "1995-06-15T00:00:00Z"}`,
			wantStatus:   http.StatusCreated,
			wantName:     "Test Client",
			wantMethod:   testutil.StrPtr("Instagram"),
			wantNotes:    testutil.StrPtr("allergic to gel"),
			wantBirthday: testutil.TimePtr(time.Date(1995, 6, 15, 0, 0, 0, 0, time.UTC)),
		},
		{
			name:       "blank name",
			body:       `{"client_name": "   "}`,
			wantStatus: http.StatusBadRequest,
		},
		{
			name:       "malformed json",
			body:       `{"client_name":`,
			wantStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, _ := http.NewRequest("POST", server.URL+"/clients", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(cookie)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusCreated {
				return
			}

			var got model.Client
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if got.ID == "" {
				t.Error("expected a non-empty ID")
			}
			if got.ClientName != tt.wantName {
				t.Errorf("got client_name %q, want %q", got.ClientName, tt.wantName)
			}
			if !testutil.EqualStrPtr(got.ContactMethod, tt.wantMethod) {
				t.Errorf("got contact_method %v, want %v", testutil.DerefStr(got.ContactMethod), testutil.DerefStr(tt.wantMethod))
			}
			if !testutil.EqualStrPtr(got.Notes, tt.wantNotes) {
				t.Errorf("got notes %v, want %v", testutil.DerefStr(got.Notes), testutil.DerefStr(tt.wantNotes))
			}
			if !testutil.EqualTimePtr(got.Birthday, tt.wantBirthday) {
				t.Errorf("got birthday %v, want %v", got.Birthday, tt.wantBirthday)
			}

		})
	}
}
func TestCreateClient_Unauthenticated(t *testing.T) {
	server, _ := testutil.NewTestServer(t)
	req, _ := http.NewRequest("POST", server.URL+"/clients", strings.NewReader(`{"client_name":"Test"}`))
	req.Header.Set("Content-Type", "application/json")
	// no cookie attached
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}
func TestGetClient(t *testing.T) {
	server, pool := testutil.NewTestServer(t)
	testutil.CreateTestUser(t, pool, "test@example.com", "testpassword123")
	cookie := testutil.Login(t, server, "test@example.com", "testpassword123")
	clientID := testutil.CreateTestClient(t, pool, "Jane Doe", testutil.StrPtr("Whatsapp"), testutil.StrPtr("Allergies"), testutil.TimePtr(time.Date(2004, 6, 15, 0, 0, 0, 0, time.UTC)))
	tests := []struct {
		name         string
		id           string
		wantStatus   int
		wantName     string
		wantMethod   *string
		wantNotes    *string
		wantBirthday *time.Time
	}{
		{name: "existing client", id: clientID, wantStatus: http.StatusOK, wantName: "Jane Doe",
			wantMethod:   testutil.StrPtr("Whatsapp"),
			wantNotes:    testutil.StrPtr("Allergies"),
			wantBirthday: testutil.TimePtr(time.Date(2004, 6, 15, 0, 0, 0, 0, time.UTC))},
		{name: "nonexistent client", id: "00000000-0000-0000-0000-000000000000", wantStatus: http.StatusNotFound},
		{name: "invalid id", id: "not-an-id", wantStatus: http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := fmt.Sprintf("/clients/%s", tt.id)
			req, _ := http.NewRequest(http.MethodGet, server.URL+path, nil)
			req.AddCookie(cookie)

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatalf("request failed: %v", err)
			}
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Errorf("got status %d, want %d", resp.StatusCode, tt.wantStatus)
			}

			if tt.wantStatus != http.StatusOK {
				return
			}

			var got model.Client
			if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
				t.Fatalf("failed to decode response: %v", err)
			}

			if got.ID != tt.id {
				t.Errorf("got id %q, want %q", got.ID, tt.id)
			}
			if got.ClientName != tt.wantName {
				t.Errorf("got client_name %q, want %q", got.ClientName, tt.wantName)
			}
			if !testutil.EqualStrPtr(got.ContactMethod, tt.wantMethod) {
				t.Errorf("got contact_method %v, want %v", testutil.DerefStr(got.ContactMethod), testutil.DerefStr(tt.wantMethod))
			}
			if !testutil.EqualStrPtr(got.Notes, tt.wantNotes) {
				t.Errorf("got notes %v, want %v", testutil.DerefStr(got.Notes), testutil.DerefStr(tt.wantNotes))
			}
			if !testutil.EqualTimePtr(got.Birthday, tt.wantBirthday) {
				t.Errorf("got birthday %v, want %v", got.Birthday, tt.wantBirthday)
			}

		})
	}
}
func TestGetClient_Unauthenticated(t *testing.T) {
	server, _ := testutil.NewTestServer(t)
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/clients/00000000-0000-0000-0000-000000000000", nil)
	// no cookie attached
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestListClients(t *testing.T) {
	server, pool := testutil.NewTestServer(t)
	testutil.CreateTestUser(t, pool, "test@example.com", "testpassword123")
	cookie := testutil.Login(t, server, "test@example.com", "testpassword123")
	clientID := testutil.CreateTestClient(t, pool, "Jane Doe", testutil.StrPtr("Whatsapp"), testutil.StrPtr("Allergies"), testutil.TimePtr(time.Date(2004, 6, 15, 0, 0, 0, 0, time.UTC)))
	clientTwoID := testutil.CreateTestClient(t, pool, "John Doe", testutil.StrPtr("Tik Tok"), nil, nil)
	wantList := []model.ClientSummary{
		{ID: clientID, ClientName: "Jane Doe", ContactMethod: testutil.StrPtr("Whatsapp")},
		{ID: clientTwoID, ClientName: "John Doe", ContactMethod: testutil.StrPtr("Tik Tok")},
	}
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/clients", nil)
	req.AddCookie(cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	var got []model.ClientSummary
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) != 2 {
		t.Errorf("got slice length %d, want %d", len(got), 2)
	}
	for i, c := range got {
		wantClient := wantList[i]
		if wantClient.ID != c.ID {
			t.Errorf("%d: got client id %q, want %q", i, c.ID, wantClient.ID)

		}
		if wantClient.ClientName != c.ClientName {
			t.Errorf("%d: got client name %q, want %q", i, c.ClientName, wantClient.ClientName)

		}
		if !testutil.EqualStrPtr(c.ContactMethod, wantClient.ContactMethod) {
			t.Errorf("%d: got contact_method %v, want %v", i, testutil.DerefStr(c.ContactMethod), testutil.DerefStr(wantClient.ContactMethod))
		}
	}
}

func TestListClients_Empty(t *testing.T) {
	server, pool := testutil.NewTestServer(t)
	testutil.CreateTestUser(t, pool, "test@example.com", "testpassword123")
	cookie := testutil.Login(t, server, "test@example.com", "testpassword123")

	req, _ := http.NewRequest(http.MethodGet, server.URL+"/clients", nil)
	req.AddCookie(cookie)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusOK)
	}

	got := []model.ClientSummary{}
	if err := json.NewDecoder(resp.Body).Decode(&got); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(got) > 0 {
		t.Errorf("got slice length %d, want %d", len(got), 0)
	}
}

func TestListClient_Unauthenticated(t *testing.T) {
	server, _ := testutil.NewTestServer(t)
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/clients", nil)
	// no cookie attached
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("got status %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

// ├── Get Appointments
// │   ├── Client with appointments
// │   ├── Client with no appointments
// │   ├── Nonexistent client
// │   ├── Invalid ID
// │   └── Unauthenticated
// │
// ├── Update
// │   ├── Valid update
// │   ├── Update all fields
// │   ├── Optional fields omitted
// │   ├── Blank name
// │   ├── Malformed JSON
// │   ├── Nonexistent client
// │   ├── Invalid ID
// │   └── Unauthenticated
// │
// └── Delete
//     ├── Existing client
//     ├── Nonexistent client
//     ├── Invalid ID
//     ├── Unauthenticated
//     └── Soft-deleted client no longer accessible
