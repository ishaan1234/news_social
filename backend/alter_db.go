package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Error loading .env file (ignoring): %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`
		ALTER TABLE users
		ADD COLUMN IF NOT EXISTS role TEXT,
		ADD COLUMN IF NOT EXISTS bio TEXT,
		ADD COLUMN IF NOT EXISTS location TEXT,
		ADD COLUMN IF NOT EXISTS website TEXT;
	`)
	if err != nil {
		log.Fatalf("Error altering table: %v", err)
	}
	fmt.Println("Successfully altered users table!")
}
