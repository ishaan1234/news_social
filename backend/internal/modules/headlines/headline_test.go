package headlines

import (
	"context"
	"errors"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type mockHeadlineRepo struct {
	createdTitle string
	err          error
}

func (m *mockHeadlineRepo) Create(ctx context.Context, title string) (int, error) {
	m.createdTitle = title
	return 7, m.err
}

func (m *mockHeadlineRepo) GetAll(ctx context.Context) ([]models.Headline, error) {
	return []models.Headline{{ID: 1, Title: "One"}}, m.err
}

func (m *mockHeadlineRepo) GetByID(ctx context.Context, id int) (models.Headline, error) {
	if m.err != nil {
		return models.Headline{}, m.err
	}
	return models.Headline{ID: id, Title: "One"}, nil
}

func TestService_CreateHeadline(t *testing.T) {
	repo := &mockHeadlineRepo{}
	service := NewService(repo)

	id, err := service.CreateHeadline(context.Background(), "  Breaking News  ")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id != 7 || repo.createdTitle != "Breaking News" {
		t.Fatalf("headline was not created correctly")
	}
}

func TestService_CreateHeadline_EmptyTitle(t *testing.T) {
	service := NewService(&mockHeadlineRepo{})

	_, err := service.CreateHeadline(context.Background(), "   ")
	if err == nil {
		t.Fatalf("expected validation error")
	}
}

func TestService_GetHeadline_Error(t *testing.T) {
	service := NewService(&mockHeadlineRepo{err: errors.New("not found")})

	_, err := service.GetHeadline(context.Background(), 1)
	if err == nil {
		t.Fatalf("expected error")
	}
}
