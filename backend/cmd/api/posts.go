package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type createPostRequest struct {
	UserEmail string `json:"user_email"`
	ArticleID string `json:"article_id"`
	Caption   string `json:"caption"`
}

type postResponse struct {
	ID        string    `json:"id"`
	UserEmail string    `json:"user_email"`
	ArticleID string    `json:"article_id"`
	Caption   string    `json:"caption,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

func createPostHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database is not configured")
			return
		}

		var req createPostRequest
		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		if err := decoder.Decode(&req); err != nil {
			writeJSONError(w, http.StatusBadRequest, "invalid request body")
			return
		}

		userEmail, err := requestUserEmail(r, req.UserEmail)
		if err != nil {
			if hasAuthorizationHeader(r) {
				writeJSONError(w, http.StatusUnauthorized, "valid authenticated user is required")
			} else {
				writeJSONError(w, http.StatusBadRequest, err.Error())
			}
			return
		}

		articleID := strings.TrimSpace(req.ArticleID)
		if _, err := uuid.Parse(articleID); err != nil {
			writeJSONError(w, http.StatusBadRequest, "valid article_id is required")
			return
		}

		post, err := insertPost(r, db, userEmail, articleID, strings.TrimSpace(req.Caption))
		if err != nil {
			status, message := postInsertError(err)
			writeJSONError(w, status, message)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"post":    post,
		})
	}
}

func insertPost(r *http.Request, db *sql.DB, userEmail, articleID, caption string) (postResponse, error) {
	var post postResponse
	err := db.QueryRowContext(r.Context(), `
		INSERT INTO posts (user_email, article_id, caption)
		VALUES ($1, $2, NULLIF($3, ''))
		RETURNING id::text, user_email, article_id::text, COALESCE(caption, ''), created_at
	`, userEmail, articleID, caption).Scan(
		&post.ID,
		&post.UserEmail,
		&post.ArticleID,
		&post.Caption,
		&post.CreatedAt,
	)
	if err != nil {
		return postResponse{}, err
	}
	return post, nil
}

func postInsertError(err error) (int, string) {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return http.StatusInternalServerError, "failed to create post"
	}

	switch pqErr.Code {
	case "23505":
		return http.StatusConflict, "user has already posted this article"
	case "23503":
		return http.StatusBadRequest, "user_email or article_id does not exist"
	default:
		return http.StatusInternalServerError, "failed to create post"
	}
}
