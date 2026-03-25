// package main

// import (
// 	"io"
// 	"log"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"time"
// )

// func newsHandler(w http.ResponseWriter, r *http.Request) {
// 	apiKey := os.Getenv("NEWSAPI_KEY")
// 	if apiKey == "" {
// 		writeJSONError(w, http.StatusInternalServerError, "missing NEWSAPI_KEY environment variable")
// 		return
// 	}

// 	// Optional query param: /news?q=tesla
// 	query := r.URL.Query().Get("q")
// 	if query == "" {
// 		query = "tesla"
// 	}

// 	newsURL := "https://newsapi.org/v2/everything?q=" +
// 		url.QueryEscape(query) +
// 		"&sortBy=publishedAt&pageSize=2&language=en"

// 	req, err := http.NewRequest(http.MethodGet, newsURL, nil)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to create NewsAPI request")
// 		return
// 	}

// 	req.Header.Set("X-Api-Key", apiKey)

// 	client := &http.Client{Timeout: 8 * time.Second}

// 	res, err := client.Do(req)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to call NewsAPI")
// 		return
// 	}
// 	defer res.Body.Close()

// 	contentType := res.Header.Get("Content-Type")
// 	if contentType == "" {
// 		contentType = "application/json"
// 	}

// 	w.Header().Set("Content-Type", contentType)
// 	w.WriteHeader(res.StatusCode)

// 	if _, err := io.Copy(w, res.Body); err != nil {
// 		log.Printf("failed to copy NewsAPI response: %v", err)
// 	}
// }

// package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strings"
// 	"time"
// )

// type NewsAPIResponse struct {
// 	Status       string    `json:"status"`
// 	TotalResults int       `json:"totalResults"`
// 	Articles     []Article `json:"articles"`
// }

// type Article struct {
// 	Source      Source `json:"source"`
// 	Author      string `json:"author"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	URL         string `json:"url"`
// 	URLToImage  string `json:"urlToImage"`
// 	PublishedAt string `json:"publishedAt"`
// 	Content     string `json:"content"`
// 	Summary     string `json:"summary,omitempty"`
// }

// type Source struct {
// 	ID   any    `json:"id"`
// 	Name string `json:"name"`
// }

// func newsHandler(w http.ResponseWriter, r *http.Request) {
// 	apiKey := os.Getenv("NEWSAPI_KEY")
// 	if apiKey == "" {
// 		writeJSONError(w, http.StatusInternalServerError, "missing NEWSAPI_KEY environment variable")
// 		return
// 	}

// 	query := r.URL.Query().Get("q")
// 	if query == "" {
// 		query = "tesla"
// 	}

// 	newsURL := "https://newsapi.org/v2/everything?q=" +
// 		url.QueryEscape(query) +
// 		"&sortBy=publishedAt&pageSize=5&language=en"

// 	req, err := http.NewRequest(http.MethodGet, newsURL, nil)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to create NewsAPI request")
// 		return
// 	}

// 	req.Header.Set("X-Api-Key", apiKey)

// 	client := &http.Client{Timeout: 8 * time.Second}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to call NewsAPI")
// 		return
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		writeJSONError(w, http.StatusBadGateway, "NewsAPI returned non-200 response")
// 		return
// 	}

// 	var newsResp NewsAPIResponse
// 	if err := json.NewDecoder(res.Body).Decode(&newsResp); err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to decode NewsAPI response")
// 		return
// 	}

// 	if len(newsResp.Articles) == 0 {
// 		_ = json.NewEncoder(w).Encode(map[string]any{
// 			"success": true,
// 			"summary": "No articles found.",
// 		})
// 		return
// 	}

// 	var parts []string
// 	for _, article := range newsResp.Articles {
// 		part := "Title: " + article.Title + "\n" +
// 			"Description: " + article.Description + "\n" +
// 			"Content: " + article.Content
// 		parts = append(parts, part)
// 	}

// 	combinedText := strings.Join(parts, "\n\n---\n\n")

// 	summary, err := summarizeWithGroq(combinedText)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to summarize news: "+err.Error())
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(map[string]any{
// 		"success": true,
// 		"query":   query,
// 		"summary": summary,
// 	})
// }

// package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"time"
// )

// type NewsAPIResponse struct {
// 	Status       string    `json:"status"`
// 	TotalResults int       `json:"totalResults"`
// 	Articles     []Article `json:"articles"`
// }

// type Article struct {
// 	Source      Source `json:"source"`
// 	Author      string `json:"author"`
// 	Title       string `json:"title"`
// 	Description string `json:"description"`
// 	URL         string `json:"url"`
// 	URLToImage  string `json:"urlToImage"`
// 	PublishedAt string `json:"publishedAt"`
// 	Content     string `json:"content"`
// 	Summary     string `json:"summary,omitempty"`
// }

// type Source struct {
// 	ID   any    `json:"id"`
// 	Name string `json:"name"`
// }

// func newsHandler(w http.ResponseWriter, r *http.Request) {
// 	apiKey := os.Getenv("NEWSAPI_KEY")
// 	if apiKey == "" {
// 		writeJSONError(w, http.StatusInternalServerError, "missing NEWSAPI_KEY environment variable")
// 		return
// 	}

// 	query := r.URL.Query().Get("q")
// 	if query == "" {
// 		query = "tesla"
// 	}

// 	newsURL := "https://newsapi.org/v2/everything?q=" +
// 		url.QueryEscape(query) +
// 		"&sortBy=publishedAt&pageSize=2&language=en"

// 	req, err := http.NewRequest(http.MethodGet, newsURL, nil)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to create NewsAPI request")
// 		return
// 	}

// 	req.Header.Set("X-Api-Key", apiKey)

// 	client := &http.Client{Timeout: 8 * time.Second}
// 	res, err := client.Do(req)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to call NewsAPI")
// 		return
// 	}
// 	defer res.Body.Close()

// 	if res.StatusCode != http.StatusOK {
// 		writeJSONError(w, http.StatusBadGateway, "NewsAPI returned non-200 response")
// 		return
// 	}

// 	var newsResp NewsAPIResponse
// 	if err := json.NewDecoder(res.Body).Decode(&newsResp); err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to decode NewsAPI response")
// 		return
// 	}

// 	for i := range newsResp.Articles {
// 		articleText := "Title: " + newsResp.Articles[i].Title + "\n" +
// 			"Description: " + newsResp.Articles[i].Description + "\n" +
// 			"Content: " + newsResp.Articles[i].Content

// 		summary, err := summarizeWithGroq(articleText)
// 		if err != nil {
// 			newsResp.Articles[i].Summary = "summary unavailable"
// 			continue
// 		}

// 		newsResp.Articles[i].Summary = summary
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(newsResp)
// }

package main

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type NewsAPIResponse struct {
	Status       string    `json:"status"`
	TotalResults int       `json:"totalResults"`
	Articles     []Article `json:"articles"`
}

type Article struct {
	Source      Source `json:"source"`
	Author      string `json:"author"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	URLToImage  string `json:"urlToImage"`
	PublishedAt string `json:"publishedAt"`
	Content     string `json:"content"`
	Summary     string `json:"summary,omitempty"`
}

type Source struct {
	ID   any    `json:"id"`
	Name string `json:"name"`
}

// Handles /news endpoint
func newsHandler(w http.ResponseWriter, r *http.Request) {
	// Get API key from environment
	apiKey := os.Getenv("NEWSAPI_KEY")
	if apiKey == "" {
		writeJSONError(w, http.StatusInternalServerError, "missing NEWSAPI_KEY environment variable")
		return
	}

	// Read query param (?q=...)
	query := r.URL.Query().Get("q")
	if query == "" {
		query = "tesla"
	}

	// Build NewsAPI request URL
	newsURL := "https://newsapi.org/v2/everything?q=" +
		url.QueryEscape(query) +
		"&sortBy=publishedAt&pageSize=2&language=en"

	req, err := http.NewRequest(http.MethodGet, newsURL, nil)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to create NewsAPI request")
		return
	}

	req.Header.Set("X-Api-Key", apiKey)

	// Call NewsAPI
	client := &http.Client{Timeout: 8 * time.Second}
	res, err := client.Do(req)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to call NewsAPI")
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		writeJSONError(w, http.StatusBadGateway, "NewsAPI returned non-200 response")
		return
	}

	var newsResp NewsAPIResponse
	if err := json.NewDecoder(res.Body).Decode(&newsResp); err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to decode NewsAPI response")
		return
	}

	// For each article:
	// 1. Try scraping full article text
	// 2. Fallback to NewsAPI snippet if scraping fails
	// 3. Send text to Groq for summarization
	for i := range newsResp.Articles {
		fullText, err := extractArticleText(newsResp.Articles[i].URL)
		if err != nil || strings.TrimSpace(fullText) == "" {
			fullText = "Title: " + newsResp.Articles[i].Title + "\n" +
				"Description: " + newsResp.Articles[i].Description + "\n" +
				"Content: " + newsResp.Articles[i].Content
		}

		// Send the article text (scraped or fallback) to Groq for summarization
		summary, err := summarizeWithGroq(fullText)

		if err != nil {
			// If summarization fails, store a fallback message
			newsResp.Articles[i].Summary = "summary unavailable"
			continue
		}
		// Save the generated summary into the article
		newsResp.Articles[i].Summary = summary
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(newsResp)
}
