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
	Label string `json:"label" db:"label"`
}

type Product struct {
	ID          string  `json:"id" db:"id"`
	Label       string  `json:"label" db:"label"`
	Brand       string  `json:"brand" db:"brand"`
	Category    string  `json:"category" db:"category"`
	Color       string  `json:"color" db:"color"`
	Description string  `json:"description" db:"description"`
	Image       string  `json:"image" db:"image"`
	ListedBy    string  `json:"listed_by" db:"listedby"`
	Price       float64 `json:"price" db:"price"`
	Size        string  `json:"size" db:"size"`
}

type Session struct {
	ID     string `json:"id" db:"id"`
	Valid  bool   `json:"valid" db:"valid"`
	UserID string `json:"user_id" db:"user_id"`
}

type AuthCode struct {
	ID           int    `db:"id"`
	Code         string `db:"code"`
	Used         bool   `db:"used"`
	GeneratedFor string `db:"generated_for"`
}
