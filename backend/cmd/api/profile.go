package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
)

type UserProfile struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	AvatarURL   string `json:"avatar_url"`
	Role        string `json:"role"`
	Bio         string `json:"bio"`
	Location    string `json:"location"`
	Website     string `json:"website"`
}

func profileHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database is not configured")
			return
		}

		switch r.Method {
		case http.MethodGet:
			getProfile(w, r, db)
		case http.MethodPut:
			updateProfile(w, r, db)
		default:
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		}
	}
}

func getProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	email := r.URL.Query().Get("email")
	if email == "" {
		writeJSONError(w, http.StatusBadRequest, "email query parameter is required")
		return
	}

	var p UserProfile
	var role, bio, location, website, avatar sql.NullString
	err := db.QueryRowContext(r.Context(), `
		SELECT email, username, display_name, avatar_url, role, bio, location, website
		FROM users WHERE email = $1
	`, email).Scan(&p.Email, &p.Username, &p.DisplayName, &avatar, &role, &bio, &location, &website)

	if err != nil {
		if err == sql.ErrNoRows {
			writeJSONError(w, http.StatusNotFound, "user profile not found")
		} else {
			writeJSONError(w, http.StatusInternalServerError, "failed to get profile")
		}
		return
	}

	p.AvatarURL = avatar.String
	p.Role = role.String
	p.Bio = bio.String
	p.Location = location.String
	p.Website = website.String

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    p,
	})
}

func updateProfile(w http.ResponseWriter, r *http.Request, db *sql.DB) {
	if !hasAuthorizationHeader(r) {
		writeJSONError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req UserProfile
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	// Verify the user is updating their own profile
	sessionEmail, err := requestUserEmail(r, req.Email)
	if err != nil || sessionEmail != req.Email {
		writeJSONError(w, http.StatusForbidden, "not allowed to update this profile")
		return
	}

	_, err = db.ExecContext(r.Context(), `
		UPDATE users
		SET display_name = $1, username = $2, role = $3, bio = $4, location = $5, website = $6
		WHERE email = $7
	`, req.DisplayName, req.Username, req.Role, req.Bio, req.Location, req.Website, sessionEmail)

	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to update profile")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string]any{
		"success": true,
		"data":    req,
	})
}
