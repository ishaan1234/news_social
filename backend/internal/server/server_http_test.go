package server

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/config"
)

func mockConfig() *config.Config {
	return &config.Config{
		Port:         "8080",
		JWTSecret:    "test-secret",
		NewsAPIKey:   "test",
		OpenAIAPIKey: "test",
		RateLimitRPS: 100,
		DBUrl:        "postgres://test",
	}
}

func TestProtectedRouteUnauthorized(t *testing.T) {
	srv := NewHTTPServer(mockConfig(), nil)
	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/api/summaries", "application/json", nil)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestPublicRouteValidatesInput(t *testing.T) {
	srv := NewHTTPServer(mockConfig(), nil)
	ts := httptest.NewServer(srv.mux)
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/api/articles")
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}
