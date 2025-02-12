package models

import "time"

type User struct {
	ID          string    `json:"id" db:"id"`
	Username    string    `json:"username" db:"username"`
	PhoneNumber string    `json:"phone_number" db:"phone_number"`
	Password    string    `json:"password" db:"password"`
	AccountType string    `json:"account_type" db:"account_type"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	Verified    bool      `json:"verified" db:"verified"`
}
