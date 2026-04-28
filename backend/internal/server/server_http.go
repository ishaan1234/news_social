package server

import (
	"database/sql"
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
	"github.com/ishaan1234/news_social/backend/internal/utils"
	"github.com/ishaan1234/news_social/backend/pkg/clients/ai"
	"github.com/ishaan1234/news_social/backend/pkg/clients/newsapi"
)

type HTTPServer struct {
	addr string
	mux  *http.ServeMux
}

func NewHTTPServer(cfg *config.Config, postgres *db.Postgres) *HTTPServer {
	mux := http.NewServeMux()

	var sqlDB *sql.DB
	if postgres != nil {
		sqlDB = postgres.DB
	}

	newsClient := newsapi.NewClient(cfg.NewsAPIKey)
	aiClient := ai.NewClient(cfg.OpenAIAPIKey)

	headlineRepo := headlines.NewRepository(sqlDB)
	articleRepo := articles.NewRepository(sqlDB)
	summaryRepo := summaries.NewRepository(sqlDB)
	socialRepo := social.NewRepository(sqlDB)

	headlineService := headlines.NewService(headlineRepo)
	articleService := articles.NewService(newsClient, articleRepo)
	summaryService := summaries.NewService(aiClient, summaryRepo)
	socialService := social.NewService(socialRepo)

	headlineHandler := headlines.NewHandler(headlineService)
	articleHandler := articles.NewHandler(articleService)
	summaryHandler := summaries.NewHandler(summaryService)
	socialHandler := social.NewHandler(socialService)

	aggregator := headlines.NewAggregator(headlineService, articleService, summaryService, socialService)

	authMiddleware := middleware.AuthMiddleware(cfg.JWTSecret)
	loggingMiddleware := middleware.LoggingMiddleware
	rateLimitMiddleware := middleware.RateLimitMiddleware(cfg.RateLimitRPS)

	chain := func(h http.Handler, m ...func(http.Handler) http.Handler) http.Handler {
		for i := len(m) - 1; i >= 0; i-- {
			h = m[i](h)
		}
		return h
	}

	mux.Handle("/api/headlines",
		chain(http.HandlerFunc(headlineHandler.GetHeadlines), loggingMiddleware, rateLimitMiddleware),
	)

	mux.Handle("/api/headlines/create",
		chain(http.HandlerFunc(headlineHandler.CreateHeadline), loggingMiddleware, rateLimitMiddleware),
	)

	mux.Handle("/api/articles",
		chain(http.HandlerFunc(articleHandler.GetArticles), loggingMiddleware, rateLimitMiddleware),
	)

	mux.Handle("/api/headline/full",
		chain(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil || id <= 0 {
				utils.WriteError(w, http.StatusBadRequest, "valid id is required")
				return
			}

			result, err := aggregator.GetFullView(r.Context(), id)
			if err != nil {
				utils.WriteError(w, http.StatusInternalServerError, err.Error())
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(result)
		}), loggingMiddleware, rateLimitMiddleware),
	)

	mux.Handle("/api/summaries",
		chain(http.HandlerFunc(summaryHandler.GenerateSummaryHandler), loggingMiddleware, rateLimitMiddleware, authMiddleware),
	)

	mux.Handle("/api/comments",
		chain(http.HandlerFunc(socialHandler.CreateComment), loggingMiddleware, rateLimitMiddleware, authMiddleware),
	)

	mux.Handle("/api/comments/list",
		chain(http.HandlerFunc(socialHandler.GetComments), loggingMiddleware, rateLimitMiddleware),
	)

	return &HTTPServer{addr: ":" + cfg.Port, mux: mux}
}

func (s *HTTPServer) Start() {
	log.Printf("server running on %s", s.addr)
	if err := http.ListenAndServe(s.addr, s.mux); err != nil {
		log.Fatal(err)
	}
}
