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

func (r *ProductCategoryRepo) InsertCreationRequest(req *models.ProductCategoryCreationRequest) error {
	const query = `
    insert into product_category_creation_requests (
      id, requester, label, description
    )
    values (
      :id, :requester, :label, :description
    )
  `
	_, err := r.db.NamedExec(query, req)
	if err != nil {
		return fmt.Errorf("Error while inserting request creation: %w", err)
	}
	return nil
}

func (r *ProductCategoryRepo) GetAllProductCategoryCreationRequestsByUserAccountType(
	user *models.User,
) ([]models.ProductCategoryCreationRequest, error) {
	data := make([]models.ProductCategoryCreationRequest, 0)
	var err error
	if user.AccountType == "seller" {
		err = r.db.Select(&data, "select * from product_category_creation_requests where requester=$1", user.ID)
	} else {
		err = r.db.Select(&data, "select * from product_category_creation_requests")
	}
	if err != nil {
		return nil, fmt.Errorf("Error while getting product category creation requests: %w", err)
	}
	return data, nil
}

func (r *ProductCategoryRepo) SetRequestStatus(id, status string) error {
	const query = "update product_category_creation_requests set status='$2' where id=$1"
	_, err := r.db.Exec(query, id, status)
	if err != nil {
		return fmt.Errorf("Error while setting request status to rejected: %w", err)
	}
	return nil
}

func (r *ProductCategoryRepo) GetProductCategoryCreationRequestByID(id string) (*models.ProductCategoryCreationRequest, error) {
	const query = "select * from product_category_creation_requests where id=$1"
	req := &models.ProductCategoryCreationRequest{}
	err := r.db.Get(req, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("Error while getting product category creation request by id: %w", err)
	}
	return req, nil
}
