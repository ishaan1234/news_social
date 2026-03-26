package main

import (
	"log"
	"net/http"
)

func main() {
	loadDotEnv()

	mux := http.NewServeMux()
	mux.HandleFunc("/news", newsHandler)

	log.Println("server listening on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
