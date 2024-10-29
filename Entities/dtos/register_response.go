package dtos

type RegisterResponse struct {
	ID        string `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}
