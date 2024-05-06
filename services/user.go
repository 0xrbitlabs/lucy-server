package services

import (
	"lucy/dtos"
	"lucy/interfaces"
	"lucy/models"
	"lucy/repositories"
	"lucy/types"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Insert(*models.User) error
	GetUser(filter repositories.Filter) (*models.User, error)
	CountByPhone(phone string) (int, error)
	GetAll() (*[]models.User, error)
}

type UserService struct {
	userRepo UserRepo
	logger   interfaces.Logger
}

func NewUserService(userRepo UserRepo, logger interfaces.Logger) UserService {
	return UserService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (svc UserService) CreateAdminAccount(data dtos.CreateAdminDTO) error {
	count, err := svc.userRepo.CountByPhone(data.Phone)
	if err != nil {
		svc.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	if count != 0 {
		return types.ServiceErrDuplicatePhone
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data.Password), bcrypt.DefaultCost)
	if err != nil {
		svc.logger.Error("Error while hashing password: ", err)
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
	err = svc.userRepo.Insert(newUser)
	if err != nil {
		svc.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	return nil
}

func (svc UserService) GetAllUsers() (*[]models.User, error) {
	return nil, nil
}
func (svc UserService) GetUserByID(id string) (*models.User, error) {
	return nil, nil
}
