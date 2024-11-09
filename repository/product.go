package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/lucy/domain"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db}
}

func (r *ProductRepo) GetAllProductCategory() ([]domain.ProductCategory, error) {
	const query = "select * from product_categories"
	data := make([]domain.ProductCategory, 0)
	err := r.db.Select(&data, query)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all product categories: %w", err)
	}
	return data, nil
}

func (r *ProductRepo) InsertProductCategory(pc *domain.ProductCategory) error {
	const query = `
    insert into product_category(
      id, label, active
    )
    values (
      :id, :label, :active
    )
  `
	_, err := r.db.NamedExec(query, pc)
	if err != nil {
		return fmt.Errorf("Error while inserting product category: %w", err)
	}
	return nil
}

func (r *ProductRepo) SetProductCategoryStatus(id string, active bool) error {
	const query = "update product_categories set active=$1 where id=$2"
	_, err := r.db.Exec(query, active, id)
	if err != nil {
		return fmt.Errorf("Error while setting product category status: %w", err)
	}
	return nil
}

func (r *ProductRepo) GetAllProducts() ([]domain.Product, error) {
	const query = "select * from products"
	data := make([]domain.Product, 0)
	err := r.db.Select(&data, query)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all products: %w", err)
	}
	return data, nil
}

func (r *ProductRepo) GetProductByID(id string) (*domain.Product, error) {
	product := &domain.Product{}
	const query = "select * from products where id=$1"
	err := r.db.Get(product, query, id)
	if err != nil {
		return nil, fmt.Errorf("Error while getting product by id: %w", err)
	}
	return product, nil
}

func (r *ProductRepo) GetProductsByCategory(categoryID string) ([]domain.Product, error) {
	const query = "select * from products where category=$1"
	data := make([]domain.Product, 0)
	err := r.db.Select(&data, query, categoryID)
	if err != nil {
		return nil, fmt.Errorf("Error while getting products by category: %w", err)
	}
	return data, nil
}

func (r *ProductRepo) InsertProduct(product *domain.Product) error {
	const query = `
    insert into products (
      id, label, category, description,
      price, listed_by
    )
    values (
      :id, :label, :category, :description,
      :price, :listed_by
    )
  `
	_, err := r.db.NamedExec(query, product)
	if err != nil {
		return fmt.Errorf("Error while inserting product: %w", err)
	}
	return nil
}

func (r *ProductRepo) DeleteProduct(productID string) error {
	return nil
}
