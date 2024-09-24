package entities

type LoginRequest struct {
	UsernameOrEmail string `json:"username_or_email" binding:"required"` // Accepts either username or email
	Password        string `json:"password" binding:"required,min=8"`    // Password must be provided and at least 8 characters long
}

type LoginResponse struct {
	AccessToken string `json:"access_token"`
	VendorID    string `json:"vendor_id"`
}
