package dtos

type MarketProviderRequest struct {
	Username string `json:"username" validate:"required"`    // Required, username of the provider
	Password string `json:"password" validate:"required"`    // Required, password of the provider
	Phone    string `json:"phone" validate:"required"`       // Required, phone number of the provider
	Email    string `json:"email" validate:"required,email"` // Required, email address of the provider
}
