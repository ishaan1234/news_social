package headlines

import (
	"context"
	"fmt"
	"strings"

	"github.com/ishaan1234/news_social/backend/internal/modules/articles"
	"github.com/ishaan1234/news_social/backend/internal/modules/social"
	"github.com/ishaan1234/news_social/backend/internal/modules/summaries"
)

type Aggregator struct {
	headlineSvc *Service
	articleSvc  *articles.Service
	summarySvc  *summaries.Service
	socialSvc   *social.Service
}

func NewAggregator(h *Service, a *articles.Service, s *summaries.Service, so *social.Service) *Aggregator {
	return &Aggregator{headlineSvc: h, articleSvc: a, summarySvc: s, socialSvc: so}
}

func (ag *Aggregator) GetFullView(ctx context.Context, headlineID int) (map[string]interface{}, error) {
	headline, err := ag.headlineSvc.GetHeadline(ctx, headlineID)
	if err != nil {
		return nil, fmt.Errorf("get headline: %w", err)
	}

	articleItems, err := ag.articleSvc.GetOrFetchArticles(ctx, headlineID, headline.Title)
	if err != nil {
		return nil, fmt.Errorf("get articles: %w", err)
	}

	var combined []string
	for _, article := range articleItems {
		if article.Content != "" {
			combined = append(combined, article.Content)
		}
	}

	summary, err := ag.summarySvc.GetOrGenerateSummary(ctx, headlineID, strings.Join(combined, "\n\n"))
	if err != nil {
		return nil, fmt.Errorf("get summary: %w", err)
	}

	comments, err := ag.socialSvc.GetComments(ctx, headlineID)
	if err != nil {
		return nil, fmt.Errorf("get comments: %w", err)
	}

	return map[string]interface{}{
		"headline": headline,
		"articles": articleItems,
		"summary":  summary,
		"comments": comments,
	}, nil
}
