package dtos

type MarketProviderRequest struct {
	Name    string `json:"name" validate:"required"`        // Required, name of the market provider
	Phone   string `json:"phone" validate:"required"`       // Required, phone number of the provider
	Email   string `json:"email" validate:"required,email"` // Required, email address of the provider
	Address string `json:"address" validate:"required"`     // Required, physical address of the provider
}
