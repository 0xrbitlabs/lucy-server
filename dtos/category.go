package dtos

type CreateCategoryDTO struct {
	Label string `json:"label"`
}

func (dto CreateCategoryDTO) Validate() map[string]string {
	errors := make(map[string]string)
	if dto.Label == "" {
		errors["label"] = "label can not be empty"
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
