package posts

import (
	"context"
	"errors"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type mockPostRepo struct {
	post    models.Post
	comment models.PostComment
	vote    models.PostVoteSummary
	err     error
}

func (m *mockPostRepo) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	post.ID = 22
	m.post = post
	return post, m.err
}

func (m *mockPostRepo) GetPosts(ctx context.Context, viewerID string) ([]models.Post, error) {
	return []models.Post{m.post}, m.err
}

func (m *mockPostRepo) GetPostByID(ctx context.Context, postID int, viewerID string) (models.Post, error) {
	m.post.ID = postID
	return m.post, m.err
}

func (m *mockPostRepo) CreateComment(ctx context.Context, comment models.PostComment) (models.PostComment, error) {
	comment.ID = 3
	m.comment = comment
	return comment, m.err
}

func (m *mockPostRepo) GetComments(ctx context.Context, postID int) ([]models.PostComment, error) {
	return []models.PostComment{{ID: 1, PostID: postID, Content: "hello"}}, m.err
}

func (m *mockPostRepo) SetVote(ctx context.Context, postID int, voterID string, value int) (models.PostVoteSummary, error) {
	m.vote = models.PostVoteSummary{PostID: postID, VoteScore: value, ViewerVote: value}
	return m.vote, m.err
}

func (m *mockPostRepo) IncrementShare(ctx context.Context, postID int) (int, error) {
	return 4, m.err
}

func TestService_CreatePost(t *testing.T) {
	repo := &mockPostRepo{}
	service := NewService(repo)

	post, err := service.CreatePost(context.Background(), models.Post{
		AuthorName: "  Avery  ",
		Body:       "  This matters  ",
		Article: models.LinkedArticle{
			URL:   " https://example.com/story ",
			Title: " Story title ",
		},
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if post.ID != 22 || repo.post.Body != "This matters" || repo.post.Article.URL != "https://example.com/story" {
		t.Fatalf("post was not normalized and created: %#v", repo.post)
	}
}

func TestService_CreatePostRequiresArticle(t *testing.T) {
	service := NewService(&mockPostRepo{})

	_, err := service.CreatePost(context.Background(), models.Post{Body: "take"})
	if err == nil {
		t.Fatalf("expected article validation error")
	}
}

func TestService_CreateComment(t *testing.T) {
	repo := &mockPostRepo{}
	service := NewService(repo)

	comment, err := service.CreateComment(context.Background(), models.PostComment{
		PostID:  10,
		Content: "  agreed  ",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if comment.ID != 3 || repo.comment.AuthorName != "Anonymous" || repo.comment.Content != "agreed" {
		t.Fatalf("comment was not normalized and created: %#v", repo.comment)
	}
}

func TestService_SetVote(t *testing.T) {
	service := NewService(&mockPostRepo{})

	summary, err := service.SetVote(context.Background(), 10, "viewer-1", 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.PostID != 10 || summary.ViewerVote != 1 {
		t.Fatalf("unexpected vote summary: %#v", summary)
	}
}

func TestService_SetVoteValidation(t *testing.T) {
	service := NewService(&mockPostRepo{})

	_, err := service.SetVote(context.Background(), 10, "", 1)
	if err == nil {
		t.Fatalf("expected voter validation error")
	}

	_, err = service.SetVote(context.Background(), 10, "viewer-1", 2)
	if err == nil {
		t.Fatalf("expected value validation error")
	}
}

func TestService_RepoError(t *testing.T) {
	service := NewService(&mockPostRepo{err: errors.New("db down")})

	_, err := service.SharePost(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected repository error")
	}
}
