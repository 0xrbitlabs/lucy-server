package services

import (
	"lucy/dtos"
	"lucy/types"
)

type AuthService struct {
	userRepo UserRepo
	logger   types.ILogger
}

func NewAuthService(userRepo UserRepo, logger types.ILogger) AuthService {
	return AuthService{
		userRepo: userRepo,
		logger:   logger,
	}
}

func (s AuthService) Login(data dtos.LoginDTO) (*string, error) {
	return nil, nil
}
