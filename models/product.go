package models

type Product struct {
	ID          string `json:"id" db:"id"`
	Owner       string `json:"owner" db:"owner"`
	CategoryID  string `json:"category_id" db:"category_id"`
	Label       string `json:"label" db:"label"`
	Description string `json:"description" db:"description"`
	Price       int    `json:"price" db:"price"`
	Image       string `json:"image" db:"image"`
	Enabled     bool   `json:"enabled" db:"enabled"`
}
