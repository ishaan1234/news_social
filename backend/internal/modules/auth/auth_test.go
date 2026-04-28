package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signedToken(t *testing.T, secret string, userID int) string {
	t.Helper()
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour).Unix(),
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("failed to sign token: %v", err)
	}
	return token
}

func TestService_ValidateToken(t *testing.T) {
	service := NewService("secret")
	token := signedToken(t, "secret", 42)

	userID, err := service.ValidateToken(token)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if userID != 42 {
		t.Fatalf("expected user id 42, got %d", userID)
	}
}

func TestService_ValidateToken_Empty(t *testing.T) {
	service := NewService("secret")

	_, err := service.ValidateToken("")
	if err == nil {
		t.Fatalf("expected error")
	}
}
