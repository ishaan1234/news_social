package posts

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

func (s *Service) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	if s.repo == nil {
		return models.Post{}, fmt.Errorf("post repository is not configured")
	}

	post.Body = strings.TrimSpace(post.Body)
	if post.Body == "" {
		return models.Post{}, fmt.Errorf("post body is required")
	}

	post.Article.URL = strings.TrimSpace(post.Article.URL)
	post.Article.Title = strings.TrimSpace(post.Article.Title)
	post.Article.Summary = strings.TrimSpace(post.Article.Summary)
	if post.Article.URL == "" {
		return models.Post{}, fmt.Errorf("article url is required")
	}
	if post.Article.Title == "" {
		return models.Post{}, fmt.Errorf("article title is required")
	}

	post.AuthorID = strings.TrimSpace(post.AuthorID)
	post.AuthorName = strings.TrimSpace(post.AuthorName)
	if post.AuthorName == "" {
		post.AuthorName = "Anonymous"
	}
	post.AuthorHandle = strings.TrimSpace(post.AuthorHandle)
	post.Article.Source = strings.TrimSpace(post.Article.Source)
	post.Article.ImageURL = strings.TrimSpace(post.Article.ImageURL)
	post.Article.PublishedAt = strings.TrimSpace(post.Article.PublishedAt)

	return s.repo.CreatePost(ctx, post)
}

func (s *Service) GetPosts(ctx context.Context, viewerID string) ([]models.Post, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("post repository is not configured")
	}
	return s.repo.GetPosts(ctx, strings.TrimSpace(viewerID))
}

func (s *Service) GetPost(ctx context.Context, postID int, viewerID string) (models.Post, error) {
	if s.repo == nil {
		return models.Post{}, fmt.Errorf("post repository is not configured")
	}
	if postID <= 0 {
		return models.Post{}, fmt.Errorf("valid post_id is required")
	}
	return s.repo.GetPostByID(ctx, postID, strings.TrimSpace(viewerID))
}

func (s *Service) CreateComment(ctx context.Context, comment models.PostComment) (models.PostComment, error) {
	if s.repo == nil {
		return models.PostComment{}, fmt.Errorf("post repository is not configured")
	}
	if comment.PostID <= 0 {
		return models.PostComment{}, fmt.Errorf("valid post_id is required")
	}

	comment.Content = strings.TrimSpace(comment.Content)
	if comment.Content == "" {
		return models.PostComment{}, fmt.Errorf("comment content is required")
	}

	comment.AuthorID = strings.TrimSpace(comment.AuthorID)
	comment.AuthorName = strings.TrimSpace(comment.AuthorName)
	if comment.AuthorName == "" {
		comment.AuthorName = "Anonymous"
	}

	return s.repo.CreateComment(ctx, comment)
}

func (s *Service) GetComments(ctx context.Context, postID int) ([]models.PostComment, error) {
	if s.repo == nil {
		return nil, fmt.Errorf("post repository is not configured")
	}
	if postID <= 0 {
		return nil, fmt.Errorf("valid post_id is required")
	}
	return s.repo.GetComments(ctx, postID)
}

func (s *Service) SetVote(ctx context.Context, postID int, voterID string, value int) (models.PostVoteSummary, error) {
	if s.repo == nil {
		return models.PostVoteSummary{}, fmt.Errorf("post repository is not configured")
	}
	if postID <= 0 {
		return models.PostVoteSummary{}, fmt.Errorf("valid post_id is required")
	}

	voterID = strings.TrimSpace(voterID)
	if voterID == "" {
		return models.PostVoteSummary{}, fmt.Errorf("voter_id is required")
	}
	if value < -1 || value > 1 {
		return models.PostVoteSummary{}, fmt.Errorf("vote value must be -1, 0, or 1")
	}

	return s.repo.SetVote(ctx, postID, voterID, value)
}

func (s *Service) SharePost(ctx context.Context, postID int) (int, error) {
	if s.repo == nil {
		return 0, fmt.Errorf("post repository is not configured")
	}
	if postID <= 0 {
		return 0, fmt.Errorf("valid post_id is required")
	}
	return s.repo.IncrementShare(ctx, postID)
}
