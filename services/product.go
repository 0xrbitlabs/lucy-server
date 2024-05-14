package services

import (
	"errors"
	"lucy/dtos"
	"lucy/handlers"
	"lucy/models"
	"lucy/repositories"
	"lucy/types"
	"net/http"

	"github.com/oklog/ulid/v2"
)

type ProductRepo interface {
	Insert(*models.Product) error
	GetAll(*models.User) (*[]models.Product, error)
	ToggleStatus(ids []string, status bool) error
}

type ProductService struct {
	productRepo  ProductRepo
	categoryRepo CategoryRepo
	logger       handlers.Logger
}

func NewProductService(
	productRepo ProductRepo,
	categoryRepo CategoryRepo,
	logger handlers.Logger,
) ProductService {
	return ProductService{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
		logger:       logger,
	}
}

func (s ProductService) CreateProduct(data *dtos.CreateProductDTO, owner *models.User) error {
	_, err := s.categoryRepo.GetCategory(repositories.Filter{
		Field: "id",
		Value: data.CategoryID,
	})
	if err != nil {
		if errors.Is(err, types.ErrResourceNotFound) {
			return types.ServiceError{
				StatusCode: http.StatusBadRequest,
				ErrorCode:  types.CategoryNotFound,
			}
		}
		s.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	product := &models.Product{
		ID:          ulid.Make().String(),
		Owner:       owner.ID,
		CategoryID:  data.CategoryID,
		Label:       data.Label,
		Description: data.Description,
		Price:       data.Price,
		Image:       data.Image,
		Enabled:     true,
	}
	err = s.productRepo.Insert(product)
	if err != nil {
		s.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	return nil
}

func (s ProductService) GetAll(user *models.User) (*[]models.Product, error) {
	data, err := s.productRepo.GetAll(user)
	if err != nil {
		s.logger.Error(err.Error())
		return nil, types.ServiceErrInternal
	}
	return data, nil
}

func (s ProductService) ToggleStatus(data *dtos.ToggleProductStatusDTO) error {
	err := s.productRepo.ToggleStatus(data.IDs, data.Status)
	if err != nil {
		s.logger.Error(err.Error())
		return err
	}
	return nil
}
