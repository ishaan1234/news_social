package main

import (
	"log"
	"net/http"

	"github.com/ishaan1234/news_social/backend/internal/config"
)

func main() {
	loadDotEnv()
	cfg := config.Load()

	mux := http.NewServeMux()
	mux.HandleFunc("/news", newsHandler)
	registerFirebaseEmailPasswordRoutes(mux)

	addr := ":" + cfg.Port
	log.Printf("server listening on http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(addr, mux))
}
