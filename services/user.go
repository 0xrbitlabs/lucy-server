package services

import (
	"errors"
	"lucy/dtos"
	"lucy/interfaces"
	"lucy/models"
	"lucy/repositories"
	"lucy/types"
	"net/http"

	"github.com/oklog/ulid/v2"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo interfaces.UserRepo
	logger   interfaces.Logger
}

func NewUserService(userRepo interfaces.UserRepo, logger interfaces.Logger) UserService {
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

func (svc UserService) ChangePassword(data dtos.ChangeUserPasswordDTO) error {
	user, err := svc.userRepo.GetUser(repositories.Filter{Field: "id", Value: data.UserID})
	if err != nil {
		if errors.Is(err, types.ErrResourceNotFound) {
			return err
		}
		svc.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.OldPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return types.ServiceError{
				StatusCode: http.StatusBadRequest,
				ErrorCode:  types.ErrWrongPassword,
			}
		}
		svc.logger.Error("Error while comparing hash and password:", err)
		return types.ServiceErrInternal
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(data.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		svc.logger.Error("Error while hashing password", err)
		return types.ServiceErrInternal
	}
	err = svc.userRepo.UpdatePassword(data.UserID, string(hash))
	if err != nil {
		svc.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	return nil
}
