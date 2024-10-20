package dtos

type MarketRequest struct {
	ProviderID  string `json:"provider_id" validate:"required,uuid"`          // Required, UUID of the provider
	Name        string `json:"name" validate:"required"`                      // Required, name of the market
	Address     string `json:"address" validate:"required"`                   // Required, address of the market
	Description string `json:"description,omitempty"`                         // Optional, description of the market
	Image       string `json:"image,omitempty"`                               // Optional, URL or path to the market image
	OpenTime    string `json:"open_time" validate:"required,datetime=15:04"`  // Required, opening time in HH:mm format
	CloseTime   string `json:"close_time" validate:"required,datetime=15:04"` // Required, closing time in HH:mm format
	Latitude    string `json:"latitude,omitempty"`                            // Optional, latitude coordinate
	Longitude   string `json:"longitude,omitempty"`                           // Optional, longitude coordinate
}

type MarketEditRequest struct {
	ProviderID  string `json:"provider_id" validate:"required,uuid"`          // Required, UUID of the provider
	Name        string `json:"name" validate:"required"`                      // Required, name of the market
	Address     string `json:"address" validate:"required"`                   // Required, address of the market
	Description string `json:"description,omitempty"`                         // Optional, description of the market
	Image       string `json:"image,omitempty"`                               // Optional, URL or path to the market image
	LayoutImage string `json:"layout_image,omitempty"`                        // Optional, URL or path to the market layout image
	Phone       string `json:"phone,omitempty"`                               // Optional, phone number of the market
	OpenTime    string `json:"open_time" validate:"required,datetime=15:04"`  // Required, opening time in HH:mm format
	CloseTime   string `json:"close_time" validate:"required,datetime=15:04"` // Required, closing time in HH:mm format
	Latitude    string `json:"latitude,omitempty"`                            // Optional, latitude coordinate
	Longitude   string `json:"longitude,omitempty"`                           // Optional, longitude coordinate
}
