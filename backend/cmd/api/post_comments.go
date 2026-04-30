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

type createPostCommentRequest struct {
	PostID    string `json:"post_id"`
	UserEmail string `json:"user_email"`
	Content   string `json:"content"`
}

type postCommentResponse struct {
	ID          string    `json:"id"`
	PostID      string    `json:"post_id"`
	UserEmail   string    `json:"user_email"`
	Username    string    `json:"username,omitempty"`
	DisplayName string    `json:"display_name,omitempty"`
	AvatarURL   string    `json:"avatar_url,omitempty"`
	Content     string    `json:"content"`
	CreatedAt   time.Time `json:"created_at"`
}

func postCommentsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database is not configured")
			return
		}

		switch r.Method {
		case http.MethodPost:
			createPostComment(w, r, db)
		case http.MethodGet:
			getPostComments(w, r, db)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func createPostComment(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	var req createPostCommentRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	postID := strings.TrimSpace(req.PostID)
	if _, err := uuid.Parse(postID); err != nil {
		writeJSONError(w, http.StatusBadRequest, "valid post_id is required")
		return
	}

	userEmail, err := requestUserEmail(r, req.UserEmail)
	if err != nil {
		if hasAuthorizationHeader(r) {
			writeJSONError(w, http.StatusUnauthorized, "valid authenticated user is required")
		} else {
			writeJSONError(w, http.StatusBadRequest, "valid user_email is required")
		}
		return
	}

	content := strings.TrimSpace(req.Content)
	if content == "" {
		writeJSONError(w, http.StatusBadRequest, "content is required")
		return
	}

	comment, err := insertPostComment(r, db, postID, userEmail, content)
	if err != nil {
		status, message := postCommentInsertError(err)
		writeJSONError(w, status, message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"comment": comment,
	})
}

func getPostComments(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	postID := strings.TrimSpace(r.URL.Query().Get("post_id"))
	if _, err := uuid.Parse(postID); err != nil {
		writeJSONError(w, http.StatusBadRequest, "valid post_id is required")
		return
	}

	comments, err := fetchPostComments(r, db, postID)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to fetch comments")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success":  true,
		"comments": comments,
	})
}

func insertPostComment(r *http.Request, db *sql.DB, postID, userEmail, content string) (postCommentResponse, error) {
	var comment postCommentResponse
	err := db.QueryRowContext(r.Context(), `
		INSERT INTO post_comments (post_id, user_email, content)
		VALUES ($1, $2, $3)
		RETURNING id::text, post_id::text, user_email, content, created_at
	`, postID, userEmail, content).Scan(
		&comment.ID,
		&comment.PostID,
		&comment.UserEmail,
		&comment.Content,
		&comment.CreatedAt,
	)
	if err != nil {
		return postCommentResponse{}, err
	}
	return comment, nil
}

func fetchPostComments(r *http.Request, db *sql.DB, postID string) ([]postCommentResponse, error) {
	rows, err := db.QueryContext(r.Context(), `
		SELECT
			c.id::text,
			c.post_id::text,
			c.user_email,
			COALESCE(u.username, ''),
			COALESCE(u.display_name, ''),
			COALESCE(u.avatar_url, ''),
			c.content,
			c.created_at
		FROM post_comments c
		JOIN users u ON u.email = c.user_email
		WHERE c.post_id = $1
		ORDER BY c.created_at ASC
	`, postID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []postCommentResponse
	for rows.Next() {
		var comment postCommentResponse
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.UserEmail,
			&comment.Username,
			&comment.DisplayName,
			&comment.AvatarURL,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func postCommentInsertError(err error) (int, string) {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return http.StatusInternalServerError, "failed to create comment"
	}

	switch pqErr.Code {
	case "23503":
		return http.StatusBadRequest, "post_id or user_email does not exist"
	default:
		return http.StatusInternalServerError, "failed to create comment"
	}
}
