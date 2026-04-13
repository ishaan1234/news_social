package headlines

import (
	"context"
	"database/sql"
	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, title string) (int, error) {
	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO headlines (title) VALUES ($1) RETURNING id`,
		title,
	).Scan(&id)
	return id, err
}

func (r *Repository) GetAll(ctx context.Context) ([]models.Headline, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, created_at FROM headlines ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Headline
	for rows.Next() {
		var h models.Headline
		rows.Scan(&h.ID, &h.Title, &h.CreatedAt)
		result = append(result, h)
	}
	return result, nil
}

func (r *Repository) GetByID(ctx context.Context, id int) (models.Headline, error) {
	var h models.Headline
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, created_at FROM headlines WHERE id=$1`, id).
		Scan(&h.ID, &h.Title, &h.CreatedAt)
	return h, err
}