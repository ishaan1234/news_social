package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/ishaan1234/news_social/backend/internal/models"
	// "github.com/ishaan1234/news_social/backend/internal/repository"
)

type AIClient interface {
	Summarize(ctx context.Context, articles []models.Article) (string, error)
}

type summaryService struct {
	// headlineRepo repository.HeadlineRepository
	// articleRepo  repository.ArticleRepository
	// summaryRepo  repository.SummaryRepository
	aiClient AIClient
}

// func NewSummaryService(h repository.HeadlineRepository,
// 	a repository.ArticleRepository,
// 	s repository.SummaryRepository,
// 	ai AIClient) SummaryService {
//
// 	return &summaryService{
// 		headlineRepo: h,
// 		articleRepo:  a,
// 		summaryRepo:  s,
// 		aiClient:     ai,
// 	}
// }

func (s *summaryService) Generate(ctx context.Context, headlineID uuid.UUID) (*models.Summary, error) {
	// articles, err := s.articleRepo.GetByHeadlineID(ctx, headlineID)
	// if err != nil {
	// 	return nil, err
	// }

	// text, err := s.aiClient.Summarize(ctx, articles)
	// if err != nil {
	// 	return nil, err
	// }

	summary := &models.Summary{
		ID:         uuid.New(),
		HeadlineID: headlineID,
		// Content:    text,
	}

	// return summary, s.summaryRepo.Create(ctx, summary)
	return summary, nil
}
