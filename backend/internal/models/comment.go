package models

import (
	"time"

	"github.com/google/uuid"
)

type Comment struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"user_id"`
	HeadlineID uuid.UUID `json:"headline_id"`
	Content    string    `json:"content"`
	CreatedAt  time.Time `json:"created_at"`
}
