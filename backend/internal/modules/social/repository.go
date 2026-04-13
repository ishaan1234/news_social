package social

import (
	"context"
	"database/sql"
	"fmt"
	"time"
	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

//Create Comment
func (r *Repository) CreateComment(ctx context.Context, comment models.Comment) (int, error) {
	var id int

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO comments (user_id, headline_id, content, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`,
		comment.UserID,
		comment.HeadlineID,
		comment.Content,
		time.Now(),
	).Scan(&id)

	if err != nil {
		return 0, fmt.Errorf("failed to create comment: %w", err)
	}

	return id, nil
}

//Get Comments by Headline
func (r *Repository) GetCommentsByHeadline(ctx context.Context, headlineID int) ([]models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, headline_id, content, created_at
		FROM comments
		WHERE headline_id = $1
		ORDER BY created_at DESC
	`, headlineID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.HeadlineID,
			&c.Content,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	return comments, nil
}

// Get Comments by User
func (r *Repository) GetCommentsByUser(ctx context.Context, userID int) ([]models.Comment, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, headline_id, content, created_at
		FROM comments
		WHERE user_id = $1
		ORDER BY created_at DESC
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user comments: %w", err)
	}
	defer rows.Close()

	var comments []models.Comment

	for rows.Next() {
		var c models.Comment
		if err := rows.Scan(
			&c.ID,
			&c.UserID,
			&c.HeadlineID,
			&c.Content,
			&c.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan comment: %w", err)
		}
		comments = append(comments, c)
	}

	return comments, nil
}

//Delete Comment
func (r *Repository) DeleteComment(ctx context.Context, commentID int, userID int) error {
	result, err := r.db.ExecContext(ctx, `
		DELETE FROM comments
		WHERE id = $1 AND user_id = $2
	`, commentID, userID)

	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no comment found or unauthorized")
	}

	return nil
}