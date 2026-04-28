package articles

import (
	"context"
	"fmt"

	"github.com/ishaan1234/news_social/backend/internal/models"
	"github.com/ishaan1234/news_social/backend/pkg/clients/newsapi"
)

type NewsClient interface {
	GetTopHeadlines(topic string) ([]newsapi.Article, error)
}

type Service struct {
	newsClient NewsClient
	repo       Repository
}

func NewService(newsClient NewsClient, repo Repository) *Service {
	return &Service{newsClient: newsClient, repo: repo}
}

func (s *Service) FetchAndSaveArticles(ctx context.Context, headlineID int, topic string) ([]models.Article, error) {
	if s.newsClient == nil {
		return nil, fmt.Errorf("news client is not configured")
	}

	rawArticles, err := s.newsClient.GetTopHeadlines(topic)
	if err != nil {
		return nil, fmt.Errorf("fetch articles: %w", err)
	}

	items := make([]models.Article, 0, len(rawArticles))
	for _, a := range rawArticles {
		items = append(items, models.Article{
			HeadlineID: headlineID,
			Source:     a.Source,
			Title:      a.Title,
			URL:        a.URL,
			Content:    a.Content,
		})
	}

	if s.repo != nil && len(items) > 0 {
		if err := s.repo.SaveArticles(ctx, headlineID, items); err != nil {
			return nil, fmt.Errorf("save articles: %w", err)
		}
	}

	return items, nil
}

func (s *Service) GetArticles(ctx context.Context, headlineID int) ([]models.Article, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("article repository is not configured")
	}
	return s.repo.GetArticlesByHeadline(ctx, headlineID)
}

func (s *Service) GetOrFetchArticles(ctx context.Context, headlineID int, topic string) ([]models.Article, error) {
	articles, err := s.GetArticles(ctx, headlineID)
	if err == nil && len(articles) > 0 {
		return articles, nil
	}
	return s.FetchAndSaveArticles(ctx, headlineID, topic)
}
