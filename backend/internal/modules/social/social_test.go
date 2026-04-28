package social

import (
	"context"
	"errors"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type mockSocialRepo struct {
	comment models.Comment
	err     error
}

func (m *mockSocialRepo) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	comment.ID = 11
	m.comment = comment
	return comment, m.err
}

func (m *mockSocialRepo) GetCommentsByHeadline(ctx context.Context, headlineID int) ([]models.Comment, error) {
	return []models.Comment{{ID: 1, HeadlineID: headlineID, Content: "hello"}}, m.err
}

func (m *mockSocialRepo) GetCommentsByUser(ctx context.Context, userID int) ([]models.Comment, error) {
	return nil, m.err
}

func (m *mockSocialRepo) DeleteComment(ctx context.Context, commentID int, userID int) error {
	return m.err
}

func TestService_CreateComment(t *testing.T) {
	repo := &mockSocialRepo{}
	service := NewService(repo)

	comment, err := service.CreateComment(context.Background(), CreateCommentRequest{UserID: 1, HeadlineID: 2, Content: "  Nice article  "})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 11 || comment.Content != "Nice article" {
		t.Fatalf("unexpected comment: %#v", comment)
	}
}

func TestService_CreateComment_Validation(t *testing.T) {
	service := NewService(&mockSocialRepo{})

	_, err := service.CreateComment(context.Background(), CreateCommentRequest{UserID: 0, HeadlineID: 2, Content: "x"})
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestService_CreateComment_RepoError(t *testing.T) {
	service := NewService(&mockSocialRepo{err: errors.New("db error")})

	_, err := service.CreateComment(context.Background(), CreateCommentRequest{UserID: 1, HeadlineID: 2, Content: "x"})
	if err == nil {
		t.Fatalf("expected repository error")
	}
}
