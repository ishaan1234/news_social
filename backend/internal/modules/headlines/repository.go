package headlines

import "github.com/ishaan1234/news_social/backend/internal/models"

type ArticleAggregator interface {
	ListHeadlines() ([]models.Headline, error)
	FetchHeadlineBundle(id string) (any, error)
}
