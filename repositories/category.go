package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"lucy/models"
	"lucy/types"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type CategoryRepo struct {
	db *sqlx.DB
}

func NewCategoryRepository(db *sqlx.DB) CategoryRepo {
	return CategoryRepo{db}
}

func (r CategoryRepo) Insert(category *models.Category) error {
	_, err := r.db.NamedExec(
		`
      insert into categories(id, label, description, enabled)
      values(:id, :label, :description, :enabled)
    `,
		category,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting category: %w", err)
	}
	return nil
}

func (r CategoryRepo) GetCategory(filter Filter) (*models.Category, error) {
	category := &models.Category{}
	query := fmt.Sprintf("select * from category where %s=$1", filter.Field)
	err := r.db.Get(
		category,
		query,
		filter.Value,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrResourceNotFound
		}
		return nil, fmt.Errorf("Error while getting category by %s: %w", filter.Field, err)
	}
	return category, nil
}

func (r CategoryRepo) GetAll(callerAccountType types.AccountType) (*[]models.Category, error) {
	categories := []models.Category{}
	query := "select * from categories "
	if callerAccountType != types.AdminAccount {
		query += "where enabled=true"
	}
	err := r.db.Select(&categories, query)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all categories: %w", err)
	}
	return &categories, nil
}

func (r CategoryRepo) CountByLabel(label string) (int, error) {
	count := 0
	err := r.db.QueryRowx("select count(*) from categories where label=$1", label).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("Error while counting categories by label: %w", err)
	}
	return count, nil
}

func (r CategoryRepo) ToggleEnabled(ids []string, status bool) error {
	_, err := r.db.Exec(
		"update categories set enabled = $1 where id = ANY($2)",
    status,
		pq.StringArray(ids),
	)
	if err != nil {
		return fmt.Errorf("Error while toggling enabled status: %w", err)
	}
	return nil
}
