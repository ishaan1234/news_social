package main

import (
    "fmt"
    "net/http"
)

func main() {
    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        fmt.Fprint(w, `{"ok":true}`)
    })

    fmt.Println("Listening on :8080")
    _ = http.ListenAndServe(":8080", nil)
}
