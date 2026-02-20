package main

import (
	"bufio"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// main is the entry point of the application.
// It loads environment variables, registers HTTP routes,
// main loads environment variables from .env (if present), registers the /news HTTP handler, and starts the HTTP server on :8080.
// It logs the listening address and exits if the server fails to start.
func main() {
	// Load environment variables from .env file (if present)
	loadDotEnv()

	// Create a new HTTP request multiplexer (router)
	mux := http.NewServeMux()

	// Register /news endpoint and bind it to newsHandler
	mux.HandleFunc("/news", newsHandler)

	log.Println("server listening on http://localhost:8080")

	// Start HTTP server on port 8080
	// If the server fails, log.Fatal will exit the program
	log.Fatal(http.ListenAndServe(":8080", mux))
}

// newsHandler handles GET requests to /news.
// newsHandler handles HTTP requests to /news by proxying NewsAPI results for the query "tesla" to the client.
// It requires the NEWSAPI_KEY environment variable; on success it mirrors NewsAPI's HTTP status code and Content-Type
// (defaults to "application/json" if absent) and streams the response body directly to the client. On failure it sends a
// JSON-formatted error response.
func newsHandler(w http.ResponseWriter, _ *http.Request) {

	// Read NewsAPI key from environment variables
	apiKey := os.Getenv("NEWSAPI_KEY")
	if apiKey == "" {
		writeJSONError(w, http.StatusInternalServerError, "missing NEWSAPI_KEY environment variable")
		return
	}

	// Create a new HTTP request to NewsAPI
	req, err := http.NewRequest(http.MethodGet, "https://newsapi.org/v2/everything?q=tesla", nil)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to create NewsAPI request")
		return
	}

	// Set API key in request header
	req.Header.Set("X-Api-Key", apiKey)

	// Create HTTP client with timeout to avoid hanging requests
	client := &http.Client{Timeout: 8 * time.Second}

	// Execute request to NewsAPI
	res, err := client.Do(req)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to call NewsAPI")
		return
	}
	defer res.Body.Close()

	// Forward response content-type header
	contentType := res.Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/json"
	}

	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(res.StatusCode)

	// Stream NewsAPI response body directly to client
	if _, err := io.Copy(w, res.Body); err != nil {
		log.Printf("failed to copy NewsAPI response: %v", err)
	}
}

// writeJSONError writes an HTTP response with the provided status and a JSON
// body indicating failure.
//
// The response Content-Type is set to "application/json". The JSON body is an
// object with "success": false and "error": <message>.
func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	// Encode error response as JSON
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": false,
		"error":   message,
	})
}

// loadDotEnv attempts to load environment variables
// loadDotEnv attempts to load environment variables from common .env file locations.
// It looks for backend/.env first and then .env, stops after the first file that is successfully parsed, and leaves existing environment variables unchanged when parsing.
func loadDotEnv() {
	paths := []string{"backend/.env", ".env"}

	for _, p := range paths {
		if err := parseDotEnv(p); err == nil {
			return
		}
	}
}

// parseDotEnv reads a .env file and sets environment variables.
// parseDotEnv reads a dotenv-format file at the given path and sets environment
// variables for keys that are not already defined in the process environment.
// It ignores empty lines and lines beginning with `#`. Lines are parsed using the
// first `=` as the separator into `KEY=VALUE`; keys and values are trimmed of
// surrounding whitespace and values have surrounding single or double quotes
// removed. Lines without a key or without `=` are skipped. Any error from
// opening the file or from scanning its contents is returned.
func parseDotEnv(path string) error {
	f, err := os.Open(filepath.Clean(path))
	if err != nil {
		return err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split key=value format
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		val := strings.Trim(strings.TrimSpace(parts[1]), `"'`)

		if key == "" {
			continue
		}

		// Only set variable if not already defined
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, val)
		}
	}

	return scanner.Err()
}