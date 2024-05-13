package services

import (
	"errors"
	"lucy/dtos"
	"lucy/handlers"
	"lucy/models"
	"lucy/repositories"
	"lucy/types"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type UserRepo interface {
	Insert(*models.User) error
	GetUser(filter repositories.Filter) (*models.User, error)
	CountByPhone(phone string) (int, error)
	GetAll() (*[]models.User, error)
	UpdatePassword(userId, password string) error
}

type AuthService struct {
	userRepo UserRepo
	logger   handlers.Logger
	jwt      handlers.JWTProvider
}

func NewAuthService(userRepo UserRepo, logger handlers.Logger, jwt handlers.JWTProvider) AuthService {
	return AuthService{
		userRepo: userRepo,
		logger:   logger,
		jwt:      jwt,
	}
}

func (s AuthService) Login(data dtos.LoginDTO) (*string, error) {
	user, err := s.userRepo.GetUser(repositories.Filter{Field: "phone", Value: data.Phone})
	if err != nil {
		if errors.Is(err, types.ErrResourceNotFound) {
			return nil, types.ServiceError{
				StatusCode: http.StatusBadRequest,
				ErrorCode:  types.ErrPhoneNotFound,
			}
		}
		s.logger.Error(err.Error())
		return nil, types.ServiceErrInternal
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(data.Password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, types.ServiceError{
				StatusCode: http.StatusBadRequest,
				ErrorCode:  types.ErrWrongPassword,
			}
		}
		s.logger.Error("Error while comparing hash and password:", err)
		return nil, types.ServiceErrInternal
	}
	authToken, err := s.jwt.Encode(map[string]interface{}{
		"id": user.ID,
	})
	if err != nil {
		s.logger.Error(err.Error())
		return nil, types.ServiceError{
			StatusCode: http.StatusInternalServerError,
			ErrorCode:  types.ErrTokenEncoding,
		}
	}
	return &authToken, nil
}
