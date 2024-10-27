package domain

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

type ProductCategory struct {
	ID          string `json:"id" db:"id"`
	Label       string `json:"label" db:"label"`
	Description string `json:"description" db:"description"`
	Active      bool   `json:"active" db:"active"`
}

type Product struct {
	ID          string  `json:"id" db:"id"`
	Label       string  `json:"label" db:"label"`
	Description string  `json:"description" db:"description"`
	Images      string  `json:"image" db:"images"`
	Price       float64 `json:"price" db:"price"`
	ListedBy    string  `json:"listed_by" db:"listedby"`
}
