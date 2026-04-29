package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func withMockTransport(t *testing.T, fn roundTripFunc) {
	t.Helper()

	original := http.DefaultTransport
	http.DefaultTransport = fn

	t.Cleanup(func() {
		http.DefaultTransport = original
	})
}

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"application/json"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

func htmlResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": []string{"text/html"}},
		Body:       io.NopCloser(bytes.NewBufferString(body)),
	}
}

// 1) writeJSONError
func TestWriteJSONError(t *testing.T) {
	rr := httptest.NewRecorder()

	writeJSONError(rr, http.StatusBadRequest, "bad request")

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected status %d, got %d", http.StatusBadRequest, rr.Code)
	}

	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("expected content-type application/json, got %s", ct)
	}

	var got map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode response JSON: %v", err)
	}

	if got["success"] != false {
		t.Fatalf("expected success=false, got %v", got["success"])
	}

	if got["error"] != "bad request" {
		t.Fatalf("expected error 'bad request', got %v", got["error"])
	}
}

// 2) parseDotEnv
func TestParseDotEnv(t *testing.T) {
	dir := t.TempDir()
	envPath := filepath.Join(dir, ".env")

	content := `
# comment
NEWSAPI_KEY=test-news-key
GROQ_API_KEY="test-groq-key"
EMPTY_VALUE=
`
	if err := os.WriteFile(envPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write temp env file: %v", err)
	}

	_ = os.Unsetenv("NEWSAPI_KEY")
	_ = os.Unsetenv("GROQ_API_KEY")
	_ = os.Unsetenv("EMPTY_VALUE")

	if err := parseDotEnv(envPath); err != nil {
		t.Fatalf("parseDotEnv returned error: %v", err)
	}

	if got := os.Getenv("NEWSAPI_KEY"); got != "test-news-key" {
		t.Fatalf("expected NEWSAPI_KEY=test-news-key, got %q", got)
	}

	if got := os.Getenv("GROQ_API_KEY"); got != "test-groq-key" {
		t.Fatalf("expected GROQ_API_KEY=test-groq-key, got %q", got)
	}

	if got := os.Getenv("EMPTY_VALUE"); got != "" {
		t.Fatalf("expected EMPTY_VALUE='', got %q", got)
	}
}

// 3) loadDotEnv
func TestLoadDotEnv(t *testing.T) {
	dir := t.TempDir()

	originalWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalWD)
	}()

	if err := os.Chdir(dir); err != nil {
		t.Fatalf("failed to chdir to temp dir: %v", err)
	}

	envContent := "NEWSAPI_KEY=loaded-from-loadDotEnv\n"
	if err := os.WriteFile(".env", []byte(envContent), 0o600); err != nil {
		t.Fatalf("failed to write .env: %v", err)
	}

	_ = os.Unsetenv("NEWSAPI_KEY")

	loadDotEnv()

	if got := os.Getenv("NEWSAPI_KEY"); got != "loaded-from-loadDotEnv" {
		t.Fatalf("expected NEWSAPI_KEY=loaded-from-loadDotEnv, got %q", got)
	}
}

// 4) summarizeWithGroq
func TestSummarizeWithGroq(t *testing.T) {
	t.Setenv("GROQ_API_KEY", "test-groq-key")

	withMockTransport(t, func(req *http.Request) (*http.Response, error) {
		if req.URL.String() != "https://api.groq.com/openai/v1/chat/completions" {
			t.Fatalf("unexpected URL: %s", req.URL.String())
		}

		if got := req.Header.Get("Authorization"); got != "Bearer test-groq-key" {
			t.Fatalf("unexpected Authorization header: %s", got)
		}

		body, err := io.ReadAll(req.Body)
		if err != nil {
			t.Fatalf("failed to read request body: %v", err)
		}

		if !strings.Contains(string(body), "Tesla expands in Canada") {
			t.Fatalf("expected request body to contain article text, got: %s", string(body))
		}

		mockGroqJSON := `{
			"choices": [
				{
					"message": {
						"content": "Tesla expands in Canada amid stronger EV competition."
					}
				}
			]
		}`

		return jsonResponse(http.StatusOK, mockGroqJSON), nil
	})

	summary, err := summarizeWithGroq("Tesla expands in Canada")
	if err != nil {
		t.Fatalf("summarizeWithGroq returned error: %v", err)
	}

	expected := "Tesla expands in Canada amid stronger EV competition."
	if summary != expected {
		t.Fatalf("expected %q, got %q", expected, summary)
	}
}

// 5) extractArticleText
func TestExtractArticleText(t *testing.T) {
	withMockTransport(t, func(req *http.Request) (*http.Response, error) {
		return htmlResponse(http.StatusForbidden, "<html><body>forbidden</body></html>"), nil
	})

	_, err := extractArticleText("https://example.com/article")
	if err == nil {
		t.Fatal("expected extractArticleText to return an error for non-200 response")
	}

	if !strings.Contains(err.Error(), "article page returned status 403") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// 6) newsHandler
func TestNewsHandler(t *testing.T) {
	t.Setenv("NEWSAPI_KEY", "test-news-key")
	t.Setenv("GROQ_API_KEY", "test-groq-key")

	withMockTransport(t, func(req *http.Request) (*http.Response, error) {
		switch {
		case strings.Contains(req.URL.Host, "newsapi.org"):
			mockNewsAPIJSON := `{
				"status": "ok",
				"totalResults": 1,
				"articles": [
					{
						"source": { "id": null, "name": "Electrek" },
						"author": "Reporter",
						"title": "Tesla expands in Canada",
						"description": "Expansion news",
						"url": "https://example.com/article",
						"urlToImage": "https://example.com/image.jpg",
						"publishedAt": "2026-03-25T10:00:00Z",
						"content": "Snippet from NewsAPI"
					}
				]
			}`
			return jsonResponse(http.StatusOK, mockNewsAPIJSON), nil

		case req.URL.String() == "https://example.com/article":
			html := `
			<!doctype html>
			<html>
			<head><title>Tesla Article</title></head>
			<body>
				<article>
					<p>Tesla is opening new locations in Canada and expanding its EV presence.</p>
					<p>The move reflects stronger market momentum.</p>
				</article>
			</body>
			</html>`
			return htmlResponse(http.StatusOK, html), nil

		case strings.Contains(req.URL.Host, "api.groq.com"):
			mockGroqJSON := `{
				"choices": [
					{
						"message": {
							"content": "Tesla is expanding its Canadian presence through new locations and broader EV market growth."
						}
					}
				]
			}`
			return jsonResponse(http.StatusOK, mockGroqJSON), nil

		default:
			t.Fatalf("unexpected request URL: %s", req.URL.String())
			return nil, nil
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/news?q=tesla", nil)
	rr := httptest.NewRecorder()

	newsHandler(nil)(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, rr.Code)
	}

	var got NewsAPIResponse
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("failed to decode handler response: %v", err)
	}

	if len(got.Articles) != 1 {
		t.Fatalf("expected 1 article, got %d", len(got.Articles))
	}

	if got.Articles[0].Summary == "" {
		t.Fatal("expected article summary to be populated")
	}

	expectedSummary := "Tesla is expanding its Canadian presence through new locations and broader EV market growth."
	if got.Articles[0].Summary != expectedSummary {
		t.Fatalf("expected summary %q, got %q", expectedSummary, got.Articles[0].Summary)
	}
}
