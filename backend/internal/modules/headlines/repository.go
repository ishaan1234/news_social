package headlines

import "social-news/internal/models"

type ArticleAggregator interface {
	ListHeadlines() ([]models.Headline, error)
	FetchHeadlineBundle(id string) (any, error)
}
