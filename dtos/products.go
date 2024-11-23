package dtos

type CreateProductDTO struct {
	Brand       string `json:"brand"`
	Category    string `json:"category"`
	Color       string `json:"color"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Label       string `json:"label"`
	Size        string `json:"size"`
	Price       string `json:"price"`
}

func (d *CreateProductDTO) IsValid() bool {
	if d.Brand == "" {
		return false
	}
	if d.Category == "" {
		return false
	}
	if d.Label == "" {
		return false
	}
	if d.Color == "" {
		return false
	}
	if d.Image == "" {
		return false
	}
	if d.Size == "" {
		return false
	}
	if d.Price == "" {
		return false
	}
	return true
}
