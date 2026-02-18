package service

import (
	"context"
	"github.com/google/uuid"
	"socialnews/internal/models"
)

type AuthService interface {
	Register(ctx context.Context, email, password string) (*models.User, error)
	Login(ctx context.Context, email, password string) (string, error)
	ValidateToken(ctx context.Context, token string) (*models.User, error)
}

type HeadlineService interface {
	CreateFromArticles(ctx context.Context, articles []models.Article) (*models.Headline, error)
	GetByID(ctx context.Context, id uuid.UUID) (*models.Headline, error)
	GetTrending(ctx context.Context, limit int) ([]models.Headline, error)
}

type ArticleService interface {
	Store(ctx context.Context, article *models.Article) error
	GetByHeadline(ctx context.Context, headlineID uuid.UUID) ([]models.Article, error)
}

type SummaryService interface {
	Generate(ctx context.Context, headlineID uuid.UUID) (*models.Summary, error)
}

type SocialService interface {
	AddComment(ctx context.Context, userID, headlineID uuid.UUID, content string) error
	Vote(ctx context.Context, userID, headlineID uuid.UUID, value int) error
	GetComments(ctx context.Context, headlineID uuid.UUID) ([]models.Comment, error)
}
