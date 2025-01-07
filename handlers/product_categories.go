package handlers

import (
	"encoding/json"
	"log/slog"
	"lucy/middlewares"
	"lucy/models"
	"lucy/repo"
	"lucy/utils"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/ulid/v2"
)

type ProductCategoryHandler struct {
	productCategories *repo.ProductCategoryRepo
	logger            *slog.Logger
}

func NewProductCategoryHandler(
	productCategories *repo.ProductCategoryRepo,
	logger *slog.Logger,
) *ProductCategoryHandler {
	return &ProductCategoryHandler{
		productCategories: productCategories,
		logger:            logger,
	}
}

func (h *ProductCategoryHandler) CreateProductCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	payload := &struct {
		Label       string `json:"label"`
		Description string `json:"description"`
	}{}
	err := json.NewDecoder(r.Body).Decode(payload)
	if err != nil {
		h.logger.Error("Error while decoding payload body:", slog.Any("err", err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if payload.Label == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteError("empty_label", w)
		return
	}
	if payload.Description == "" {
		w.WriteHeader(http.StatusBadRequest)
		utils.WriteError("empty_description", w)
		return
	}
	labelIsUnique, err := h.productCategories.LabelIsUnique(payload.Label)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if !labelIsUnique {
		w.WriteHeader(http.StatusConflict)
		utils.WriteError("label_exists", w)
		return
	}
	productCategory := &models.ProductCategory{
		ID:          ulid.Make().String(),
		Label:       payload.Label,
		Description: payload.Description,
		CreatedAt:   time.Now().UTC(),
	}
	err = h.productCategories.Insert(productCategory)
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *ProductCategoryHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	data, err := h.productCategories.GetAll()
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		h.logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	utils.WriteData(map[string]interface{}{
		"data": data,
	}, w)
}

func (h *ProductCategoryHandler) RegisterRoutes(r chi.Router, m *middlewares.AuthMiddleware) {
	adminOnlyAuth := m.AuthenticateWithRole("admin")
	auth := m.AuthenticateWithRole()
	r.Route("/categories", func(r chi.Router) {
		r.With(adminOnlyAuth).Post("/", h.CreateProductCategory)
		r.With(auth).Get("/", h.GetAll)
	})
}
