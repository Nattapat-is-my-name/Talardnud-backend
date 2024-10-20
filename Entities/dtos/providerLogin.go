package dtos

type ProviderLoginRequest struct {
	Username string `json:"username" validate:"required,username"`
	Password string `json:"password" validate:"required,min=8"`
}

// ProviderLoginResponse represents the response body for successful provider login
type ProviderLoginResponse struct {
	AccessToken string `json:"access_token"`
	ProviderID  string `json:"provider_id"`
}
