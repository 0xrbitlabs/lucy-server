package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"lucy/models"

	"github.com/jmoiron/sqlx"
)

type ProductCategoryRepo struct {
	db *sqlx.DB
}

func NewProductCategoryRepo(db *sqlx.DB) *ProductCategoryRepo {
	return &ProductCategoryRepo{db}
}

func (r *ProductCategoryRepo) Insert(pc *models.ProductCategory) error {
	const query = `
    INSERT INTO product_categories (
      id, label, description, created_at
    )
    VALUES (
      :id, :label, :description, :created_at
    )
  `
	_, err := r.db.NamedExec(query, pc)
	if err != nil {
		return fmt.Errorf("Error while inserting product category: %w", err)
	}
	return nil
}

func (r *ProductCategoryRepo) GetAll() ([]models.ProductCategory, error) {
	const query = "select * from product_categories"
	data := make([]models.ProductCategory, 0)
	err := r.db.Select(&data, query)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all product categories: %w", err)
	}
	return data, nil
}

func (r *ProductCategoryRepo) GetByLabel(label string) (*models.ProductCategory, error) {
	pc := &models.ProductCategory{}
	const query = "select * from product_categories where label=$1"
	err := r.db.Get(pc, query, label)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting product category by label: %w", err)
	}
	return pc, nil
}

func (r *ProductCategoryRepo) GetByID(id string) (*models.ProductCategory, error) {
	pc := &models.ProductCategory{}
	const query = "select * from product_categories where id=$1"
	err := r.db.Get(pc, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting product category by id: %w", err)
	}
	return pc, nil
}

func (r *ProductCategoryRepo) LabelIsUnique(label string) (bool, error) {
	pc, err := r.GetByLabel(label)
	if err != nil {
		return false, err
	}
	return pc == nil, nil
}
