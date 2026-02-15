package models

import "time"

type Article struct {
	ID          int64     `json:"id"`
	HeadlineID  int64     `json:"headline_id"`
	Source      string    `json:"source"`
	Title       string    `json:"title"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
}
