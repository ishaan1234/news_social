package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/ishaan1234/news_social/backend/internal/config"
)

// Mock config
func mockConfig() *config.Config {
	return &config.Config{
		Port:          "8080",
		JWTSecret:     "test-secret",
		NewsAPIKey:    "test",
		OpenAIAPIKey:  "test",
		RateLimitRPS:  100,
		DBUrl:         "postgres://test",
	}
}

// Basic server start test
func TestServerRoutes(t *testing.T) {
	cfg := mockConfig()

	// Pass nil DB for now if your constructor allows it,
	// otherwise create a mock Postgres struct
	srv := NewHTTPServer(cfg, nil)

	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	// Test public endpoint
	resp, err := http.Get(ts.URL + "/api/headlines")
	if err != nil {
		t.Fatalf("failed to call endpoint: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected 200, got %d", resp.StatusCode)
	}
}

// Test protected route without token
func TestProtectedRouteUnauthorized(t *testing.T) {
	cfg := mockConfig()
	srv := NewHTTPServer(cfg, nil)

	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	req, _ := http.NewRequest("POST", ts.URL+"/api/summaries", nil)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}
}