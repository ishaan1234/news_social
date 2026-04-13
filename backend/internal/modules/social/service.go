package social

import "github.com/ishaan1234/news_social/backend/internal/models"

type Repository interface {
	Create(comment *models.Comment) error
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateCommentRequest struct {
	UserID     string `json:"user_id"`
	HeadlineID string `json:"headline_id"`
	Content    string `json:"content"`
}

func (s *Service) CreateComment(req CreateCommentRequest) (*models.Comment, error) {
	comment := &models.Comment{
		Content: req.Content,
	}
	return comment, s.repo.Create(comment)
}
