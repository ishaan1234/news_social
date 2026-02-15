package summaries

import "social-news/internal/models"

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

