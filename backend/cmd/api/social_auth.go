package main

import (
	"fmt"
	"net/http"
	"strings"
)

func requestUserEmail(r *http.Request, fallbackEmail string) (string, error) {
	if hasAuthorizationHeader(r) {
		authClient, err := newFirebaseAuth(r.Context())
		if err != nil {
			return "", fmt.Errorf("initialize auth: %w", err)
		}

		user, err := authClient.verifyRequest(r)
		if err != nil {
			return "", err
		}

		return normalizeEmail(user.Email)
	}

	return normalizeEmail(fallbackEmail)
}

func hasAuthorizationHeader(r *http.Request) bool {
	return strings.TrimSpace(r.Header.Get("Authorization")) != ""
}
