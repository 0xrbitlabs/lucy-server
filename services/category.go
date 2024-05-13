package services

import (
	"lucy/dtos"
	"lucy/handlers"
	"lucy/models"
	"lucy/repositories"
	"lucy/types"
	"net/http"

	"github.com/oklog/ulid/v2"
)

type CategoryRepo interface {
	Insert(*models.Category) error
	GetCategory(repositories.Filter) (*models.Category, error)
	GetAll(types.AccountType) (*[]models.Category, error)
	CountByLabel(string) (int, error)
}

type CategoryService struct {
	categoryRepo CategoryRepo
	logger       handlers.Logger
}

func NewCategoryService(categoryRepo CategoryRepo, logger handlers.Logger) CategoryService {
	return CategoryService{
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s CategoryService) CreateCategory(data dtos.CreateCategoryDTO) error {
	count, err := s.categoryRepo.CountByLabel(data.Label)
	if err != nil {
		s.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	if count != 0 {
		return types.ServiceError{
			StatusCode: http.StatusConflict,
			ErrorCode:  types.ErrDuplicateLabel,
		}
	}
	category := &models.Category{
		ID:          ulid.Make().String(),
		Label:       data.Label,
		Description: data.Description,
		Enabled:     data.Enabled,
	}
	err = s.categoryRepo.Insert(category)
	if err != nil {
		s.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	return nil
}

func (s CategoryService) GetAllCategories(currUser *models.User) (*[]models.Category, error) {
	categories, err := s.categoryRepo.GetAll(currUser.AccountType)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, types.ServiceErrInternal
	}
	return categories, nil
}
