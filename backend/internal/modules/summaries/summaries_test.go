package summaries

import (
	"context"
	"errors"
	"testing"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type mockAIClient struct {
	text string
	err  error
}

func (m *mockAIClient) GenerateSummary(content string) (string, error) {
	return m.text, m.err
}

type mockSummaryRepo struct {
	summary models.Summary
	err     error
	saved   models.Summary
}

func (m *mockSummaryRepo) SaveSummary(ctx context.Context, summary models.Summary) error {
	m.saved = summary
	return m.err
}

func (m *mockSummaryRepo) GetSummary(ctx context.Context, headlineID int) (models.Summary, error) {
	return m.summary, m.err
}

func TestService_GenerateAndSaveSummary(t *testing.T) {
	repo := &mockSummaryRepo{}
	service := NewService(&mockAIClient{text: "AI generated summary"}, repo)

	summary, err := service.GenerateAndSaveSummary(context.Background(), 9, "article text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Content != "AI generated summary" || repo.saved.HeadlineID != 9 {
		t.Fatalf("summary was not generated and saved correctly")
	}
}

func TestService_GenerateAndSaveSummary_AIError(t *testing.T) {
	service := NewService(&mockAIClient{err: errors.New("ai failure")}, &mockSummaryRepo{})

	_, err := service.GenerateAndSaveSummary(context.Background(), 9, "article text")
	if err == nil {
		t.Fatalf("expected error")
	}
}

func TestService_GetOrGenerateSummary_UsesExisting(t *testing.T) {
	repo := &mockSummaryRepo{summary: models.Summary{HeadlineID: 9, Content: "cached"}}
	service := NewService(&mockAIClient{text: "new"}, repo)

	summary, err := service.GetOrGenerateSummary(context.Background(), 9, "article text")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if summary.Content != "cached" {
		t.Fatalf("expected cached summary")
	}
}
