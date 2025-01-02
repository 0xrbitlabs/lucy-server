package models

type Session struct {
	ID     string `json:"id" db:"id"`
	UserID string `json:"user_id" db:"user_id"`
	Valid  bool   `json:"valid" db:"valid"`
}
