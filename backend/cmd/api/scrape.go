package main

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	readability "github.com/mackee/go-readability"
)

// extractArticleText takes a news article URL,
// fetches the webpage, extracts the main readable content,
// and returns it as plain text.
func extractArticleText(articleURL string) (string, error) {

	// Create an HTTP client with timeout
	// so requests don’t hang forever like a bad conversation
	client := &http.Client{
		Timeout: 12 * time.Second,
	}

	// Build a GET request for the article URL
	req, err := http.NewRequest(http.MethodGet, articleURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create article request: %w", err)
	}

	// Set User-Agent header to mimic a real browser
	// some sites block “unknown” clients otherwise
	req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; NewsSummarizer/1.0)")

	// Send the HTTP request to fetch the webpage
	res, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch article page: %w", err)
	}
	defer res.Body.Close()

	// Ensure we got a valid response (HTTP 200 OK)
	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("article page returned status %d", res.StatusCode)
	}

	// Read the entire HTML body of the webpage
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read article page: %w", err)
	}

	// Initialize readability options (default settings)
	options := readability.DefaultOptions()

	// Extract the main article content from raw HTML
	// this removes ads, navbars, sidebars, etc.
	article, err := readability.Extract(string(body), options)
	if err != nil {
		return "", fmt.Errorf("failed to extract readable content: %w", err)
	}

	// Convert extracted content into plain text
	// (clean, readable article body)
	text := strings.TrimSpace(readability.ExtractTextContent(article.Root))

	// Fallback: if text extraction fails,
	// convert content into markdown instead
	if text == "" && article.Root != nil {
		text = strings.TrimSpace(readability.ToMarkdown(article.Root))
	}

	// If still empty, something went wrong
	if text == "" {
		return "", fmt.Errorf("extracted article text is empty")
	}

	// Return the cleaned article text
	return text, nil
}
