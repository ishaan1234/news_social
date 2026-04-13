package summaries

import {
	"context"
	"fmt"
	"github.com/ishaan1234/news_social/backend/internal/models"
	"github.com/ishaan1234/news_social/backend/pkg/clients/openai"
}

type AIClient interface {
	Summarize(text string) (string, error)
}

type Repository interface {
	FindByHeadline(headlineID int64) (*models.Summary, error)
	Save(summary *models.Summary) error
}

type Service struct {
	repo Repository
	ai   AIClient
}

func NewService(repo Repository, ai AIClient) *Service {
	return &Service{repo: repo, ai: ai}
}

// Generate summary for a headline by fetching articles, combining text, and sending to AI
func (s *Service) GenerateSummary(headlineID string) (*models.Summary, error) {
	// Fetch articles, combine text, send to AI
	content, err := s.ai.Summarize("combined article text")
	if err != nil {
		return nil, err
	}

	summary := &models.Summary{
		Content: content,
		Model:   "gpt-4",
	}

	return summary, s.repo.Save(summary)
}

// Generate summary and save to DB
func (s *Service) GenerateAndSaveSummary(ctx context.Context, headlineID int, content string) (string, error) {
	summaryText, err := s.aiClient.GenerateSummary(content)
	if err != nil {
		return "", fmt.Errorf("failed to generate summary: %w", err)
	}

	summary := models.Summary{
		HeadlineID: headlineID,
		Summary:    summaryText,
	}

	if err := s.repo.SaveSummary(ctx, summary); err != nil {
		return "", fmt.Errorf("failed to save summary: %w", err)
	}

	return summaryText, nil
}

// Fetch existing summary
func (s *Service) GetSummary(ctx context.Context, headlineID int) (models.Summary, error) {
	return s.repo.GetSummary(ctx, headlineID)
}
