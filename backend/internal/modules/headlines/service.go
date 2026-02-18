package headlines

import "github.com/ishaan1234/news_social/backend/internal/models"

type Service struct {
	repo ArticleAggregator
}

func NewService(repo ArticleAggregator) *Service {
	return &Service{repo: repo}
}

func (s *Service) ListHeadlines() ([]models.Headline, error) {
	return s.repo.ListHeadlines()
}

func (s *Service) GetHeadlineDetails(id string) (any, error) {
	return s.repo.FetchHeadlineBundle(id)
}
