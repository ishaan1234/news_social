package social

import (
	"context"
	"fmt"
	"strings"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

type CreateCommentRequest struct {
	UserID     int    `json:"user_id"`
	HeadlineID int    `json:"headline_id"`
	Content    string `json:"content"`
}

func (s *Service) CreateComment(ctx context.Context, req CreateCommentRequest) (models.Comment, error) {
	if req.UserID <= 0 {
		return models.Comment{}, fmt.Errorf("valid user_id is required")
	}
	if req.HeadlineID <= 0 {
		return models.Comment{}, fmt.Errorf("valid headline_id is required")
	}
	if strings.TrimSpace(req.Content) == "" {
		return models.Comment{}, fmt.Errorf("comment content is required")
	}

	comment := models.Comment{
		UserID:     req.UserID,
		HeadlineID: req.HeadlineID,
		Content:    strings.TrimSpace(req.Content),
	}
	return s.repo.CreateComment(ctx, comment)
}

func (s *Service) GetComments(ctx context.Context, headlineID int) ([]models.Comment, error) {
	if headlineID <= 0 {
		return nil, fmt.Errorf("valid headline_id is required")
	}
	return s.repo.GetCommentsByHeadline(ctx, headlineID)
}

func (s *Service) DeleteComment(ctx context.Context, commentID int, userID int) error {
	if commentID <= 0 || userID <= 0 {
		return fmt.Errorf("valid comment_id and user_id are required")
	}
	return s.repo.DeleteComment(ctx, commentID, userID)
}
