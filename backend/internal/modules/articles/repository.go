package articles

import {
	"context"
	"database/sql"
	"fmt"
	"github.com/ishaan1234/news_social/backend/internal/models"}

type Repository interface {
	FindByHeadline(headlineID int64) ([]models.Article, error)
	SaveBulk([]models.Article) error
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save multiple articles linked to a headline
func (r *Repository) SaveArticles(ctx context.Context, headlineID int, articles []models.Article) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	stmt, err := tx.PrepareContext(ctx, `
		INSERT INTO articles (headline_id, source, url, content)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (url) DO NOTHING
	`)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmt.Close()

	for _, a := range articles {
		if _, err := stmt.ExecContext(ctx, headlineID, a.Source, a.URL, a.Content); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert article: %w", err)
		}
	}

	return tx.Commit()
}

// Fetch articles by headline ID
func (r *Repository) GetArticlesByHeadline(ctx context.Context, headlineID int) ([]models.Article, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, headline_id, source, url, content, created_at
		FROM articles
		WHERE headline_id = $1
	`, headlineID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Article
	for rows.Next() {
		var a models.Article
		if err := rows.Scan(&a.ID, &a.HeadlineID, &a.Source, &a.URL, &a.Content, &a.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, a)
	}
	return result, nil
}