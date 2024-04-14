package types

type Category struct {
	ID          string `json:"id" db:"id"`
	Label       string `json:"label" db:"label"`
	Description string `json:"description" db:"description"`
}
