package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ishaan1234/news_social/backend/internal/config"
	"github.com/ishaan1234/news_social/backend/internal/db"
	"github.com/ishaan1234/news_social/backend/internal/middleware"
	"github.com/ishaan1234/news_social/backend/internal/modules/articles"
	"github.com/ishaan1234/news_social/backend/internal/modules/headlines"
	"github.com/ishaan1234/news_social/backend/internal/modules/social"
	"github.com/ishaan1234/news_social/backend/internal/modules/summaries"
	"github.com/ishaan1234/news_social/backend/pkg/clients/ai"
	"github.com/ishaan1234/news_social/backend/pkg/clients/newsapi"
)

type HTTPServer struct {
	addr string
	mux  *http.ServeMux
}

func NewHTTPServer(cfg *config.Config, postgres *db.Postgres) *HTTPServer {
	mux := http.NewServeMux()

	// External Clients
	newsClient := newsapi.NewClient(cfg.NewsAPIKey)
	aiClient := ai.NewClient(cfg.OpenAIAPIKey)

	// Repositories
	headlineRepo := headlines.NewRepository(postgres.DB)
	articleRepo := articles.NewRepository(postgres.DB)
	summaryRepo := summaries.NewRepository(postgres.DB)
	socialRepo := social.NewRepository(postgres.DB)

	// Services
	headlineService := headlines.NewService(headlineRepo)
	articleService := articles.NewService(newsClient, articleRepo)
	summaryService := summaries.NewService(aiClient, summaryRepo)
	socialService := social.NewService(socialRepo)

	// Handlers
	headlineHandler := headlines.NewHandler(headlineService)
	articleHandler := articles.NewHandler(articleService)
	summaryHandler := summaries.NewHandler(summaryService)
	socialHandler := social.NewHandler(socialService)

	// Aggregator
	aggregator := headlines.NewAggregator(
		headlineService,
		articleService,
		summaryService,
		socialService,
	)

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	loggingMiddleware := middleware.LoggingMiddleware
	rateLimitMiddleware := middleware.RateLimitMiddleware(cfg.RateLimitRPS)

	// Helper for chaining middleware
	chain := func(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			h = m[i](h)
		}
		return h
	}

	// Headlines
	mux.Handle("/api/headlines",
		chain(http.HandlerFunc(headlineHandler.GetHeadlines),
			loggingMiddleware,
			rateLimitMiddleware,
		),
	)

	mux.Handle("/api/headlines/create",
		chain(http.HandlerFunc(headlineHandler.CreateHeadline),
			loggingMiddleware,
			rateLimitMiddleware,
		),
	)

	// Articles (fetch or DB)
	mux.Handle("/api/articles",
		chain(http.HandlerFunc(articleHandler.GetArticles),
			loggingMiddleware,
			rateLimitMiddleware,
		),
	)

	// Full Aggregation Endpoint
	mux.Handle("/api/headline/full",
		chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			idStr := r.URL.Query().Get("id")
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(w, "invalid id", http.StatusBadRequest)
				return
			}

			result, err := aggregator.GetFullView(r.Context(), id)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}),
			loggingMiddleware,
			rateLimitMiddleware,
		),
	)

	// Generate Summary
	mux.Handle("/api/summaries",
		chain(http.HandlerFunc(summaryHandler.GenerateSummaryHandler),
			loggingMiddleware,
			rateLimitMiddleware,
			authMiddleware,
		),
	)

	// Comments
	mux.Handle("/api/comments",
		chain(http.HandlerFunc(socialHandler.CreateComment),
			loggingMiddleware,
			rateLimitMiddleware,
			authMiddleware,
		),
	)

	mux.Handle("/api/comments/list",
		chain(http.HandlerFunc(socialHandler.GetComments),
			loggingMiddleware,
			rateLimitMiddleware,
		),
	)

	return &HTTPServer{
		addr: ":" + cfg.Port,
		mux:  mux,
	}
}

func (s *HTTPServer) Start() {
	log.Printf("Server running on %s", s.addr)

	if err := http.ListenAndServe(s.addr, s.mux); err != nil {
		log.Fatal(err)
	}
}