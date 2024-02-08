package types

type UpdateUserInfoPayload struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Password    string `db:"password"`
	Description string `db:"description"`
	Town        string `db:"town"`
	Country     string `db:"country"`
}
