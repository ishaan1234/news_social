// package main

// import (
// 	"log"
// 	"net/http"

// 	"github.com/ishaan1234/news_social/backend/internal/config"
// )

// func main() {
// 	loadDotEnv()
// 	cfg := config.Load()

// 	mux := http.NewServeMux()
// 	mux.HandleFunc("/news", newsHandler)
// 	registerFirebaseEmailPasswordRoutes(mux)

// 	addr := ":" + cfg.Port
// 	log.Printf("server listening on http://localhost:%s", cfg.Port)
// 	log.Fatal(http.ListenAndServe(addr, mux))
// }

package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/ishaan1234/news_social/backend/internal/config"

	_ "github.com/lib/pq"
)

func main() {
	loadDotEnv()
	cfg := config.Load()

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Fatal("failed to open database:", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("failed to connect to database:", err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/news", newsHandler(db))
	mux.HandleFunc("/posts", createPostHandler(db))
	mux.HandleFunc("/following", followingHandler(db))
	mux.HandleFunc("/feed", feedHandler(db))
	mux.HandleFunc("/post-likes", postLikesHandler(db))
	mux.HandleFunc("/post-comments", postCommentsHandler(db))

	registerFirebaseEmailPasswordRoutes(mux)

	addr := ":" + cfg.Port
	log.Printf("server listening on http://localhost:%s", cfg.Port)
	log.Fatal(http.ListenAndServe(addr, mux))
}
