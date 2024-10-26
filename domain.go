package main

type AccountType string

const (
	SellerAccountType  AccountType = "seller"
	RegularAccountType AccountType = "regular"
	AdminAccountType   AccountType = "admin"
)

type User struct {
	ID          string `json:"id" db:"id"`
	Phone       string `json:"phone" db:"phone"`
	Username    string `json:"username" db:"username"`
	Password    string `json:"password" db:"password"`
	AccountType string `json:"account_type" db:"account_type"`
}
