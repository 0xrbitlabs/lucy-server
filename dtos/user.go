package dtos

type CreateAdminDTO struct {
	Username string `json:"username" db:"username"`
	Phone    string `json:"phone" db:"phone"`
	Password string `json:"password" db:"password"`
}

func (dto CreateAdminDTO) Validate() map[string]string {
	errors := make(map[string]string)
	if dto.Username == "" {
		errors["username"] = "username can not be empty"
	}
	if dto.Password == "" {
		errors["password"] = "password not be empty"
	}
	if dto.Phone == "" {
		errors["phone"] = "phone can not be empty"
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
