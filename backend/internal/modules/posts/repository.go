package posts

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ishaan1234/news_social/backend/internal/models"
)

type Repository interface {
	CreatePost(ctx context.Context, post models.Post) (models.Post, error)
	GetPosts(ctx context.Context, viewerID string) ([]models.Post, error)
	GetPostByID(ctx context.Context, postID int, viewerID string) (models.Post, error)
	CreateComment(ctx context.Context, comment models.PostComment) (models.PostComment, error)
	GetComments(ctx context.Context, postID int) ([]models.PostComment, error)
	SetVote(ctx context.Context, postID int, voterID string, value int) (models.PostVoteSummary, error)
	IncrementShare(ctx context.Context, postID int) (int, error)
}

type SQLRepository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *SQLRepository {
	return &SQLRepository{db: db}
}

func (r *SQLRepository) CreatePost(ctx context.Context, post models.Post) (models.Post, error) {
	if r == nil || r.db == nil {
		return models.Post{}, fmt.Errorf("post repository database is not configured")
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO posts (
			author_id,
			author_name,
			author_handle,
			body,
			article_url,
			article_title,
			article_source,
			article_summary,
			article_image_url,
			article_published_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id, share_count, created_at
	`,
		post.AuthorID,
		post.AuthorName,
		post.AuthorHandle,
		post.Body,
		post.Article.URL,
		post.Article.Title,
		post.Article.Source,
		post.Article.Summary,
		post.Article.ImageURL,
		post.Article.PublishedAt,
	).Scan(&post.ID, &post.ShareCount, &post.CreatedAt)
	if err != nil {
		return models.Post{}, fmt.Errorf("create post: %w", err)
	}

	return post, nil
}

func (r *SQLRepository) GetPosts(ctx context.Context, viewerID string) ([]models.Post, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("post repository database is not configured")
	}

	rows, err := r.db.QueryContext(ctx, postSelectSQL()+` ORDER BY p.created_at DESC`, viewerID)
	if err != nil {
		return nil, fmt.Errorf("get posts: %w", err)
	}
	defer rows.Close()

	result := []models.Post{}
	for rows.Next() {
		post, err := scanPost(rows)
		if err != nil {
			return nil, err
		}
		result = append(result, post)
	}

	return result, rows.Err()
}

func (r *SQLRepository) GetPostByID(ctx context.Context, postID int, viewerID string) (models.Post, error) {
	if r == nil || r.db == nil {
		return models.Post{}, fmt.Errorf("post repository database is not configured")
	}

	row := r.db.QueryRowContext(ctx, postSelectSQL()+` WHERE p.id = $2`, viewerID, postID)
	post, err := scanPost(row)
	if err != nil {
		return models.Post{}, fmt.Errorf("get post: %w", err)
	}
	return post, nil
}

func (r *SQLRepository) CreateComment(ctx context.Context, comment models.PostComment) (models.PostComment, error) {
	if r == nil || r.db == nil {
		return models.PostComment{}, fmt.Errorf("post repository database is not configured")
	}

	err := r.db.QueryRowContext(ctx, `
		INSERT INTO post_comments (post_id, author_id, author_name, content)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at
	`, comment.PostID, comment.AuthorID, comment.AuthorName, comment.Content).
		Scan(&comment.ID, &comment.CreatedAt)
	if err != nil {
		return models.PostComment{}, fmt.Errorf("create post comment: %w", err)
	}

	return comment, nil
}

func (r *SQLRepository) GetComments(ctx context.Context, postID int) ([]models.PostComment, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("post repository database is not configured")
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT id, post_id, COALESCE(author_id, ''), author_name, content, created_at
		FROM post_comments
		WHERE post_id = $1
		ORDER BY created_at ASC
	`, postID)
	if err != nil {
		return nil, fmt.Errorf("get post comments: %w", err)
	}
	defer rows.Close()

	comments := []models.PostComment{}
	for rows.Next() {
		var comment models.PostComment
		if err := rows.Scan(
			&comment.ID,
			&comment.PostID,
			&comment.AuthorID,
			&comment.AuthorName,
			&comment.Content,
			&comment.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan post comment: %w", err)
		}
		comments = append(comments, comment)
	}

	return comments, rows.Err()
}

func (r *SQLRepository) SetVote(ctx context.Context, postID int, voterID string, value int) (models.PostVoteSummary, error) {
	if r == nil || r.db == nil {
		return models.PostVoteSummary{}, fmt.Errorf("post repository database is not configured")
	}

	if value == 0 {
		if _, err := r.db.ExecContext(ctx,
			`DELETE FROM post_votes WHERE post_id = $1 AND voter_id = $2`,
			postID,
			voterID,
		); err != nil {
			return models.PostVoteSummary{}, fmt.Errorf("clear post vote: %w", err)
		}
		return r.getVoteSummary(ctx, postID, voterID)
	}

	if _, err := r.db.ExecContext(ctx, `
		INSERT INTO post_votes (post_id, voter_id, value)
		VALUES ($1, $2, $3)
		ON CONFLICT (post_id, voter_id) DO UPDATE
		SET value = EXCLUDED.value,
		    updated_at = NOW()
	`, postID, voterID, value); err != nil {
		return models.PostVoteSummary{}, fmt.Errorf("set post vote: %w", err)
	}

	return r.getVoteSummary(ctx, postID, voterID)
}

func (r *SQLRepository) IncrementShare(ctx context.Context, postID int) (int, error) {
	if r == nil || r.db == nil {
		return 0, fmt.Errorf("post repository database is not configured")
	}

	var shareCount int
	if err := r.db.QueryRowContext(ctx, `
		UPDATE posts
		SET share_count = share_count + 1
		WHERE id = $1
		RETURNING share_count
	`, postID).Scan(&shareCount); err != nil {
		return 0, fmt.Errorf("share post: %w", err)
	}

	return shareCount, nil
}

func (r *SQLRepository) getVoteSummary(ctx context.Context, postID int, voterID string) (models.PostVoteSummary, error) {
	var summary models.PostVoteSummary
	summary.PostID = postID

	err := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE((SELECT SUM(value) FROM post_votes WHERE post_id = $1), 0)::int,
			COALESCE((SELECT value FROM post_votes WHERE post_id = $1 AND voter_id = $2), 0)::int
	`, postID, voterID).Scan(&summary.VoteScore, &summary.ViewerVote)
	if err != nil {
		return models.PostVoteSummary{}, fmt.Errorf("get post vote summary: %w", err)
	}

	return summary, nil
}

func postSelectSQL() string {
	return `
		SELECT
			p.id,
			COALESCE(p.author_id, ''),
			p.author_name,
			COALESCE(p.author_handle, ''),
			p.body,
			p.article_url,
			p.article_title,
			COALESCE(p.article_source, ''),
			COALESCE(p.article_summary, ''),
			COALESCE(p.article_image_url, ''),
			COALESCE(p.article_published_at, ''),
			COALESCE((SELECT SUM(value) FROM post_votes WHERE post_id = p.id), 0)::int,
			COALESCE((SELECT value FROM post_votes WHERE post_id = p.id AND voter_id = $1), 0)::int,
			COALESCE((SELECT COUNT(*) FROM post_comments WHERE post_id = p.id), 0)::int,
			p.share_count,
			p.created_at
		FROM posts p
	`
}

type postScanner interface {
	Scan(dest ...interface{}) error
}

func scanPost(scanner postScanner) (models.Post, error) {
	var post models.Post
	if err := scanner.Scan(
		&post.ID,
		&post.AuthorID,
		&post.AuthorName,
		&post.AuthorHandle,
		&post.Body,
		&post.Article.URL,
		&post.Article.Title,
		&post.Article.Source,
		&post.Article.Summary,
		&post.Article.ImageURL,
		&post.Article.PublishedAt,
		&post.VoteScore,
		&post.ViewerVote,
		&post.CommentCount,
		&post.ShareCount,
		&post.CreatedAt,
	); err != nil {
		return models.Post{}, fmt.Errorf("scan post: %w", err)
	}
	return post, nil
}
