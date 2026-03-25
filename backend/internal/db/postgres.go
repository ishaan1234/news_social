package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

type PostgresDB struct {
	Conn *sql.DB
}

func NewPostgres(dbURL string) *PostgresDB {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	// Optional: set max connections
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Ping to verify connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	return &PostgresDB{
		Conn: db,
	}
}

// Close the database connection
func (p *PostgresDB) Close() error {
	return p.Conn.Close()
}

// Helper function to execute a query with context
func (p *PostgresDB) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return p.Conn.ExecContext(ctx, query, args...)
}

// Helper function to query rows with context
func (p *PostgresDB) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return p.Conn.QueryContext(ctx, query, args...)
}

// Helper function to query single row
func (p *PostgresDB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	return p.Conn.QueryRowContext(ctx, query, args...)
}

// Example migration runner: executes all SQL files in migrations folder
func (p *PostgresDB) RunMigrations(migrationFiles []string) error {
	for _, file := range migrationFiles {
		sqlStmt, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		if _, err := p.Conn.Exec(string(sqlStmt)); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
		log.Printf("migration applied: %s", file)
	}
	return nil
}