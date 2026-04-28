package summaries

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Repository interface {
	SaveSummary(ctx context.Context, summary models.Summary) error
	GetSummary(ctx context.Context, headlineID int) (models.Summary, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) SaveSummary(ctx context.Context, summary models.Summary) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("summary repository database is not configured")
	}

	_, err := r.db.ExecContext(ctx, `
		INSERT INTO summaries (headline_id, content, model)
		VALUES ($1, $2, $3)
		ON CONFLICT (headline_id) DO UPDATE
		SET content = EXCLUDED.content,
		    model = EXCLUDED.model
	`, summary.HeadlineID, summary.Content, summary.Model)
	if err != nil {
		return fmt.Errorf("save summary: %w", err)
	}
	return nil
}

func (r *SQLRepository) GetSummary(ctx context.Context, headlineID int) (models.Summary, error) {
	if r == nil || r.db == nil {
		return models.Summary{}, fmt.Errorf("summary repository database is not configured")
	}

	var summary models.Summary
	err := r.db.QueryRowContext(ctx, `
		SELECT id, headline_id, content, model, created_at
		FROM summaries
		WHERE headline_id = $1
	`, headlineID).Scan(&summary.ID, &summary.HeadlineID, &summary.Content, &summary.Model, &summary.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return models.Summary{}, nil
		}
		return models.Summary{}, fmt.Errorf("get summary: %w", err)
	}
	return summary, nil
}
