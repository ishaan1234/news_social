package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go/v4"
	firebaseauth "firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

type firebaseAuth struct {
	client *firebaseauth.Client
}

type firebaseUser struct {
	UID           string         `json:"uid"`
	Email         string         `json:"email,omitempty"`
	Name          string         `json:"name,omitempty"`
	Picture       string         `json:"picture,omitempty"`
	EmailVerified bool           `json:"email_verified"`
	Claims        map[string]any `json:"claims,omitempty"`
}

type firebaseUserContextKey struct{}

// newFirebaseAuth creates a Firebase auth client.
// It uses GOOGLE_APPLICATION_CREDENTIALS by default and also supports
// FIREBASE_CREDENTIALS_FILE and FIREBASE_PROJECT_ID when set.
func newFirebaseAuth(ctx context.Context) (*firebaseAuth, error) {
	config := &firebase.Config{}

	if projectID := strings.TrimSpace(os.Getenv("FIREBASE_PROJECT_ID")); projectID != "" {
		config.ProjectID = projectID
	}

	var opts []option.ClientOption
	if credentialsPath := strings.TrimSpace(os.Getenv("FIREBASE_CREDENTIALS_FILE")); credentialsPath != "" {
		opts = append(opts, option.WithCredentialsFile(credentialsPath))
	}

	app, err := firebase.NewApp(ctx, config, opts...)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %w", err)
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase auth client: %w", err)
	}

	return &firebaseAuth{client: client}, nil
}

// verifyToken validates a Firebase ID token or an Authorization header value.
func (a *firebaseAuth) verifyToken(ctx context.Context, rawToken string) (*firebaseUser, error) {
	if a == nil || a.client == nil {
		return nil, errors.New("firebase auth client is not initialized")
	}

	idToken, err := normalizeBearerToken(rawToken)
	if err != nil {
		return nil, err
	}

	token, err := a.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("failed to verify firebase id token: %w", err)
	}

	return firebaseUserFromToken(token), nil
}

// verifyRequest reads the Authorization header and verifies the Firebase ID token.
func (a *firebaseAuth) verifyRequest(r *http.Request) (*firebaseUser, error) {
	if r == nil {
		return nil, errors.New("request is nil")
	}

	return a.verifyToken(r.Context(), r.Header.Get("Authorization"))
}

// middleware protects an HTTP handler with Firebase auth and stores the user in context.
func (a *firebaseAuth) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, err := a.verifyRequest(r)
		if err != nil {
			writeJSONError(w, http.StatusUnauthorized, "invalid or missing firebase auth token")
			return
		}

		ctx := context.WithValue(r.Context(), firebaseUserContextKey{}, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func firebaseUserFromContext(ctx context.Context) (*firebaseUser, bool) {
	user, ok := ctx.Value(firebaseUserContextKey{}).(*firebaseUser)
	return user, ok
}

func normalizeBearerToken(rawToken string) (string, error) {
	token := strings.TrimSpace(rawToken)
	if token == "" {
		return "", errors.New("missing authorization token")
	}

	if strings.HasPrefix(strings.ToLower(token), "bearer ") {
		token = strings.TrimSpace(token[7:])
	}

	if token == "" {
		return "", errors.New("missing bearer token")
	}

	return token, nil
}

func firebaseUserFromToken(token *firebaseauth.Token) *firebaseUser {
	if token == nil {
		return nil
	}

	return &firebaseUser{
		UID:           token.UID,
		Email:         stringClaim(token.Claims, "email"),
		Name:          stringClaim(token.Claims, "name"),
		Picture:       stringClaim(token.Claims, "picture"),
		EmailVerified: boolClaim(token.Claims, "email_verified"),
		Claims:        token.Claims,
	}
}

func stringClaim(claims map[string]any, key string) string {
	value, ok := claims[key]
	if !ok {
		return ""
	}

	text, ok := value.(string)
	if !ok {
		return ""
	}

	return text
}

func boolClaim(claims map[string]any, key string) bool {
	value, ok := claims[key]
	if !ok {
		return false
	}

	flag, ok := value.(bool)
	if ok {
		return flag
	}

	return false
}
