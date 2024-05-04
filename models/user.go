package models

type User struct {
	ID          string `json:"id" db:"id"`
	Username    string `json:"username" db:"username"`
	Phone       string `json:"phone" db:"phone"`
	Password    string `json:"password" db:"password"`
	AccountType string `json:"account_type" db:"account_type"`
	Description string `json:"description" db:"description"`
	Country     string `json:"country" db:"country"`
	Town        string `json:"town" db:"town"`
}
