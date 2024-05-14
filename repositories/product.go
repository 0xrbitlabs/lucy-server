package repositories

import (
	"fmt"
	"lucy/models"
	"lucy/types"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
        description, price, image, enabled
      )
      values(
        :id, :owner, :category_id, :label,
        :description, :price, :image, :enabled
      )
    `,
		data,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting product: %w", err)
	}
	return nil
}

func (r ProductRepo) GetAll(user *models.User) (*[]models.Product, error) {
	var err error
	products := []models.Product{}
	switch user.AccountType {
	case types.AdminAccount:
		err = r.db.Select(
			&products,
			"select * from products",
		)
	case types.SellerAccount:
		err = r.db.Select(
			&products,
			"select * from products where owner=$1",
			user.ID,
		)
	}
	if err != nil {
		return nil, fmt.Errorf("Error while retrieving all products: %w", err)
	}
	return &products, nil
}

func (r ProductRepo) ToggleStatus(ids []string, status bool) error {
	_, err := r.db.Exec(
		"update products set enabled = $1 where id = ANY($2)",
		status,
		pq.StringArray(ids),
	)
	if err != nil {
		return fmt.Errorf("Error while toggling product status: %w", err)
	}
	return nil
}
