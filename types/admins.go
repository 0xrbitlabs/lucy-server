package types

type Admin struct {
	ID       int    `json:"id" db:"id"`
	Username string `json:"username" db:"username"`
	Password string `json:"password" db:"password"`
	IsSuper  bool   `json:"is_super" db:"is_super"`
}
