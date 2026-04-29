// package main

// import (
// 	"encoding/json"
// 	"net/http"
// 	"net/url"
// 	"os"
// 	"strings"
// 	"time"
// )
// import "database/sql"

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

// func saveArticleToDB(db *sql.DB, a Article) error {
// 	_, err := db.Exec(`
// 		INSERT INTO articles (
// 			title,
// 			description,
// 			content,
// 			summary,
// 			author,
// 			source_name,
// 			source_id,
// 			url,
// 			image_url,
// 			published_at
// 		)
// 		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
// 		ON CONFLICT (url) DO UPDATE SET
// 			title = EXCLUDED.title,
// 			description = EXCLUDED.description,
// 			content = EXCLUDED.content,
// 			summary = EXCLUDED.summary,
// 			author = EXCLUDED.author,
// 			source_name = EXCLUDED.source_name,
// 			source_id = EXCLUDED.source_id,
// 			image_url = EXCLUDED.image_url,
// 			published_at = EXCLUDED.published_at
// 	`,
// 		a.Title,
// 		a.Description,
// 		a.Content,
// 		a.Summary,
// 		a.Author,
// 		a.Source.Name,
// 		a.Source.ID,
// 		a.URL,
// 		a.URLToImage,
// 		a.PublishedAt,
// 	)

// 	return err
// }

// // Handles /news endpoint
// func newsHandler(w http.ResponseWriter, r *http.Request) {
// 	// Get API key from environment
// 	apiKey := os.Getenv("NEWSAPI_KEY")
// 	if apiKey == "" {
// 		writeJSONError(w, http.StatusInternalServerError, "missing NEWSAPI_KEY environment variable")
// 		return
// 	}

// 	// Read query param (?q=...)
// 	query := r.URL.Query().Get("q")
// 	if query == "" {
// 		query = "tesla"
// 	}

// 	// Build NewsAPI request URL
// 	newsURL := "https://newsapi.org/v2/everything?q=" +
// 		url.QueryEscape(query) +
// 		"&sortBy=publishedAt&pageSize=10&language=en"

// 	req, err := http.NewRequest(http.MethodGet, newsURL, nil)
// 	if err != nil {
// 		writeJSONError(w, http.StatusBadGateway, "failed to create NewsAPI request")
// 		return
// 	}

// 	req.Header.Set("X-Api-Key", apiKey)

// 	// Call NewsAPI
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

// 	// For each article:
// 	// 1. Try scraping full article text
// 	// 2. Fallback to NewsAPI snippet if scraping fails
// 	// 3. Send text to Groq for summarization
// 	for i := range newsResp.Articles {
// 		fullText, err := extractArticleText(newsResp.Articles[i].URL)
// 		if err != nil || strings.TrimSpace(fullText) == "" {
// 			fullText = "Title: " + newsResp.Articles[i].Title + "\n" +
// 				"Description: " + newsResp.Articles[i].Description + "\n" +
// 				"Content: " + newsResp.Articles[i].Content
// 		}

// 		// Send the article text (scraped or fallback) to Groq for summarization
// 		summary, err := summarizeWithGroq(fullText)

// 		if err != nil {
// 			// If summarization fails, store a fallback message
// 			newsResp.Articles[i].Summary = "summary unavailable"
// 			continue
// 		}
// 		// Save the generated summary into the article
// 		newsResp.Articles[i].Summary = summary

// 		for i := range newsResp.Articles {
// 			fullText, err := extractArticleText(newsResp.Articles[i].URL)
// 			if err != nil || strings.TrimSpace(fullText) == "" {
// 				fullText = "Title: " + newsResp.Articles[i].Title + "\n" +
// 					"Description: " + newsResp.Articles[i].Description + "\n" +
// 					"Content: " + newsResp.Articles[i].Content
// 			}

// 			summary, err := summarizeWithGroq(fullText)
// 			if err != nil {
// 				newsResp.Articles[i].Summary = "summary unavailable"
// 			} else {
// 				newsResp.Articles[i].Summary = summary
// 			}

// 			// Save to Supabase/Postgres AFTER summary is ready
// 			if err := saveArticleToDB(db, newsResp.Articles[i]); err != nil {
// 				log.Println("failed to save article:", err)
// 			}
// 		}

// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	_ = json.NewEncoder(w).Encode(newsResp)
// }

package main

import (
	"database/sql"
	"encoding/json"
	"log"
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

func saveArticleToDB(db *sql.DB, a Article) error {
	if db == nil {
		return nil
	}

	_, err := db.Exec(`
		INSERT INTO articles (
			title,
			description,
			content,
			summary,
			author,
			source_name,
			source_id,
			url,
			image_url,
			published_at
		)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (url) DO UPDATE SET
			title = EXCLUDED.title,
			description = EXCLUDED.description,
			content = EXCLUDED.content,
			summary = EXCLUDED.summary,
			author = EXCLUDED.author,
			source_name = EXCLUDED.source_name,
			source_id = EXCLUDED.source_id,
			image_url = EXCLUDED.image_url,
			published_at = EXCLUDED.published_at
	`,
		a.Title,
		a.Description,
		a.Content,
		a.Summary,
		a.Author,
		a.Source.Name,
		a.Source.ID,
		a.URL,
		a.URLToImage,
		a.PublishedAt,
	)

	return err
}

func newsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("NEWSAPI_KEY")
		if apiKey == "" {
			writeJSONError(w, http.StatusInternalServerError, "missing NEWSAPI_KEY environment variable")
			return
		}

		query := r.URL.Query().Get("q")
		if query == "" {
			query = "tesla"
		}

		newsURL := "https://newsapi.org/v2/everything?q=" +
			url.QueryEscape(query) +
			"&sortBy=publishedAt&pageSize=10&language=en"

		req, err := http.NewRequest(http.MethodGet, newsURL, nil)
		if err != nil {
			writeJSONError(w, http.StatusBadGateway, "failed to create NewsAPI request")
			return
		}

		req.Header.Set("X-Api-Key", apiKey)

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

		for i := range newsResp.Articles {
			fullText, err := extractArticleText(newsResp.Articles[i].URL)
			if err != nil || strings.TrimSpace(fullText) == "" {
				fullText = "Title: " + newsResp.Articles[i].Title + "\n" +
					"Description: " + newsResp.Articles[i].Description + "\n" +
					"Content: " + newsResp.Articles[i].Content
			}

			summary, err := summarizeWithGroq(fullText)
			if err != nil {
				newsResp.Articles[i].Summary = "summary unavailable"
			} else {
				newsResp.Articles[i].Summary = summary
			}

			// Save to Supabase/Postgres after summary is ready
			if err := saveArticleToDB(db, newsResp.Articles[i]); err != nil {
				log.Println("failed to save article:", err)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(newsResp)
	}
}
