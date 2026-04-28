package models

import "time"

type Article struct {
	ID          int       `json:"id"`
	HeadlineID  int       `json:"headline_id"`
	Source      string    `json:"source"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	Content     string    `json:"content"`
	PublishedAt time.Time `json:"published_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}
