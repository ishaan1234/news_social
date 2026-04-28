package models

import "time"

type Summary struct {
	ID         int       `json:"id"`
	HeadlineID int       `json:"headline_id"`
	Content    string    `json:"content"`
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
}
