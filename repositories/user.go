package repositories

import (
	"database/sql"
	"errors"
	"fmt"
	"lucy/models"
	"lucy/types"

	"github.com/jmoiron/sqlx"
)

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepo {
	return UserRepo{db: db}
}

type Filter struct {
	Field string
	Value string
}

func (repo UserRepo) Insert(user *models.User) error {
	_, err := repo.db.NamedExec(
		`
      insert into users(
        id, username, phone, password, account_type, description, country, town
      )
      values(
        :id, :username, :phone, :password, :account_type, :description, :country, :town
      )
    `,
		user,
	)
	if err != nil {
		return fmt.Errorf("Error while inserting user: %w", err)
	}
	return nil
}

func (repo UserRepo) GetUser(filter Filter) (*models.User, error) {
	user := &models.User{}
	query := fmt.Sprintf("select * from users where %s=$1", filter.Field)
	err := repo.db.Get(
		user,
		query,
		filter.Value,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrResourceNotFound
		}
		return nil, fmt.Errorf("Error while getting user by %s: %w", filter.Field, err)
	}
	return user, nil
}

func (repo UserRepo) GetAll() (*[]models.User, error) {
	users := []models.User{}
	err := repo.db.Select(
		&users,
		"select * from users",
	)
	if err != nil {
		return nil, fmt.Errorf("Error while getting all users: %w", err)
	}
	return &users, nil
}

func (repo UserRepo) CountByPhone(phone string) (int, error) {
	count := 0
	err := repo.db.QueryRowx("select count(*) from users where phone=$1", phone).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("Error while counting users by phone: %w", err)
	}
	return count, nil
}

func (repo UserRepo) UpdatePassword(userId, password string) error {
	_, err := repo.db.Exec("update users set password=$1 where id=$2", password, userId)
	if err != nil {
		return fmt.Errorf("Error while updating user password: %w", err)
	}
	return nil
}
