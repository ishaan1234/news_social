package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/ishaan1234/news_social/backend/internal/models"
	// "github.com/ishaan1234/news_social/backend/internal/repository"
)

type socialService struct {
	// commentRepo repository.CommentRepository
	// voteRepo    repository.VoteRepository
}

// func NewSocialService(
// 	c repository.CommentRepository,
// 	v repository.VoteRepository,
// ) SocialService {
// 	return &socialService{
// 		commentRepo: c,
// 		voteRepo:    v,
// 	}
// }

func (s *socialService) AddComment(ctx context.Context, userID, headlineID uuid.UUID, content string) error {
	_ = &models.Comment{
		ID:         uuid.New(),
		UserID:     userID,
		HeadlineID: headlineID,
		Content:    content,
	}
	// return s.commentRepo.Create(ctx, comment)
	return nil
}

func (s *socialService) Vote(ctx context.Context, userID, headlineID uuid.UUID, value int) error {
	if value != 1 && value != -1 {
		return errors.New("invalid vote value")
	}

	// return s.voteRepo.Upsert(ctx, &models.Vote{
	// 	UserID:     userID,
	// 	HeadlineID: headlineID,
	// 	Value:      value,
	// })
	return nil
}

func (s *socialService) GetComments(ctx context.Context, headlineID uuid.UUID) ([]models.Comment, error) {
	// return s.commentRepo.GetByHeadlineID(ctx, headlineID)
	return nil, nil
}
