package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestNormalizeEmail(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:  "trim and lowercase",
			input: "  Ritik@Example.COM ",
			want:  "ritik@example.com",
		},
		{
			name:    "empty email",
			input:   "   ",
			wantErr: true,
		},
		{
			name:    "invalid email",
			input:   "not-an-email",
			wantErr: true,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := normalizeEmail(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tc.want {
				t.Fatalf("got %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNewFirebaseIdentityClient(t *testing.T) {
	t.Parallel()

	old := os.Getenv("FIREBASE_WEB_API_KEY")
	t.Cleanup(func() {
		_ = os.Setenv("FIREBASE_WEB_API_KEY", old)
	})

	t.Run("missing api key", func(t *testing.T) {
		t.Parallel()

		_ = os.Unsetenv("FIREBASE_WEB_API_KEY")
		client, err := newFirebaseIdentityClient()
		if err == nil {
			t.Fatalf("expected error, got nil and client=%v", client)
		}
	})

	t.Run("api key present", func(t *testing.T) {
		t.Parallel()

		if err := os.Setenv("FIREBASE_WEB_API_KEY", "test-key"); err != nil {
			t.Fatalf("failed to set env: %v", err)
		}

		client, err := newFirebaseIdentityClient()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if client == nil {
			t.Fatal("expected client, got nil")
		}
		if client.apiKey != "test-key" {
			t.Fatalf("got api key %q, want %q", client.apiKey, "test-key")
		}
		if client.client == nil {
			t.Fatal("expected http client to be initialized")
		}
	})
}

func TestSignupHandlerValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		body           string
		wantStatusCode int
		wantContains   string
	}{
		{
			name:           "reject non post",
			method:         http.MethodGet,
			body:           "",
			wantStatusCode: http.StatusMethodNotAllowed,
			wantContains:   "method not allowed",
		},
		{
			name:           "reject invalid json",
			method:         http.MethodPost,
			body:           "{bad json",
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "invalid request body",
		},
		{
			name:           "reject invalid email",
			method:         http.MethodPost,
			body:           `{"email":"nope","password":"123456"}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "invalid email address",
		},
		{
			name:           "reject short password",
			method:         http.MethodPost,
			body:           `{"email":"user@example.com","password":"123"}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "password must be at least 6 characters",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/auth/signup", strings.NewReader(tc.body))
			rr := httptest.NewRecorder()

			signupHandler(rr, req)

			if rr.Code != tc.wantStatusCode {
				t.Fatalf("got status %d, want %d. body=%s", rr.Code, tc.wantStatusCode, rr.Body.String())
			}
			if !strings.Contains(strings.ToLower(rr.Body.String()), strings.ToLower(tc.wantContains)) {
				t.Fatalf("expected body to contain %q, got %q", tc.wantContains, rr.Body.String())
			}
		})
	}
}

func TestLoginHandlerValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		body           string
		wantStatusCode int
		wantContains   string
	}{
		{
			name:           "reject non post",
			method:         http.MethodGet,
			body:           "",
			wantStatusCode: http.StatusMethodNotAllowed,
			wantContains:   "method not allowed",
		},
		{
			name:           "reject invalid json",
			method:         http.MethodPost,
			body:           "{bad json",
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "invalid request body",
		},
		{
			name:           "reject invalid email",
			method:         http.MethodPost,
			body:           `{"email":"bad-email","password":"123456"}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "invalid email address",
		},
		{
			name:           "reject empty password",
			method:         http.MethodPost,
			body:           `{"email":"user@example.com","password":"   "}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "password is required",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/auth/login", strings.NewReader(tc.body))
			rr := httptest.NewRecorder()

			loginHandler(rr, req)

			if rr.Code != tc.wantStatusCode {
				t.Fatalf("got status %d, want %d. body=%s", rr.Code, tc.wantStatusCode, rr.Body.String())
			}
			if !strings.Contains(strings.ToLower(rr.Body.String()), strings.ToLower(tc.wantContains)) {
				t.Fatalf("expected body to contain %q, got %q", tc.wantContains, rr.Body.String())
			}
		})
	}
}

func TestResendVerificationEmailHandlerValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name           string
		method         string
		body           string
		wantStatusCode int
		wantContains   string
	}{
		{
			name:           "reject non post",
			method:         http.MethodGet,
			body:           "",
			wantStatusCode: http.StatusMethodNotAllowed,
			wantContains:   "method not allowed",
		},
		{
			name:           "reject invalid json",
			method:         http.MethodPost,
			body:           "{bad json",
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "invalid request body",
		},
		{
			name:           "reject invalid email",
			method:         http.MethodPost,
			body:           `{"email":"wrong","password":"123456"}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "invalid email address",
		},
		{
			name:           "reject empty password",
			method:         http.MethodPost,
			body:           `{"email":"user@example.com","password":"   "}`,
			wantStatusCode: http.StatusBadRequest,
			wantContains:   "password is required",
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tc.method, "/auth/verify-email/resend", strings.NewReader(tc.body))
			rr := httptest.NewRecorder()

			resendVerificationEmailHandler(rr, req)

			if rr.Code != tc.wantStatusCode {
				t.Fatalf("got status %d, want %d. body=%s", rr.Code, tc.wantStatusCode, rr.Body.String())
			}
			if !strings.Contains(strings.ToLower(rr.Body.String()), strings.ToLower(tc.wantContains)) {
				t.Fatalf("expected body to contain %q, got %q", tc.wantContains, rr.Body.String())
			}
		})
	}
}

func TestRegisterFirebaseEmailPasswordRoutes(t *testing.T) {
	t.Parallel()

	mux := http.NewServeMux()
	registerFirebaseEmailPasswordRoutes(mux)

	routes := []struct {
		path   string
		method string
	}{
		{path: "/auth/signup", method: http.MethodGet},
		{path: "/auth/login", method: http.MethodGet},
		{path: "/auth/verify-email/resend", method: http.MethodGet},
	}

	for _, route := range routes {
		req := httptest.NewRequest(route.method, route.path, nil)
		rr := httptest.NewRecorder()

		mux.ServeHTTP(rr, req)

		if rr.Code == http.StatusNotFound {
			t.Fatalf("route %s was not registered", route.path)
		}
	}
}

func TestValidationResponsesAreJSON(t *testing.T) {
	t.Parallel()

	req := httptest.NewRequest(http.MethodPost, "/auth/signup", strings.NewReader(`{"email":"bad","password":"123456"}`))
	rr := httptest.NewRecorder()

	signupHandler(rr, req)

	var payload map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &payload); err != nil {
		t.Fatalf("expected JSON response, got error: %v, body=%s", err, rr.Body.String())
	}
}
