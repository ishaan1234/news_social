package auth

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type Service struct {
	jwtSecret string
}

func NewService(jwtSecret string) *Service {
	return &Service{jwtSecret: jwtSecret}
}

func (s *Service) ValidateToken(tokenString string) (int, error) {
	if strings.TrimSpace(tokenString) == "" {
		return 0, fmt.Errorf("token is required")
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid token claims")
	}

	switch v := claims["user_id"].(type) {
	case float64:
		if v <= 0 {
			return 0, fmt.Errorf("invalid user_id")
		}
		return int(v), nil
	case string:
		id, err := strconv.Atoi(v)
		if err != nil || id <= 0 {
			return 0, fmt.Errorf("invalid user_id")
		}
		return id, nil
	default:
		return 0, fmt.Errorf("invalid user_id")
	}
}
