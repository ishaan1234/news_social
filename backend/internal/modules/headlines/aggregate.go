package headlines

import (
	"context"
	"github.com/ishaan1234/news_social/backend/internal/modules/articles"
	"github.com/ishaan1234/news_social/backend/internal/modules/summaries"
	"github.com/ishaan1234/news_social/backend/internal/modules/social"
)

type Aggregator struct {
	headlineSvc *Service
	articleSvc  *articles.Service
	summarySvc  *summaries.Service
	socialSvc   *social.Service
}

func NewAggregator(h *Service, a *articles.Service, s *summaries.Service, so *social.Service) *Aggregator {
	return &Aggregator{
		headlineSvc: h,
		articleSvc:  a,
		summarySvc:  s,
		socialSvc:   so,
	}
}

func (ag *Aggregator) GetFullView(ctx context.Context, headlineID int) (map[string]interface{}, error) {
	headline, _ := ag.headlineSvc.GetHeadline(ctx, headlineID)
	articles, _ := ag.articleSvc.GetArticles(ctx, headlineID)
	summary, _ := ag.summarySvc.GetSummary(ctx, headlineID)
	comments, _ := ag.socialSvc.GetComments(ctx, headlineID)

	return map[string]interface{}{
		"headline": headline,
		"articles": articles,
		"summary":  summary,
		"comments": comments,
	}, nil
}