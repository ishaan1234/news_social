package social

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Repository interface {
	CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error)
	GetCommentsByHeadline(ctx context.Context, headlineID int) ([]models.Comment, error)
	GetCommentsByUser(ctx context.Context, userID int) ([]models.Comment, error)
	DeleteComment(ctx context.Context, commentID int, userID int) error
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) CreateComment(ctx context.Context, comment models.Comment) (models.Comment, error) {
	if r == nil || r.db == nil {
		return models.Comment{}, fmt.Errorf("social repository database is not configured")
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO comments (user_id, headline_id, content)
		VALUES ($1, $2, $3)
		RETURNING id, created_at
	`, comment.UserID, comment.HeadlineID, comment.Content).Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return models.Comment{}, fmt.Errorf("create comment: %w", err)
	}
	return comment, nil
}

func (r *SQLRepository) GetCommentsByHeadline(ctx context.Context, headlineID int) ([]models.Comment, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("social repository database is not configured")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, headline_id, content, created_at
		FROM comments
		WHERE headline_id = $1
		ORDER BY created_at DESC
	`, headlineID)
	if err != nil {
		return nil, fmt.Errorf("get comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.HeadlineID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *SQLRepository) GetCommentsByUser(ctx context.Context, userID int) ([]models.Comment, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("social repository database is not configured")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, headline_id, content, created_at
		FROM comments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("get user comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(&c.ID, &c.UserID, &c.HeadlineID, &c.Content, &c.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan user comment: %w", err)
		}
		comments = append(comments, c)
	}
	return comments, rows.Err()
}

func (r *SQLRepository) DeleteComment(ctx context.Context, commentID int, userID int) error {
	if r == nil || r.db == nil {
		return fmt.Errorf("social repository database is not configured")
	}

	result, err := r.db.ExecContext(ctx, `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
	`, commentID, userID)
	if err != nil {
		return fmt.Errorf("delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check deleted comment: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("comment not found or unauthorized")
	}
	return nil
}
