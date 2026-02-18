package service

import (
	"context"

	"github.com/google/uuid"
	"socialnews/internal/models"
	"socialnews/internal/repository"
)

type headlineService struct {
	headlineRepo repository.HeadlineRepository
	articleRepo  repository.ArticleRepository
}

func NewHeadlineService(h repository.HeadlineRepository, a repository.ArticleRepository) HeadlineService {
	return &headlineService{headlineRepo: h, articleRepo: a}
}

func (s *headlineService) CreateFromArticles(ctx context.Context, articles []models.Article) (*models.Headline, error) {
	headline := &models.Headline{
		ID:    uuid.New(),
		Title: articles[0].Title,
	}

	if err := s.headlineRepo.Create(ctx, headline); err != nil {
		return nil, err
	}

	for _, article := range articles {
		article.HeadlineID = headline.ID
		_ = s.articleRepo.Create(ctx, &article)
	}

	return headline, nil
}

func (s *headlineService) GetByID(ctx context.Context, id uuid.UUID) (*models.Headline, error) {
	return s.headlineRepo.GetByIDWithArticles(ctx, id)
}

func (s *headlineService) GetTrending(ctx context.Context, limit int) ([]models.Headline, error) {
	return s.headlineRepo.GetTrending(ctx, limit)
}
