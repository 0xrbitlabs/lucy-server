package dtos

type CreateAdminDTO struct {
	Username string `json:"username"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
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

type ChangeUserPasswordDTO struct {
	UserID      string
	OldPassword string `json:"old"`
	NewPassword string `json:"new"`
}

func (dto ChangeUserPasswordDTO) Validate() map[string]string {
	errors := make(map[string]string)
	if dto.OldPassword == "" {
		errors["old"] = "old can not be empty"
	}
	if dto.NewPassword == "" {
		errors["new"] = "new can not be empty"
	}
	if len(errors) > 0 {
		return errors
	}
	return nil
}
