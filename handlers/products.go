package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joseph0x45/lucy/domain"
	"github.com/joseph0x45/lucy/dtos"
	"github.com/joseph0x45/lucy/middleware"
	"github.com/joseph0x45/lucy/repository"
	"github.com/joseph0x45/lucy/utils"
	"github.com/oklog/ulid/v2"
)

type ProductHandler struct {
	authMiddleware *middleware.AuthMiddleware
	logger         *slog.Logger
	products       *repository.ProductRepo
}

func NewProductHandler(
	authMiddleware *middleware.AuthMiddleware,
	logger *slog.Logger,
	products *repository.ProductRepo,
) *ProductHandler {
	return &ProductHandler{
		authMiddleware: authMiddleware,
		logger:         logger,
		products:       products,
	}
}

func (h *ProductHandler) HandleProductCategoryCreation(w http.ResponseWriter, r *http.Request) {
	type dto struct {
		Label string `json:"label"`
	}
	payload := &dto{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Error while decoding body: %s", err.Error()),
		)
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	exists, err := h.products.CategoryExists(payload.Label)
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if exists {
		utils.WriteData(w, http.StatusConflict, nil)
		return
	}
	err = h.products.CreateProductCategory(&domain.ProductCategory{
		Label: payload.Label,
	})
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	utils.WriteData(w, http.StatusCreated, nil)
}

func (h *ProductHandler) HandleGetAllProductCategories(w http.ResponseWriter, r *http.Request) {
	data, err := h.products.GetAllProductCategory()
	if err != nil {
		h.logger.Error(err.Error())
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	utils.WriteData(w, http.StatusOK, map[string]interface{}{
		"categories": data,
	})
}

func (h *ProductHandler) HandleProductCreation(w http.ResponseWriter, r *http.Request) {
	payload := &dtos.CreateProductDTO{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error(
			fmt.Sprintf("Error while decoding body: %s", err.Error()),
		)
		utils.WriteError(w, http.StatusInternalServerError, nil)
		return
	}
	if !payload.IsValid() {
		utils.WriteError(w, http.StatusBadRequest, nil)
		return
	}
}

func (h *ProductHandler) HandleGetAllProducts(w http.ResponseWriter, r *http.Request) {

}

func (h *ProductHandler) HandleGetAllMyProducts(w http.ResponseWriter, r *http.Request) {

}

func (h *ProductHandler) HandleUpdateProduct(w http.ResponseWriter, r *http.Request) {}

func (h *ProductHandler) HandleDeleteProduct(w http.ResponseWriter, r *http.Request) {}

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	auth := h.authMiddleware.Authenticate()
	r.Route("/categories", func(r chi.Router) {
		r.With(
			h.authMiddleware.Authenticate(string(domain.AdminAccountType)),
		).Post("/", h.HandleProductCategoryCreation)
		r.With(auth).Get("/", h.HandleGetAllProductCategories)
	})
}
