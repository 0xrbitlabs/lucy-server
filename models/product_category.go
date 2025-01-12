package models

import "time"

type ProductCategory struct {
	ID          string    `json:"id" db:"id"`
	Label       string    `json:"label" db:"label"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

type ProductCategoryCreationRequest struct {
	ID          string `json:"id" db:"id"`
	Requester   string `json:"requester" db:"requester"`
	Label       string `json:"label" db:"label"`
	Description string `json:"description" db:"description"`
	Status      string `json:"status" db:"status"`
}
