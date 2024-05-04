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
	*sqlx.DB
}

func (repo UserRepo) Insert(user *models.User) error {
	_, err := repo.NamedExec(
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

func (repo UserRepo) GetByID(id string) (*models.User, error) {
	user := &models.User{}
	err := repo.Get(
		user,
		"select * from users where id=$1",
		id,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, types.ErrResourceNotFound
		}
		return nil, fmt.Errorf("Error while getting user by id: %w", err)
	}
	return user, nil
}

func (repo UserRepo) GetAll() (*[]models.User, error) {
	users := []models.User{}
	err := repo.Select(
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
	err := repo.QueryRowx("select count(*) from users where phone=$1", phone).Scan(&count)
	if err != nil {
		return count, fmt.Errorf("Error while counting users by phone: %w", err)
	}
	return count, nil
}
