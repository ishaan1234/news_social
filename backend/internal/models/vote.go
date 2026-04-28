package models

type Vote struct {
	UserID     int `json:"user_id"`
	HeadlineID int `json:"headline_id"`
	Value      int `json:"value"`
}
