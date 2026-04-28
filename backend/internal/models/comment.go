package models

import "time"

type Comment struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	HeadlineID int       `json:"headline_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
