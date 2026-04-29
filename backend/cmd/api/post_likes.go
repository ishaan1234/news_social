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

type postLikeRequest struct {
	UserEmail string `json:"user_email"`
	PostID    string `json:"post_id"`
}

type postLikeResponse struct {
	UserEmail string    `json:"user_email"`
	PostID    string    `json:"post_id"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}

func postLikesHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database is not configured")
			return
		}

		switch r.Method {
		case http.MethodPost:
			likePost(w, r, db)
		case http.MethodDelete:
			unlikePost(w, r, db)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func likePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	req, ok := decodePostLikeRequest(w, r)
	if !ok {
		return
	}

	var like postLikeResponse
	err := db.QueryRowContext(r.Context(), `
		INSERT INTO post_likes (user_email, post_id)
		VALUES ($1, $2)
		RETURNING user_email, post_id::text, created_at
	`, req.UserEmail, req.PostID).Scan(&like.UserEmail, &like.PostID, &like.CreatedAt)
	if err != nil {
		status, message := postLikeInsertError(err)
		writeJSONError(w, status, message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"like":    like,
	})
}

func unlikePost(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	req, ok := decodePostLikeRequest(w, r)
	if !ok {
		return
	}

	result, err := db.ExecContext(r.Context(), `
		DELETE FROM post_likes
		WHERE user_email = $1 AND post_id = $2
	`, req.UserEmail, req.PostID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to unlike post")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to check unlike result")
		return
	}
	if rowsAffected == 0 {
		writeJSONError(w, http.StatusNotFound, "like does not exist")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "unliked post",
	})
}

func decodePostLikeRequest(w http.ResponseWriter, r *http.Request) (postLikeRequest, bool) {
	var req postLikeRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return postLikeRequest{}, false
	}

	userEmail, err := normalizeEmail(req.UserEmail)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "valid user_email is required")
		return postLikeRequest{}, false
	}

	postID := strings.TrimSpace(req.PostID)
	if _, err := uuid.Parse(postID); err != nil {
		writeJSONError(w, http.StatusBadRequest, "valid post_id is required")
		return postLikeRequest{}, false
	}

	return postLikeRequest{UserEmail: userEmail, PostID: postID}, true
}

func postLikeInsertError(err error) (int, string) {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return http.StatusInternalServerError, "failed to like post"
	}

	switch pqErr.Code {
	case "23505":
		return http.StatusConflict, "user has already liked this post"
	case "23503":
		return http.StatusBadRequest, "user_email or post_id does not exist"
	default:
		return http.StatusInternalServerError, "failed to like post"
	}
}
