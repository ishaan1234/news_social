package summaries

import (
	"context"
	"fmt"
	"strings"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type AIClient interface {
	GenerateSummary(content string) (string, error)
}

type Service struct {
	aiClient AIClient
	repo     Repository
	model    string
}

func NewService(aiClient AIClient, repo Repository) *Service {
	return &Service{aiClient: aiClient, repo: repo, model: "gpt-4o-mini"}
}

func (s *Service) GenerateAndSaveSummary(ctx context.Context, headlineID int, content string) (models.Summary, error) {
	if headlineID <= 0 {
		return models.Summary{}, fmt.Errorf("valid headline id is required")
	}
	if strings.TrimSpace(content) == "" {
		return models.Summary{}, fmt.Errorf("content is required")
	}
	if s.aiClient == nil {
		return models.Summary{}, fmt.Errorf("ai client is not configured")
	}

	text, err := s.aiClient.GenerateSummary(content)
	if err != nil {
		return models.Summary{}, fmt.Errorf("generate summary: %w", err)
	}

	summary := models.Summary{
		HeadlineID: headlineID,
		Content:    text,
		Model:      s.model,
	}

	if s.repo != nil {
		if err := s.repo.SaveSummary(ctx, summary); err != nil {
			return models.Summary{}, err
		}
	}

	return summary, nil
}

func (s *Service) GetSummary(ctx context.Context, headlineID int) (models.Summary, error) {
	if headlineID <= 0 {
		return models.Summary{}, fmt.Errorf("valid headline id is required")
	}
	if s.repo == nil {
		return models.Summary{}, fmt.Errorf("summary repository is not configured")
	}
	return s.repo.GetSummary(ctx, headlineID)
}

func (s *Service) GetOrGenerateSummary(ctx context.Context, headlineID int, content string) (models.Summary, error) {
	existing, err := s.GetSummary(ctx, headlineID)
	if err == nil && existing.Content != "" {
		return existing, nil
	}
	return s.GenerateAndSaveSummary(ctx, headlineID, content)
}
