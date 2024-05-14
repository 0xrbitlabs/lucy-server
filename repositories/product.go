package repositories

import (
	"fmt"
	"lucy/models"

	"github.com/jmoiron/sqlx"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) ProductRepo {
	return ProductRepo{db: db}
}

func (r ProductRepo) Insert(data *models.Product) error {
	_, err := r.db.NamedExec(
		`
      insert into products(
        id, owner, category_id, label,
        description, price, image
      )
      values(
        :id, :owner, :category_id, :label,
        :description, :price, :image
      )
    `,
		data,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting product: %w", err)
	}
	return nil
}

func (r ProductRepo) GetAll() (*[]models.Product, error) {
	products := []models.Product{}
	err := r.db.Select(
		&products,
		"select * from products",
	)
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving all products: %w", err)
	}
	return &products, nil
}
