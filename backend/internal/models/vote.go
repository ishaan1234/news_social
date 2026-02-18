package models

import "github.com/google/uuid"

type Vote struct {
	UserID     uuid.UUID `json:"user_id"`
	HeadlineID uuid.UUID `json:"headline_id"`
	Value      int       `json:"value"`
}
