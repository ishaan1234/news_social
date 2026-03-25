package articles

import "github.com/ishaan1234/news_social/backend/internal/models"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByHeadline(headlineID int64) ([]models.Article, error) {
	return s.repo.FindByHeadline(headlineID)
}
