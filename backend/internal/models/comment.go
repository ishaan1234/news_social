package models

import "time"

type Comment struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	HeadlineID int64     `json:"headline_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
