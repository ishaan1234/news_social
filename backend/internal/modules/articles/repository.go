package articles

import "github.com/ishaan1234/news_social/backend/internal/models"

type Repository interface {
	FindByHeadline(headlineID int64) ([]models.Article, error)
	SaveBulk([]models.Article) error
}
