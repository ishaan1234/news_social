package models

import (
	"time"

	"github.com/google/uuid"
)

type Headline struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	CreatedAt time.Time `json:"created_at"`
}
