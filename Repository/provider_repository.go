package Repository

import (
	"errors"
	"gorm.io/gorm"
	"log"
	entities "tln-backend/Entities"
)

type ProviderRepository struct {
	db *gorm.DB
}

func NewProviderRepository(db *gorm.DB) *ProviderRepository {
	return &ProviderRepository{db: db}
}

// CreateProvider inserts a new provider into the database.
//func (pr *ProviderRepository) CreateProvider(provider *entities.MarketProvider) error {
//	if err := pr.db.Create(provider).Error; err != nil {
//		log.Printf("Error creating provider: %v", err)
//		return err
//	}
//	return nil
//}

// UpdateProvider updates an existing provider in the database.
func (pr *ProviderRepository) UpdateProvider(provider *entities.MarketProvider) error {
	if err := pr.db.Save(provider).Error; err != nil {
		log.Printf("Error updating provider: %v", err)
		return err
	}
	return nil
}

// GetProviderByID retrieves a provider by its ID.
func (pr *ProviderRepository) GetProviderByID(id string) (*entities.MarketProvider, error) {
	var provider entities.MarketProvider
	if err := pr.db.First(&provider, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Provider not found with ID: %s", id)
			return nil, nil
		}
		log.Printf("Error retrieving provider by ID: %v", err)
		return nil, err
	}
	return &provider, nil
}

// DeleteProvider deletes a provider by its ID.
func (pr *ProviderRepository) DeleteProvider(id string) error {
	if err := pr.db.Delete(&entities.MarketProvider{}, "id = ?", id).Error; err != nil {
		log.Printf("Error deleting provider with ID: %s", id)
		return err
	}
	return nil
}

// GetAllProviders retrieves all providers from the database.
func (pr *ProviderRepository) GetAllProviders() ([]*entities.MarketProvider, error) {
	var providers []*entities.MarketProvider
	if err := pr.db.Find(&providers).Error; err != nil {
		log.Printf("Error retrieving all providers: %v", err)
		return nil, err
	}
	return providers, nil
}

// CheckProviderByEmail checks if a provider exists by its email.
func (pr *ProviderRepository) CheckProviderByUsername(username string) (*entities.MarketProvider, error) {
	var provider entities.MarketProvider
	if err := pr.db.First(&provider, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			log.Printf("Provider not found with email: %s", username)
			return nil, nil
		}
		log.Printf("Error retrieving provider by email: %v", err)
		return nil, err
	}
	return &provider, nil

}
