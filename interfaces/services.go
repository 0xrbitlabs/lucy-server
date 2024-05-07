package interfaces

import (
	"lucy/dtos"
	"lucy/models"
)

type AuthService interface {
	Login(data dtos.LoginDTO) (*string, error)
}

type UserService interface {
	CreateAdminAccount(dtos.CreateAdminDTO) error
	GetAllUsers() (*[]models.User, error)
	GetUserByID(id string) (*models.User, error)
}
