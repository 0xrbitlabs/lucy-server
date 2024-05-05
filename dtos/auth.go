package dtos

type LoginDTO struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
}

func (dto LoginDTO) Validate() map[string]string {
	errors := make(map[string]string)
	if dto.Phone == "" {
		errors["phone"] = "phone can not be empty"
	}
	if dto.Password == "" {
		errors["password"] = "phone can not be empty"
	}
	return errors
}
