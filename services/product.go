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
	GetAll() (*[]models.Product, error)
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
	}
	err = s.productRepo.Insert(product)
	if err != nil {
		s.logger.Error(err.Error())
		return types.ServiceErrInternal
	}
	return nil
}

func (s ProductService) GetAll() (*[]models.Product, error) {
	data, err := s.productRepo.GetAll()
	if err != nil {
		s.logger.Error(err.Error())
		return nil, types.ServiceErrInternal
	}
	return data, nil
}
