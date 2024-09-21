package Repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type MarketRepository struct {
	db *gorm.DB
}

func NewMarketRepository(db *gorm.DB) *MarketRepository {
	return &MarketRepository{db: db}
}

func (repo *MarketRepository) CreateMarket(market *entities.Market) error {
	var provider entities.MarketProvider
	err := repo.db.Where("id = ?", market.ProviderID).First(&provider).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("provider not found for the given provider_id: %s", market.ProviderID)
		}
		return err
	}

	// Now, create the market record
	return repo.db.Create(market).Error
}

func (repo *MarketRepository) GetMarketByName(name string) (*entities.Market, *entitiesDtos.ErrorResponse) {
	var market entities.Market
	err := repo.db.Where("name = ?", name).First(&market).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &entitiesDtos.ErrorResponse{
				Code:    404,
				Message: "Market not found",
			}
		}

		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Database error: " + err.Error(),
		}
	}

	return &market, nil
}

// GetProviderByID retrieves a provider by its ID.
func (repo *MarketRepository) GetProviderByID(providerID string) (*entities.MarketProvider, *entitiesDtos.ErrorResponse) {
	var provider entities.MarketProvider
	err := repo.db.Where("id = ?", providerID).First(&provider).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &entitiesDtos.ErrorResponse{
				Code:    404,
				Message: "Provider not found",
			}
		}
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Database error: " + err.Error(),
		}
	}

	return &provider, nil
}
func (repo *MarketRepository) GetMarketWithProviderByID(marketID string) (*entities.Market, *entitiesDtos.ErrorResponse) {
	var market entities.Market
	err := repo.db.Preload("Provider").Where("id = ?", marketID).First(&market).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, &entitiesDtos.ErrorResponse{
				Code:    404,
				Message: "Market not found",
			}
		}

		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Database error: " + err.Error(),
		}
	}

	return &market, nil
}
