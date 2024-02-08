package store

import (
	"database/sql"
	"fmt"
	"github.com/blockloop/scan/v2"
	"server/internal/types"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users {
	return &Users{
		db: db,
	}
}

func (u *Users) Insert(user *types.User) error {
	_, err := u.db.Exec(
		`
      insert into users (
        id, user_type, phone_number, password,
        name, description, country, town
      )
      values (
        ?, ?, ?, ?, ?, ?, ?, ?
      )
    `,
		user.Id, user.UserType, user.PhoneNumber, user.Password,
		user.Name, user.Description, user.Country, user.Town,
	)
	if err != nil {
		return fmt.Errorf("Failed to insert new user: %w", err)
	}
	return nil
}

func (u *Users) CountByPhoneNumber(phoneNumber string) (int, error) {
	count := 0
	err := u.db.QueryRow(
		"select count(*) from users where phone_number=?",
		phoneNumber,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("Error while counting phone numbers: %w", err)
	}
	return count, nil
}

func (u *Users) GetByPhoneNumber(phoneNumber string) (*types.User, error) {
	user := new(types.User)
	rows, err := u.db.Query("select * from users where phone_number=?", phoneNumber)
	if err != nil {
		return nil, fmt.Errorf("Error while querying user: %w", err)
	}
	defer rows.Close()
	err = scan.Row(user, rows)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrUserNotFound
		}
		return nil, fmt.Errorf("Error while scanning row into user: %w", err)
	}
	return user, nil
}

func (u *Users) UpdateInfo(data *types.UpdateUserInfoPayload) error {
	_, err := u.db.Exec(
		`
      update users set
      name=:name, password=:password, description=:description,
      country=:country, town=:town where id=:id
    `,
		data,
	)
	if err != nil {
		return fmt.Errorf("Error while updating user info: %w", err)
	}
	return nil
}
