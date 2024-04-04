package models

type User struct {
	ID          string `db:"id"`
	Type        string `db:"user_type"`
	PhoneNumber string `db:"phone_number"`
	Password    string `db:"password"`
	Name        string `db:"name"`
	Description string `db:"description"`
	Country     string `db:"country"`
	Town        string `db:"town"`
}
