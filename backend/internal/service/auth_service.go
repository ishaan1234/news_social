package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"github.com/ishaan1234/news_social/backend/internal/models"
	// "github.com/ishaan1234/news_social/backend/internal/repository"
)

type authService struct {
	// userRepo  repository.UserRepository
	jwtSecret string
}

// func NewAuthService(userRepo repository.UserRepository, secret string) AuthService {
// 	return &authService{userRepo: userRepo, jwtSecret: secret}
// }

func (s *authService) Register(ctx context.Context, email, password string) (*models.User, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		ID:       uuid.New(),
		Email:    email,
		Password: string(hashed),
	}

	// return user, s.userRepo.Create(ctx, user)
	return user, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (string, error) {
	// user, err := s.userRepo.GetByEmail(ctx, email)
	// if err != nil {
	// 	return "", errors.New("invalid credentials")
	// }
	var user *models.User
	_ = user
	_ = errors.New("")

	// if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
	// 	return "", errors.New("invalid credentials")
	// }

	claims := jwt.MapClaims{
		"user_id": "placeholder",
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) ValidateToken(ctx context.Context, tokenStr string) (*models.User, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	// claims := token.Claims.(jwt.MapClaims)
	// userID, _ := uuid.Parse(claims["user_id"].(string))
	// return s.userRepo.GetByID(ctx, userID)
	return nil, nil
}
