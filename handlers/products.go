package handlers

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/joseph0x45/lucy/domain"
	"github.com/joseph0x45/lucy/middleware"
	"github.com/joseph0x45/lucy/repository"
	"github.com/joseph0x45/lucy/utils"
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
		utils.WriteData(w, http.StatusOK, nil)
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

func (h *ProductHandler) HandleProductCreation(w http.ResponseWriter, r *http.Request)

func (h *ProductHandler) HandleGetAllProducts()

func (h *ProductHandler) HandleGetAllMyProducts()

func (h *ProductHandler) HandleUpdateProduct()

func (h *ProductHandler) HandleDeleteProduct()

func (h *ProductHandler) RegisterRoutes(r chi.Router) {
	adminAuth := h.authMiddleware.Authenticate(string(domain.AdminAccountType))
	sellerAuth := h.authMiddleware.Authenticate(string(domain.SellerAccountType))
	r.With(adminAuth).Route("/categories", func(r chi.Router) {

	})
	r.With(h.authMiddleware.Authenticate()).Route("/categories", func(r chi.Router) {
		r.Get("/", h.HandleGetAllProductCategories)
	})
	r.With(h.authMiddleware.Authenticate()).Route("/products", func(r chi.Router) {

	})
}
