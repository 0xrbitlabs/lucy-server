package types

type User struct {
	Id             string `json:"id" db:"id"`
	PhoneNumber    string `json:"phone_number" db:"phone_number"`
	FullName       string `json:"full_name" db:"full_name"`
	ProfilePicture string `json:"profile_picture" db:"profile_picture"`
}
