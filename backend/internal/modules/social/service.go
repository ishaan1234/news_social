package social

import "social-news/internal/models"

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
	UserID     int64  `json:"user_id"`
	HeadlineID int64  `json:"headline_id"`
	Content    string `json:"content"`
}

func (s *Service) CreateComment(req CreateCommentRequest) (*models.Comment, error) {
	comment := &models.Comment{
		UserID:     req.UserID,
		HeadlineID: req.HeadlineID,
		Content:    req.Content,
	}
	return comment, s.repo.Create(comment)
}
