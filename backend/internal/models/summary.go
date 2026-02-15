package models

import "time"

type Summary struct {
	ID         int64     `json:"id"`
	HeadlineID int64     `json:"headline_id"`
	Content    string    `json:"content"`
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
}
