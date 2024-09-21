package dtos

type OAuthRequest struct {
	ApplicationKey    string `json:"applicationKey"`
	ApplicationSecret string `json:"applicationSecret"`
}

type OAuthResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"status"`
	Data struct {
		UUID        string `json:"uuid"`
		AccessToken string `json:"accessToken"`
		ExpiresIn   int    `json:"expiresIn"`
		TokenType   string `json:"tokenType"`
	} `json:"data"`
}
