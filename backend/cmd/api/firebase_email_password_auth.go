package main

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/mail"
	"os"
	"regexp"
	"strings"
	"time"

	firebaseauth "firebase.google.com/go/v4/auth"
	"github.com/lib/pq"
)

const firebaseIdentityToolkitBaseURL = "https://identitytoolkit.googleapis.com/v1"

type firebaseIdentityClient struct {
	apiKey string
	client *http.Client
}

type firebaseIdentityAuthResponse struct {
	LocalID      string `json:"localId"`
	Email        string `json:"email"`
	DisplayName  string `json:"displayName"`
	IDToken      string `json:"idToken"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	Registered   bool   `json:"registered"`
}

type firebaseIdentitySendOOBResponse struct {
	Email string `json:"email"`
}

type firebaseIdentityErrorResponse struct {
	Error struct {
		Message string `json:"message"`
	} `json:"error"`
}

type signupRequest struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Username    string `json:"username,omitempty"`
	DisplayName string `json:"display_name,omitempty"`
	AvatarURL   string `json:"avatar_url,omitempty"`
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type verifyEmailResendRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type authSuccessResponse struct {
	Success      bool                         `json:"success"`
	Message      string                       `json:"message,omitempty"`
	User         *firebaseIdentityUserPayload `json:"user,omitempty"`
	IDToken      string                       `json:"id_token,omitempty"`
	RefreshToken string                       `json:"refresh_token,omitempty"`
	ExpiresIn    string                       `json:"expires_in,omitempty"`
}

type firebaseIdentityUserPayload struct {
	UID           string `json:"uid"`
	Email         string `json:"email"`
	Username      string `json:"username,omitempty"`
	DisplayName   string `json:"display_name,omitempty"`
	AvatarURL     string `json:"avatar_url,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
}

func registerFirebaseEmailPasswordRoutes(mux *http.ServeMux) {
	mux.HandleFunc("/auth/signup", signupHandler)
	mux.HandleFunc("/auth/login", loginHandler)
	mux.HandleFunc("/auth/verify-email/resend", resendVerificationEmailHandler)
}

func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req signupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	password := strings.TrimSpace(req.Password)
	if len(password) < 6 {
		writeJSONError(w, http.StatusBadRequest, "password must be at least 6 characters")
		return
	}

	displayName := strings.TrimSpace(req.DisplayName)
	avatarURL := strings.TrimSpace(req.AvatarURL)
	username := normalizeUsername(req.Username)

	client, err := newFirebaseIdentityClient()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authResp, err := client.signUpWithEmailPassword(r.Context(), email, password)
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "EMAIL_EXISTS"):
			handleExistingEmailDuringSignup(w, r, client, email, password, username, displayName, avatarURL)
			return
		case strings.Contains(err.Error(), "WEAK_PASSWORD"):
			writeJSONError(w, http.StatusBadRequest, err.Error())
			return
		default:
			writeJSONError(w, http.StatusBadGateway, err.Error())
			return
		}
	}

	if displayName != "" {
		updatedResp, err := client.updateProfile(r.Context(), authResp.IDToken, displayName)
		if err == nil && updatedResp != nil && strings.TrimSpace(updatedResp.DisplayName) != "" {
			authResp.DisplayName = updatedResp.DisplayName
		}
	}

	username = usernameOrDefault(username, email, authResp.LocalID)
	if err := saveSignupUserToSupabase(r.Context(), signupUser{
		Email:       email,
		Username:    username,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}, false); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "account created, but failed to save user details: "+err.Error())
		return
	}

	if _, err := client.sendEmailVerification(r.Context(), authResp.IDToken); err != nil {
		writeJSONError(w, http.StatusBadGateway, "account created, but failed to send verification email: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authSuccessResponse{
		Success:      true,
		Message:      "signup successful; verification email sent",
		IDToken:      authResp.IDToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    authResp.ExpiresIn,
		User: &firebaseIdentityUserPayload{
			UID:           authResp.LocalID,
			Email:         authResp.Email,
			Username:      username,
			DisplayName:   authResp.DisplayName,
			AvatarURL:     avatarURL,
			EmailVerified: false,
		},
	})
}

func handleExistingEmailDuringSignup(
	w http.ResponseWriter,
	r *http.Request,
	client *firebaseIdentityClient,
	email string,
	password string,
	username string,
	displayName string,
	avatarURL string,
) {
	ctx := r.Context()

	adminAuth, err := newFirebaseAuth(ctx)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to check existing account state: "+err.Error())
		return
	}

	u, err := adminAuth.client.GetUserByEmail(ctx, email)
	if err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to load existing account: "+err.Error())
		return
	}

	if u.EmailVerified {
		writeJSONError(w, http.StatusConflict, "account already exists. Please log in.")
		return
	}

	authResp, err := client.signInWithEmailPassword(ctx, email, password)
	if err != nil {
		writeJSONError(w, http.StatusConflict, "your account already exists, but your email is not verified. Please verify it using your original password or reset the password first.")
		return
	}

	if _, err := client.sendEmailVerification(ctx, authResp.IDToken); err != nil {
		writeJSONError(w, http.StatusBadGateway, "your account already exists, but your email is not verified. Failed to resend verification email: "+err.Error())
		return
	}

	username = usernameOrDefault(username, email, u.UID)
	if displayName == "" {
		displayName = strings.TrimSpace(u.DisplayName)
	}
	if err := saveSignupUserToSupabase(ctx, signupUser{
		Email:       email,
		Username:    username,
		DisplayName: displayName,
		AvatarURL:   avatarURL,
	}, false); err != nil {
		writeJSONError(w, http.StatusInternalServerError, "verification email sent, but failed to save user details: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authSuccessResponse{
		Success:      true,
		Message:      "your account already exists, but your email is not verified. We’ve sent a new verification email.",
		IDToken:      authResp.IDToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    authResp.ExpiresIn,
		User: &firebaseIdentityUserPayload{
			UID:           u.UID,
			Email:         u.Email,
			Username:      username,
			DisplayName:   u.DisplayName,
			AvatarURL:     avatarURL,
			EmailVerified: false,
		},
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req loginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	password := strings.TrimSpace(req.Password)
	if password == "" {
		writeJSONError(w, http.StatusBadRequest, "password is required")
		return
	}

	adminAuth, err := newFirebaseAuth(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to initialize firebase admin auth: "+err.Error())
		return
	}

	u, err := adminAuth.client.GetUserByEmail(r.Context(), email)
	if err != nil {
		if firebaseauth.IsUserNotFound(err) {
			writeJSONError(w, http.StatusNotFound, "you have not signed up yet.")
			return
		}
		writeJSONError(w, http.StatusBadGateway, "failed to load account: "+err.Error())
		return
	}

	if !u.EmailVerified {
		writeJSONError(w, http.StatusForbidden, "your email is not verified yet. Please verify it first.")
		return
	}

	client, err := newFirebaseIdentityClient()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authResp, err := client.signInWithEmailPassword(r.Context(), email, password)
	if err != nil {
		status := http.StatusUnauthorized
		if strings.Contains(err.Error(), "OPERATION_NOT_ALLOWED") {
			status = http.StatusBadRequest
		}
		writeJSONError(w, status, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authSuccessResponse{
		Success:      true,
		Message:      "login successful",
		IDToken:      authResp.IDToken,
		RefreshToken: authResp.RefreshToken,
		ExpiresIn:    authResp.ExpiresIn,
		User: &firebaseIdentityUserPayload{
			UID:           authResp.LocalID,
			Email:         authResp.Email,
			DisplayName:   authResp.DisplayName,
			EmailVerified: true,
		},
	})
}

func resendVerificationEmailHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req verifyEmailResendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	email, err := normalizeEmail(req.Email)
	if err != nil {
		writeJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	password := strings.TrimSpace(req.Password)
	if password == "" {
		writeJSONError(w, http.StatusBadRequest, "password is required")
		return
	}

	adminAuth, err := newFirebaseAuth(r.Context())
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, "failed to initialize firebase admin auth: "+err.Error())
		return
	}

	u, err := adminAuth.client.GetUserByEmail(r.Context(), email)
	if err != nil {
		if firebaseauth.IsUserNotFound(err) {
			writeJSONError(w, http.StatusNotFound, "you have not signed up yet.")
			return
		}
		writeJSONError(w, http.StatusBadGateway, "failed to load account: "+err.Error())
		return
	}

	if u.EmailVerified {
		writeJSONError(w, http.StatusBadRequest, "your email is already verified.")
		return
	}

	client, err := newFirebaseIdentityClient()
	if err != nil {
		writeJSONError(w, http.StatusInternalServerError, err.Error())
		return
	}

	authResp, err := client.signInWithEmailPassword(r.Context(), email, password)
	if err != nil {
		writeJSONError(w, http.StatusUnauthorized, "unable to resend verification email with the provided credentials")
		return
	}

	if _, err := client.sendEmailVerification(r.Context(), authResp.IDToken); err != nil {
		writeJSONError(w, http.StatusBadGateway, "failed to send verification email: "+err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(authSuccessResponse{
		Success: true,
		Message: "verification email sent",
		User: &firebaseIdentityUserPayload{
			UID:           u.UID,
			Email:         u.Email,
			DisplayName:   u.DisplayName,
			EmailVerified: false,
		},
	})
}

func newFirebaseIdentityClient() (*firebaseIdentityClient, error) {
	apiKey := strings.TrimSpace(os.Getenv("FIREBASE_WEB_API_KEY"))
	if apiKey == "" {
		return nil, errors.New("missing FIREBASE_WEB_API_KEY environment variable")
	}

	return &firebaseIdentityClient{
		apiKey: apiKey,
		client: &http.Client{Timeout: 10 * time.Second},
	}, nil
}

func (c *firebaseIdentityClient) signUpWithEmailPassword(ctx context.Context, email, password string) (*firebaseIdentityAuthResponse, error) {
	return c.postAuth(ctx, "accounts:signUp", map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})
}

func (c *firebaseIdentityClient) signInWithEmailPassword(ctx context.Context, email, password string) (*firebaseIdentityAuthResponse, error) {
	return c.postAuth(ctx, "accounts:signInWithPassword", map[string]any{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})
}

func (c *firebaseIdentityClient) updateProfile(ctx context.Context, idToken, displayName string) (*firebaseIdentityAuthResponse, error) {
	if strings.TrimSpace(displayName) == "" {
		return nil, nil
	}

	return c.postAuth(ctx, "accounts:update", map[string]any{
		"idToken":           idToken,
		"displayName":       displayName,
		"returnSecureToken": true,
	})
}

func (c *firebaseIdentityClient) sendEmailVerification(ctx context.Context, idToken string) (*firebaseIdentitySendOOBResponse, error) {
	return c.postOOB(ctx, map[string]any{
		"requestType": "VERIFY_EMAIL",
		"idToken":     idToken,
	})
}

func (c *firebaseIdentityClient) postAuth(ctx context.Context, endpoint string, payload map[string]any) (*firebaseIdentityAuthResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode firebase auth request: %w", err)
	}

	url := fmt.Sprintf("%s/%s?key=%s", firebaseIdentityToolkitBaseURL, endpoint, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create firebase auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call firebase auth api: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr firebaseIdentityErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("firebase auth error: %s", apiErr.Error.Message)
		}
		return nil, fmt.Errorf("firebase auth returned status %d", res.StatusCode)
	}

	var response firebaseIdentityAuthResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode firebase auth response: %w", err)
	}

	return &response, nil
}

func (c *firebaseIdentityClient) postOOB(ctx context.Context, payload map[string]any) (*firebaseIdentitySendOOBResponse, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to encode firebase oob request: %w", err)
	}

	url := fmt.Sprintf("%s/accounts:sendOobCode?key=%s", firebaseIdentityToolkitBaseURL, c.apiKey)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create firebase oob request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call firebase oob api: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		var apiErr firebaseIdentityErrorResponse
		if err := json.NewDecoder(res.Body).Decode(&apiErr); err == nil && apiErr.Error.Message != "" {
			return nil, fmt.Errorf("firebase auth error: %s", apiErr.Error.Message)
		}
		return nil, fmt.Errorf("firebase auth returned status %d", res.StatusCode)
	}

	var response firebaseIdentitySendOOBResponse
	if err := json.NewDecoder(res.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to decode firebase oob response: %w", err)
	}

	return &response, nil
}

func normalizeEmail(value string) (string, error) {
	email := strings.ToLower(strings.TrimSpace(value))
	if email == "" {
		return "", errors.New("email is required")
	}

	if _, err := mail.ParseAddress(email); err != nil {
		return "", errors.New("invalid email address")
	}

	return email, nil
}

type signupUser struct {
	Email       string
	Username    string
	DisplayName string
	AvatarURL   string
}

func saveSignupUserToSupabase(ctx context.Context, user signupUser, retriedUsername bool) error {
	dbURL := strings.TrimSpace(os.Getenv("DATABASE_URL"))
	if dbURL == "" {
		return nil
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return fmt.Errorf("connect database: %w", err)
	}

	if err := ensureSupabaseUsersTable(ctx, db); err != nil {
		return err
	}

	_, err = db.ExecContext(ctx, `
		INSERT INTO users (email, username, display_name, avatar_url)
		VALUES ($1, $2, NULLIF($3, ''), NULLIF($4, ''))
		ON CONFLICT (email) DO UPDATE
		SET username = EXCLUDED.username,
		    display_name = EXCLUDED.display_name,
		    avatar_url = EXCLUDED.avatar_url
	`, user.Email, user.Username, user.DisplayName, user.AvatarURL)
	if err == nil {
		return nil
	}

	if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == "23505" && !retriedUsername {
		user.Username = usernameOrDefault("", user.Email, time.Now().Format("150405"))
		return saveSignupUserToSupabase(ctx, user, true)
	}

	return fmt.Errorf("save user: %w", err)
}

func ensureSupabaseUsersTable(ctx context.Context, db *sql.DB) error {
	_, err := db.ExecContext(ctx, `
		CREATE TABLE IF NOT EXISTS users (
			email TEXT PRIMARY KEY,
			username TEXT UNIQUE,
			display_name TEXT,
			avatar_url TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("ensure users table: %w", err)
	}
	return nil
}

func normalizeUsername(value string) string {
	username := strings.ToLower(strings.TrimSpace(value))
	username = regexp.MustCompile(`[^a-z0-9_]+`).ReplaceAllString(username, "_")
	username = strings.Trim(username, "_")
	if len(username) > 30 {
		username = username[:30]
	}
	return username
}

func usernameOrDefault(username, email, suffix string) string {
	if username != "" {
		return username
	}

	base := strings.Split(email, "@")[0]
	base = normalizeUsername(base)
	if base == "" {
		base = "user"
	}

	suffix = normalizeUsername(suffix)
	if len(suffix) > 8 {
		suffix = suffix[:8]
	}
	if suffix == "" {
		return base
	}

	maxBaseLength := 30 - len(suffix) - 1
	if maxBaseLength < 1 {
		maxBaseLength = 1
	}
	if len(base) > maxBaseLength {
		base = base[:maxBaseLength]
	}

	return base + "_" + suffix
}
