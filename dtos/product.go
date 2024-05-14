package dtos

type CreateProductDTO struct {
	CategoryID  string `json:"category_id"`
	Label       string `json:"label"`
	Price       int    `json:"price"`
	Description string `json:"description"`
	Image       string `json:"image"`
}

type ToggleProductStatusDTO struct {
	IDs []string `json:"ids"`
  Status bool `json:"status"`
}
