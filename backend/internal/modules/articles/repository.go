package articles

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Repository interface {
	SaveArticles(ctx context.Context, headlineID int, articles []models.Article) error
	GetArticlesByHeadline(ctx context.Context, headlineID int) ([]models.Article, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) SaveArticles(ctx context.Context, headlineID int, articles []models.Article) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("article repository database is not configured")
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO articles (headline_id, source, title, url, content, published_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (url) DO UPDATE
		SET source = EXCLUDED.source,
		    title = EXCLUDED.title,
		    content = EXCLUDED.content,
		    published_at = EXCLUDED.published_at
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, a := range articles {
		if _, err := stmt.ExecContext(ctx, headlineID, a.Source, a.Title, a.URL, a.Content, a.PublishedAt); err != nil {
			tx.Rollback()
			return fmt.Errorf("insert article: %w", err)
		}
	}

	return tx.Commit()
}

func (r *SQLRepository) GetArticlesByHeadline(ctx context.Context, headlineID int) ([]models.Article, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("article repository database is not configured")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, headline_id, source, title, url, content, COALESCE(published_at, created_at), created_at
		FROM articles
		WHERE headline_id = $1
		ORDER BY created_at DESC
	`, headlineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Article
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(&a.ID, &a.HeadlineID, &a.Source, &a.Title, &a.URL, &a.Content, &a.PublishedAt, &a.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, a)
	}

	return result, rows.Err()
}
