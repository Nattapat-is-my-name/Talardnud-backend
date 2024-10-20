package Usecase

import (
	"github.com/google/uuid"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
)

type MarketUseCase struct {
	repo Interfaces.IMarket
}

func NewMarketUseCase(repo Interfaces.IMarket) *MarketUseCase {
	return &MarketUseCase{
		repo: repo,
	}

}

func (uc *MarketUseCase) CreateMarket(marketReq *entitiesDtos.MarketRequest) (*entities.Market, *entitiesDtos.ErrorResponse) {
	// Check if the provider exists
	_, errRes := uc.repo.GetProviderByID(marketReq.ProviderID)
	if errRes != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "Provider not found",
		}
	}

	// Check if a market with the same name already exists
	existingMarket, errRes := uc.repo.GetMarketByName(marketReq.Name)
	if errRes != nil && errRes.Code != 404 { // If error is not "not found", return it
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to check market existence: " + errRes.Message,
		}
	}

	if existingMarket != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Market already exists",
		}
	}

	// Map the MarketRequest to Market entity
	marketEntity := entities.Market{
		ID:          uuid.New().String(),
		ProviderID:  marketReq.ProviderID,
		Name:        marketReq.Name,
		Address:     marketReq.Address,
		Description: marketReq.Description,
		Image:       marketReq.Image,
		OpenTime:    marketReq.OpenTime,
		CloseTime:   marketReq.CloseTime,
		Latitude:    marketReq.Latitude,
		Longitude:   marketReq.Longitude,
	}

	// Save the market entity to the database
	err := uc.repo.CreateMarket(&marketEntity)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to create market: " + err.Error(),
		}
	}

	// Retrieve the market with full provider details
	createdMarket, errRes := uc.repo.GetMarketWithProviderByID(marketEntity.ID)
	if errRes != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve market details: " + errRes.Message,
		}
	}

	// Return the created market with provider details
	return createdMarket, nil
}

// edit market
func (uc *MarketUseCase) EditMarket(marketID string, marketReq *entitiesDtos.MarketEditRequest) (*entities.Market, *entitiesDtos.ErrorResponse) {
	// Check if the market exists
	_, errRes := uc.repo.GetMarketByID(marketID)
	if errRes != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "Market not found",
		}
	}

	// Check if the provider exists
	_, errRes = uc.repo.GetProviderByID(marketReq.ProviderID)
	if errRes != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "Provider not found",
		}
	}

	// Call the EditMarket function in the repository with marketID and marketReq
	_, err := uc.repo.EditMarket(marketID, marketReq)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to edit market: " + err.Error(),
		}
	}

	// Retrieve the market with full provider details
	editedMarket, errRes := uc.repo.GetMarketWithProviderByID(marketID)
	if errRes != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve market details: " + errRes.Message,
		}
	}

	// Return the edited market with provider details
	return editedMarket, nil
}

func (uc *MarketUseCase) GetMarket() ([]entities.Market, *entitiesDtos.ErrorResponse) {
	markets, err := uc.repo.GetMarkets()
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve markets: " + err.Error(),
		}

	}

	return markets, nil

}

func (uc *MarketUseCase) GetMarketByID(marketID string) ([]entities.Market, *entitiesDtos.ErrorResponse) {
	market, errRes := uc.repo.GetMarketByID(marketID)
	if errRes != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve market: " + errRes.Message,
		}
	}

	return market, nil
}

func (uc *MarketUseCase) GetMarketByProviderID(providerID string) ([]entities.Market, *entitiesDtos.ErrorResponse) {
	market, errRes := uc.repo.GetMarketByProviderID(providerID)
	if errRes != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to retrieve market: " + errRes.Message,
		}
	}

	return market, nil
}
