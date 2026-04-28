package articles

import (
	"context"
	"errors"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/models"
	"github.com/ishaan1234/news_social/backend/pkg/clients/newsapi"
)

type mockArticleRepo struct {
	items []models.Article
	err   error
	saved []models.Article
}

func (m *mockArticleRepo) SaveArticles(ctx context.Context, headlineID int, articles []models.Article) error {
	m.saved = articles
	return m.err
}

func (m *mockArticleRepo) GetArticlesByHeadline(ctx context.Context, headlineID int) ([]models.Article, error) {
	return m.items, m.err
}

type mockNewsClient struct {
	items []newsapi.Article
	err   error
}

func (m *mockNewsClient) GetTopHeadlines(topic string) ([]newsapi.Article, error) {
	return m.items, m.err
}

func TestService_GetArticles(t *testing.T) {
	repo := &mockArticleRepo{items: []models.Article{{ID: 1, HeadlineID: 10, Title: "Saved"}}}
	service := NewService(nil, repo)

	items, err := service.GetArticles(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || items[0].Title != "Saved" {
		t.Fatalf("unexpected articles: %#v", items)
	}
}

func TestService_FetchAndSaveArticles(t *testing.T) {
	repo := &mockArticleRepo{}
	client := &mockNewsClient{items: []newsapi.Article{{Source: "AP", Title: "News", URL: "https://example.com", Content: "content"}}}
	service := NewService(client, repo)

	items, err := service.FetchAndSaveArticles(context.Background(), 3, "topic")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 1 || len(repo.saved) != 1 {
		t.Fatalf("expected one fetched and saved article")
	}
	if repo.saved[0].HeadlineID != 3 {
		t.Fatalf("headline id was not propagated")
	}
}

func TestService_FetchAndSaveArticles_ClientError(t *testing.T) {
	service := NewService(&mockNewsClient{err: errors.New("news api down")}, &mockArticleRepo{})

	_, err := service.FetchAndSaveArticles(context.Background(), 3, "topic")
	if err == nil {
		t.Fatalf("expected error")
	}
}
