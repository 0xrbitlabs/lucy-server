package repository

import (
	"github.com/jmoiron/sqlx"
	"github.com/joseph0x45/lucy/domain"
)

type ProductRepo struct {
	db *sqlx.DB
}

func NewProductRepo(db *sqlx.DB) *ProductRepo {
	return &ProductRepo{db}
}

func (r *ProductRepo) InsertProductCategory(pc *domain.ProductCategory) error {
	const query = `
    insert into product_category(

    )
    values (

    )
  `
	return nil
}
