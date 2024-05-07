package interfaces

import (
	"lucy/models"
	"lucy/repositories"
)

type UserRepo interface {
	Insert(*models.User) error
	GetUser(filter repositories.Filter) (*models.User, error)
	CountByPhone(phone string) (int, error)
	GetAll() (*[]models.User, error)
	UpdatePassword(userId, password string) error
}
