package models

import (
	"time"

	"github.com/google/uuid"
)

type Summary struct {
	ID         uuid.UUID `json:"id"`
	HeadlineID uuid.UUID `json:"headline_id"`
	Content    string    `json:"content"`
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
}
