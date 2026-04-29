package models

import "time"

type LinkedArticle struct {
	URL         string `json:"url"`
	Title       string `json:"title"`
	Source      string `json:"source,omitempty"`
	Summary     string `json:"summary,omitempty"`
	ImageURL    string `json:"image_url,omitempty"`
	PublishedAt string `json:"published_at,omitempty"`
}

type Post struct {
	ID           int           `json:"id"`
	AuthorID     string        `json:"author_id,omitempty"`
	AuthorName   string        `json:"author_name"`
	AuthorHandle string        `json:"author_handle,omitempty"`
	Body         string        `json:"body"`
	Article      LinkedArticle `json:"article"`
	VoteScore    int           `json:"vote_score"`
	ViewerVote   int           `json:"viewer_vote,omitempty"`
	CommentCount int           `json:"comment_count"`
	ShareCount   int           `json:"share_count"`
	CreatedAt    time.Time     `json:"created_at"`
}

type PostComment struct {
	ID         int       `json:"id"`
	PostID     int       `json:"post_id"`
	AuthorID   string    `json:"author_id,omitempty"`
	AuthorName string    `json:"author_name"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}

type PostVoteSummary struct {
	PostID     int `json:"post_id"`
	VoteScore  int `json:"vote_score"`
	ViewerVote int `json:"viewer_vote"`
}
