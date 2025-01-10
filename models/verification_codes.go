package models

import "time"

type VerificationCode struct {
	ID           int       `json:"id" db:"id"`
	Code         string    `json:"code" db:"code"`
	GeneratedFor string    `json:"generated_for" db:"generated_for"`
	GeneratedAt  time.Time `json:"generated_at" db:"generated_at"`
	Used         bool      `json:"used" db:"used"`
}
