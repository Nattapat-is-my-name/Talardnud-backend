package dtos

type RegisterRequest struct {
	Username    string `json:"username" binding:"required,min=3,max=50"`      // Required, min 3, max 50 characters
	Password    string `json:"password" binding:"required,min=8"`             // Required, min 8 characters for password
	Email       string `json:"email" binding:"required,email"`                // Required, must be a valid email format
	PhoneNumber string `json:"phone_number" binding:"required,min=10,max=15"` // Required, adjust based on the expected format=
}