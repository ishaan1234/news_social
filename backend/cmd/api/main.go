package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ishaan1234/news_social/backend/internal/config"
	"github.com/ishaan1234/news_social/backend/internal/db"
	"github.com/ishaan1234/news_social/backend/internal/server"
)

func main() {
	loadDotEnv()
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/news", newsHandler)
	registerFirebaseEmailPasswordRoutes(mux)

	postgres, err := db.NewPostgres(cfg.DBUrl)
	if err != nil {
		log.Printf("database unavailable; persistent /api routes disabled, posts API using memory: %v", err)
	} else {
		defer postgres.Close()

		if err := db.RunMigrations(postgres, migrationsDir()); err != nil {
			log.Fatalf("failed to run database migrations: %v", err)
		}

		log.Println("registered database-backed /api routes")
	}

	apiServer := server.NewHTTPServer(cfg, postgres)
	mux.Handle("/api/", apiServer.Handler())

	addr := ":" + cfg.Port
	log.Printf("server listening on http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(addr, mux))
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
