package headlines

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Repository interface {
	Create(ctx context.Context, title string) (int, error)
	GetAll(ctx context.Context) ([]models.Headline, error)
	GetByID(ctx context.Context, id int) (models.Headline, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) Create(ctx context.Context, title string) (int, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("headline repository database is not configured")
	}

	var id int
	err := r.db.QueryRowContext(ctx,
		`INSERT INTO headlines (title, slug) VALUES ($1, $2) RETURNING id`,
		title,
		slugify(title),
	).Scan(&id)
	return id, err
}

func (r *SQLRepository) GetAll(ctx context.Context) ([]models.Headline, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("headline repository database is not configured")
	}

	rows, err := r.db.QueryContext(ctx,
		`SELECT id, title, COALESCE(slug, ''), created_at FROM headlines ORDER BY created_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Headline
	for rows.Next() {
		var h models.Headline
		if err := rows.Scan(&h.ID, &h.Title, &h.Slug, &h.CreatedAt); err != nil {
			return nil, err
		}
		result = append(result, h)
	}
	return result, rows.Err()
}

func (r *SQLRepository) GetByID(ctx context.Context, id int) (models.Headline, error) {
	if r == nil || r.db == nil {
		return models.Headline{}, fmt.Errorf("headline repository database is not configured")
	}

	var h models.Headline
	err := r.db.QueryRowContext(ctx,
		`SELECT id, title, COALESCE(slug, ''), created_at FROM headlines WHERE id = $1`, id).
		Scan(&h.ID, &h.Title, &h.Slug, &h.CreatedAt)
	return h, err
}

func slugify(title string) string {
	slug := strings.ToLower(strings.TrimSpace(title))
	slug = strings.ReplaceAll(slug, " ", "-")
	return slug
}
