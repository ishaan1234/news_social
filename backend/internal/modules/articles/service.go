package articles

import "github.com/ishaan1234/news_social/backend/internal/models"

type Service struct {
	newsClient newsAPIClient
	repo Repository
}

func NewService(newsClient newsAPIClient, repo Repository) *Service {
	return &Service{newsClient: newsClient, repo: repo}
}

func (s *Service) GetByHeadline(headlineID int64) ([]models.Article, error) {
	return s.repo.FindByHeadline(headlineID)
}

// Fetch and persist articles for a headline/topic
func (s *Service) FetchAndSaveArticles(ctx context.Context, headlineID int, topic string) ([]models.Article, error) {
	rawArticles, err := s.newsClient.GetTopHeadlines(topic)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch articles: %w", err)
	}

	var articles []models.Article
	for _, ra := range rawArticles {
		articles = append(articles, models.Article{
			HeadlineID: headlineID,
			Source:     ra.Source,
			URL:        ra.URL,
			Content:    ra.Content,
		})
	}

	if err := s.repo.SaveArticles(ctx, headlineID, articles); err != nil {
		return nil, fmt.Errorf("failed to save articles: %w", err)
	}

	return articles, nil
}

// Retrieve articles by headline
func (s *Service) GetArticles(ctx context.Context, headlineID int) ([]models.Article, error) {
	return s.repo.GetArticlesByHeadline(ctx, headlineID)
}

func (s *Service) GetOrFetchArticles(ctx context.Context, headlineID int, topic string) ([]models.Article, error) {

	// 1. Try DB first
	existing, _ := s.repo.GetArticlesByHeadline(ctx, headlineID)
	if len(existing) > 0 {
		return existing, nil
	}

	// 2. Fetch from API
	return s.FetchAndSaveArticles(ctx, headlineID, topic)
}