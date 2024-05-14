package handlers

import (
	"encoding/json"
	"lucy/dtos"
	"lucy/models"
	"lucy/types"
	"net/http"
)

type ProductService interface {
	CreateProduct(*dtos.CreateProductDTO, *models.User) error
	GetAll(*models.User) (*[]models.Product, error)
	ToggleStatus(*dtos.ToggleProductStatusDTO, *models.User) error
}

type ProductHandler struct {
	productService ProductService
	logger         Logger
}

func NewProductHandler(productService ProductService, logger Logger) ProductHandler {
	return ProductHandler{
		productService: productService,
		logger:         logger,
	}
}

func (h ProductHandler) HandleCreateProduct(w http.ResponseWriter, r *http.Request) {
	currUser, _ := r.Context().Value("user").(*models.User)
	payload := &dtos.CreateProductDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.productService.CreateProduct(payload, currUser)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusCreated, nil)
}

func (h ProductHandler) HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {
	currUser, _ := r.Context().Value("user").(*models.User)
	data, err := h.productService.GetAll(currUser)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusOK, map[string]interface{}{
		"products": *data,
	})
}

func (h ProductHandler) HandleToggleEnabled(w http.ResponseWriter, r *http.Request) {
	currUser, _ := r.Context().Value("user").(*models.User)
	payload := &dtos.ToggleProductStatusDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding request body: ", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = h.productService.ToggleStatus(payload, currUser)
	if err != nil {
		WriteError(w, err.(types.ServiceError))
		return
	}
	WriteData(w, http.StatusOK, nil)
}
