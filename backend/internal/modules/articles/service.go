package articles

import "social-news/internal/models"

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetByHeadline(headlineID int64) ([]models.Article, error) {
	return s.repo.FindByHeadline(headlineID)
}
