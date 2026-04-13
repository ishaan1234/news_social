package summaries

import (
	"context"
	"database/sql"
	"fmt"
	"your_project/backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

// Save summary for a headline
func (r *Repository) SaveSummary(ctx context.Context, summary models.Summary) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO summaries (headline_id, summary)
		VALUES ($1, $2)
		ON CONFLICT (headline_id) DO UPDATE
		SET summary = EXCLUDED.summary
	`, summary.HeadlineID, summary.Summary)
	if err != nil {
		return fmt.Errorf("failed to save summary: %w", err)
	}
	return nil
}

// Get summary by headline ID
func (r *Repository) GetSummary(ctx context.Context, headlineID int) (models.Summary, error) {
	var summary models.Summary
	err := r.db.QueryRowContext(ctx, `
		SELECT id, headline_id, summary, created_at
		FROM summaries
		WHERE headline_id = $1
	`, headlineID).Scan(&summary.ID, &summary.HeadlineID, &summary.Summary, &summary.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return summary, nil
		}
		return summary, fmt.Errorf("failed to fetch summary: %w", err)
	}
	return summary, nil
}