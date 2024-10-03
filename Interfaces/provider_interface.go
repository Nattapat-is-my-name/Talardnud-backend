package Interfaces

import entities "tln-backend/Entities"

type IProvider interface {
	CreateProvider(provider *entities.MarketProvider) error
	UpdateProvider(provider *entities.MarketProvider) error
	GetProviderByID(id string) (*entities.MarketProvider, error)
	DeleteProvider(id string) error
	GetAllProviders() ([]*entities.MarketProvider, error)
	CheckProviderByEmail(email string) (*entities.MarketProvider, error)
}
