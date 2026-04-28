package headlines

import (
	"context"
	"fmt"
	"strings"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateHeadline(ctx context.Context, title string) (int, error) {
	if strings.TrimSpace(title) == "" {
		return 0, fmt.Errorf("title is required")
	}
	return s.repo.Create(ctx, strings.TrimSpace(title))
}

func (s *Service) GetHeadlines(ctx context.Context) ([]models.Headline, error) {
	return s.repo.GetAll(ctx)
}

func (s *Service) GetHeadline(ctx context.Context, id int) (models.Headline, error) {
	if id <= 0 {
		return models.Headline{}, fmt.Errorf("valid headline id is required")
	}
	return s.repo.GetByID(ctx, id)
}
