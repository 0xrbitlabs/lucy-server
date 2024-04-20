package types

type User struct {
	ID          string `json:"id" db:"id"`
	Type        string `json:"type" db:"type"`
	PhoneNumber string `json:"phone_number" db:"phone_number"`
	Password    string `json:"password" db:"password"`
	Username    string `json:"username" db:"username"`
	Description string `json:"description" db:"description"`
	Country     string `json:"country" db:"country"`
	Town        string `json:"town" db:"town"`
}
