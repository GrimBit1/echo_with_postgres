package model

type User struct {
	FirstName string   `json:"first_name"`
	LastName  string   `json:"last_name"`
	Role      []string `json:"role"`
	JobTitle  string   `json:"title"`
}
type Error struct {
	Message string `json:"message"`
}
