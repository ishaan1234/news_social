package articles

import "social-news/internal/models"

type Repository interface {
	FindByHeadline(headlineID int64) ([]models.Article, error)
	SaveBulk([]models.Article) error
}
