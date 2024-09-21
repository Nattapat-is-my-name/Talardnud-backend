package Interfaces

import (
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type IMarket interface {
	CreateMarket(mae *entities.Market) error
	GetMarketByName(name string) (*entities.Market, *entitiesDtos.ErrorResponse)
	GetProviderByID(providerID string) (*entities.MarketProvider, *entitiesDtos.ErrorResponse)
	GetMarketWithProviderByID(marketID string) (*entities.Market, *entitiesDtos.ErrorResponse)
}
