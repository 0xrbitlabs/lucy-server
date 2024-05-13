package dtos

type CreateCategoryDTO struct {
	Label       string `json:"label"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

func (dto CreateCategoryDTO) Validate() map[string]string {
	errors := make(map[string]string)
	if dto.Label == "" {
		errors["label"] = "label can not be empty"
	}
	if dto.Description == "" {
		errors["description"] = "description can not be empty"
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}

type ToggleEnabledDTO struct {
	IDs    []string `json:"ids"`
	Status bool     `json:"status"`
}
