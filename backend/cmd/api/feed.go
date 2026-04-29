package main

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type feedPost struct {
	ID           string      `json:"id"`
	UserEmail    string      `json:"user_email"`
	Username     string      `json:"username,omitempty"`
	DisplayName  string      `json:"display_name,omitempty"`
	AvatarURL    string      `json:"avatar_url,omitempty"`
	Caption      string      `json:"caption,omitempty"`
	CreatedAt    time.Time   `json:"created_at"`
	LikeCount    int         `json:"like_count"`
	CommentCount int         `json:"comment_count"`
	LikedByMe    bool        `json:"liked_by_me"`
	Article      feedArticle `json:"article"`
}

type feedArticle struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description,omitempty"`
	Content     string     `json:"content,omitempty"`
	Summary     string     `json:"summary,omitempty"`
	Author      string     `json:"author,omitempty"`
	SourceName  string     `json:"source_name,omitempty"`
	SourceID    string     `json:"source_id,omitempty"`
	URL         string     `json:"url"`
	ImageURL    string     `json:"image_url,omitempty"`
	PublishedAt *time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}

func feedHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeJSONError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}

		if db == nil {
			writeJSONError(w, http.StatusInternalServerError, "database is not configured")
			return
		}

		userEmail, err := normalizeEmail(r.URL.Query().Get("user_email"))
		if err != nil {
			writeJSONError(w, http.StatusBadRequest, "valid user_email is required")
			return
		}

		feed, err := getFeed(r, db, userEmail)
		if err != nil {
			writeJSONError(w, http.StatusInternalServerError, "failed to fetch feed")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"success": true,
			"feed":    feed,
		})
	}
}

func getFeed(r *http.Request, db *sql.DB, userEmail string) ([]feedPost, error) {
	rows, err := db.QueryContext(r.Context(), `
		SELECT
			p.id::text,
			p.user_email,
			COALESCE(u.username, ''),
			COALESCE(u.display_name, ''),
			COALESCE(u.avatar_url, ''),
			COALESCE(p.caption, ''),
			p.created_at,
			COUNT(DISTINCT pl.user_email)::int AS like_count,
			COUNT(DISTINCT pc.id)::int AS comment_count,
			EXISTS (
				SELECT 1
				FROM post_likes my_like
				WHERE my_like.post_id = p.id
				  AND my_like.user_email = $1
			) AS liked_by_me,
			a.id::text,
			a.title,
			COALESCE(a.description, ''),
			COALESCE(a.content, ''),
			COALESCE(a.summary, ''),
			COALESCE(a.author, ''),
			COALESCE(a.source_name, ''),
			COALESCE(a.source_id, ''),
			a.url,
			COALESCE(a.image_url, ''),
			a.published_at,
			a.created_at
		FROM posts p
		JOIN users u ON u.email = p.user_email
		JOIN articles a ON a.id = p.article_id
		LEFT JOIN post_likes pl ON pl.post_id = p.id
		LEFT JOIN post_comments pc ON pc.post_id = p.id
		WHERE
			p.user_email = $1
			OR p.user_email IN (
				SELECT following_email
				FROM following
				WHERE follower_email = $1
			)
		GROUP BY p.id, u.email, a.id
		ORDER BY p.created_at DESC
	`, userEmail)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feed []feedPost
	for rows.Next() {
		var post feedPost
		var publishedAt sql.NullTime

		if err := rows.Scan(
			&post.ID,
			&post.UserEmail,
			&post.Username,
			&post.DisplayName,
			&post.AvatarURL,
			&post.Caption,
			&post.CreatedAt,
			&post.LikeCount,
			&post.CommentCount,
			&post.LikedByMe,
			&post.Article.ID,
			&post.Article.Title,
			&post.Article.Description,
			&post.Article.Content,
			&post.Article.Summary,
			&post.Article.Author,
			&post.Article.SourceName,
			&post.Article.SourceID,
			&post.Article.URL,
			&post.Article.ImageURL,
			&publishedAt,
			&post.Article.CreatedAt,
		); err != nil {
			return nil, err
		}

		if publishedAt.Valid {
			post.Article.PublishedAt = &publishedAt.Time
		}

		feed = append(feed, post)
	}

	return feed, rows.Err()
}
