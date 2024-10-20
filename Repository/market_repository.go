package Repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"log"
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
func (repo *MarketRepository) GetMarkets() ([]entities.Market, *entitiesDtos.ErrorResponse) {
	var markets []entities.Market

	// Add Debug logging to see actual SQL query
	err := repo.db.Debug().Find(&markets).Error
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Database error: " + err.Error(),
		}
	}

	// Log the number of markets retrieved for debugging purposes
	fmt.Printf("Markets found: %d\n", len(markets))

	return markets, nil
}

// get market by id
func (repo *MarketRepository) GetMarketByID(marketID string) ([]entities.Market, *entitiesDtos.ErrorResponse) {
	var market []entities.Market
	err := repo.db.Where("id = ?", marketID).First(&market).Error

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

	err = repo.db.Preload("Slots").Where("id = ?", marketID).First(&market).Error
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Database error: " + err.Error(),
		}

	}

	return market, nil
}

func (repo *MarketRepository) EditMarket(marketID string, marketReq *entitiesDtos.MarketEditRequest) (*entities.Market, *entitiesDtos.ErrorResponse) {
	var market entities.Market
	err := repo.db.Where("id = ?", marketID).First(&market).Error
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

	market.Name = marketReq.Name
	market.Description = marketReq.Description
	market.Address = marketReq.Address
	market.Phone = marketReq.Phone
	market.Image = marketReq.Image
	market.LayoutImage = marketReq.LayoutImage
	market.OpenTime = marketReq.OpenTime
	market.CloseTime = marketReq.CloseTime
	market.Latitude = marketReq.Latitude
	market.Longitude = marketReq.Longitude

	err = repo.db.Save(&market).Error
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to edit market: " + err.Error(),
		}
	}

	return &market, nil
}
func (repo *MarketRepository) GetMarketByProviderID(providerID string) ([]entities.Market, *entitiesDtos.ErrorResponse) {
	var markets []entities.Market

	log.Printf("Fetching markets for provider ID: %s", providerID)

	err := repo.db.Preload("Slots").Preload("Provider").Where("provider_id = ?", providerID).Find(&markets).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("No markets found for provider ID: %s", providerID)
			return nil, &entitiesDtos.ErrorResponse{
				Code:    404,
				Message: "Market not found",
			}
		}

		log.Printf("Database error when fetching markets: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Database error: " + err.Error(),
		}
	}

	log.Printf("Found %d markets for provider ID: %s", len(markets), providerID)

	for i, market := range markets {
		log.Printf("Market %d: ID=%s, Name=%s", i+1, market.ID, market.Name)
		log.Printf("  Slots count: %d", len(market.Slots))
		for j, slot := range market.Slots {
			log.Printf("    Slot %d: ID=%s, Date=%s", j+1, slot.ID, slot.Date)
		}
	}

	return markets, nil
}
