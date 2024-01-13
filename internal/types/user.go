package types

type User struct {
	Id             string `json:"id" db:"id"`
	UserType       string `json:"user_type" db:"user_type"`
	PhoneNumber    string `json:"phone_number" db:"phone_number"`
	Password       string `json:"passord" db:"passord"`
	Name           string `json:"full_name" db:"full_name"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`
	Description    string `json:"descriptiton" db:"descriptiton"`
	Country        string `json:"country" db:"country"`
	Town           string `json:"town" db:"town"`
}
