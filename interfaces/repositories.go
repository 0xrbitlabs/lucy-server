package interfaces

import (
	"lucy/models"
	"lucy/repositories"
	"lucy/types"
)

type UserRepo interface {
	Insert(*models.User) error
	GetUser(filter repositories.Filter) (*models.User, error)
	CountByPhone(phone string) (int, error)
	GetAll() (*[]models.User, error)
	UpdatePassword(userId, password string) error
}

type CategoryRepo interface {
	Insert(*models.Category) error
	GetCategory(filter repositories.Filter) (*models.Category, error)
	GetAll(callerAccountType types.AccountType) (*[]models.Category, error)
	CountByLabel(label string) (int, error)
}
