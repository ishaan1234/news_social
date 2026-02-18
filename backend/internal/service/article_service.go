package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ishaan1234/news_social/backend/internal/models"
	// "github.com/ishaan1234/news_social/backend/internal/repository"
)

type articleService struct {
	// articleRepo repository.ArticleRepository
}

// func NewArticleService(r repository.ArticleRepository) ArticleService {
// 	return &articleService{articleRepo: r}
// }

func (s *articleService) Store(ctx context.Context, article *models.Article) error {
	// exists, _ := s.articleRepo.ExistsByURL(ctx, article.URL)
	// if exists {
	// 	return nil
	// }
	// return s.articleRepo.Create(ctx, article)
	return nil
}

func (s *articleService) GetByHeadline(ctx context.Context, headlineID uuid.UUID) ([]models.Article, error) {
	// return s.articleRepo.GetByHeadlineID(ctx, headlineID)
	return nil, nil
}
