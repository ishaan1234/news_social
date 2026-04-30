package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/ishaan1234/news_social/backend/internal/config"
	internaldb "github.com/ishaan1234/news_social/backend/internal/db"
	"github.com/ishaan1234/news_social/backend/internal/server"
)

func main() {
	loadDotEnv()
	cfg := config.Load()

	var sqlDB *sql.DB
	postgres, err := internaldb.NewPostgres(cfg.DBUrl)
	if err != nil {
		log.Printf("database unavailable; database-backed routes may be disabled: %v", err)
	} else {
		defer postgres.Close()
		sqlDB = postgres.DB

		if shouldRunInternalMigrations() {
			if err := internaldb.RunMigrations(postgres, migrationsDir()); err != nil {
				log.Fatalf("failed to run database migrations: %v", err)
			}
		}

		log.Println("registered database-backed routes")
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/news", newsHandler(sqlDB))
	mux.HandleFunc("/posts", createPostHandler(sqlDB))
	mux.HandleFunc("/following", followingHandler(sqlDB))
	mux.HandleFunc("/feed", feedHandler(sqlDB))
	mux.HandleFunc("/post-likes", postLikesHandler(sqlDB))
	mux.HandleFunc("/post-comments", postCommentsHandler(sqlDB))
	mux.HandleFunc("/profile", profileHandler(sqlDB))

	registerFirebaseEmailPasswordRoutes(mux)

	apiServer := server.NewHTTPServer(cfg, postgres)
	mux.Handle("/api/", apiServer.Handler())

	addr := ":" + cfg.Port
	log.Printf("server listening on http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(addr, mux))
}

func shouldRunInternalMigrations() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("RUN_INTERNAL_MIGRATIONS")), "true")
}

func migrationsDir() string {
	candidates := []string{
		filepath.Join("internal", "db", "migrations"),
		filepath.Join("backend", "internal", "db", "migrations"),
	}

	for _, candidate := range candidates {
		info, err := os.Stat(candidate)
		if err == nil && info.IsDir() {
			return candidate
		}
	}

	return filepath.Join("internal", "db", "migrations")
}
