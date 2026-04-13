package headlines

import (
	"context"
	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateHeadline(ctx context.Context, title string) (int, error) {
	return s.repo.Create(ctx, title)
}

func (s *Service) GetHeadlines(ctx context.Context) ([]models.Headline, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetHeadline(ctx context.Context, id int) (models.Headline, error) {
	return s.repo.GetByID(ctx, id)
}