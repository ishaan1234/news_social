package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/lib/pq"
)

type followRequest struct {
	FollowerEmail  string `json:"follower_email"`
	FollowingEmail string `json:"following_email"`
}

type followResponse struct {
	FollowerEmail  string    `json:"follower_email"`
	FollowingEmail string    `json:"following_email"`
	CreatedAt      time.Time `json:"created_at,omitempty"`
}

func followingHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database is not configured")
			return
		}

		switch r.Method {
		case http.MethodPost:
			followUser(w, r, db)
		case http.MethodDelete:
			unfollowUser(w, r, db)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func followUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	req, ok := decodeFollowRequest(w, r)
	if !ok {
		return
	}

	var follow followResponse
	err := db.QueryRowContext(r.Context(), `
		INSERT INTO following (follower_email, following_email)
		VALUES ($1, $2)
		RETURNING follower_email, following_email, created_at
	`, req.FollowerEmail, req.FollowingEmail).Scan(
		&follow.FollowerEmail,
		&follow.FollowingEmail,
		&follow.CreatedAt,
	)
	if err != nil {
		status, message := followInsertError(err)
		writeJSONError(w, status, message)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"follow":  follow,
	})
}

func unfollowUser(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	req, ok := decodeFollowRequest(w, r)
	if !ok {
		return
	}

	result, err := db.ExecContext(r.Context(), `
		DELETE FROM following
		WHERE follower_email = $1 AND following_email = $2
	`, req.FollowerEmail, req.FollowingEmail)
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to unfollow user")
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to check unfollow result")
		return
	}
	if rowsAffected == 0 {
		writeJSONError(w, http.StatusNotFound, "follow relationship does not exist")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"message": "unfollowed user",
	})
}

func decodeFollowRequest(w http.ResponseWriter, r *http.Request) (followRequest, bool) {
	var req followRequest
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return followRequest{}, false
	}

	followerEmail, err := requestUserEmail(r, req.FollowerEmail)
	if err != nil {
		if hasAuthorizationHeader(r) {
			writeJSONError(w, http.StatusUnauthorized, "valid authenticated user is required")
		} else {
			writeJSONError(w, http.StatusBadRequest, "valid follower_email is required")
		}
		return followRequest{}, false
	}

	followingEmail, err := normalizeEmail(req.FollowingEmail)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, "valid following_email is required")
		return followRequest{}, false
	}

	if followerEmail == followingEmail {
		writeJSONError(w, http.StatusBadRequest, "users cannot follow themselves")
		return followRequest{}, false
	}

	return followRequest{
		FollowerEmail:  followerEmail,
		FollowingEmail: followingEmail,
	}, true
}

func followInsertError(err error) (int, string) {
	var pqErr *pq.Error
	if !errors.As(err, &pqErr) {
		return http.StatusInternalServerError, "failed to follow user"
	}

	switch pqErr.Code {
	case "23505":
		return http.StatusConflict, "user is already following this account"
	case "23503":
		return http.StatusBadRequest, "follower_email or following_email does not exist"
	case "23514":
		return http.StatusBadRequest, "users cannot follow themselves"
	default:
		return http.StatusInternalServerError, "failed to follow user"
	}
}
