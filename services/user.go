package services

import (
	"lucy/dtos"
	"lucy/models"
	"lucy/types"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Insert(*models.User) error
	GetByID(id string) (*models.User, error)
	CountByPhone(phone string) (int, error)
	GetAll() (*[]models.User, error)
}

type UserService struct {
	UserRepo
	types.ILogger
}

func (svc UserService) CreateAdminAccount(data dtos.CreateAdminDTO) error {
	count, err := svc.CountByPhone(data.Phone)
	if err != nil {
		svc.Error(err.Error())
		return types.ServiceErrInternal
	}
	if count != 0 {
		return types.ServiceErrDuplicatePhone
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		svc.Error("Error while hashing password: ", err)
		return types.ServiceErrInternal
	}
	newUser := &models.User{
		ID:          ulid.Make().String(),
		Username:    data.Username,
		Phone:       data.Phone,
		Password:    string(hash),
		AccountType: types.AdminAccount,
		Description: "",
		Country:     "",
		Town:        "",
	}
	err = svc.Insert(newUser)
	if err != nil {

	}
	return nil
}
