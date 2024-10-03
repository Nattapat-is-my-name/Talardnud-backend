package Usecase

import (
	"fmt"
	"github.com/google/uuid"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
)

type ProviderUseCase struct {
	repo Interfaces.IProvider
}

func NewProviderUseCase(repo Interfaces.IProvider) *ProviderUseCase {
	return &ProviderUseCase{
		repo: repo,
	}
}

// CreateProvider creates a new provider and returns an error response if there's an issue.
func (uc *ProviderUseCase) CreateProvider(provider *entitiesDtos.MarketProviderRequest) (*entities.MarketProvider, *entitiesDtos.ErrorResponse) {
	// Check if provider with email already exists
	existingProvider, err := uc.repo.CheckProviderByEmail(provider.Email)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Error checking existing provider: %v", err),
		}
	}

	// If a provider with this email already exists, return an error
	if existingProvider != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Provider with this email already exists",
		}
	}

	// Continue with creating the new provider
	var providerEntities entities.MarketProvider
	id, err := uuid.NewUUID()
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to generate UUID: %v", err),
		}
	}

	providerEntities = entities.MarketProvider{
		ID:      id.String(),
		Name:    provider.Name,
		Phone:   provider.Phone,
		Email:   provider.Email,
		Address: provider.Address,
	}

	if err := uc.repo.CreateProvider(&providerEntities); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create provider: %v", err),
		}
	}

	return &providerEntities, nil
}

// UpdateProvider updates an existing provider and returns an error response if there's an issue.
func (uc *ProviderUseCase) UpdateProvider(provider *entities.MarketProvider) (*entities.MarketProvider, *entitiesDtos.ErrorResponse) {
	if err := uc.repo.UpdateProvider(provider); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to update provider: %v", err),
		}
	}

	return provider, nil
}

// GetProviderByID UseCase for getting a provider by ID
func (uc *ProviderUseCase) GetProviderByID(id string) (*entities.MarketProvider, *entitiesDtos.ErrorResponse) {
	provider, err := uc.repo.GetProviderByID(id)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to get provider: %v", err),
		}
	}
	if provider == nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "provider not found",
		}
	}
	return provider, nil
}

// DeleteProvider UseCase for deleting a provider by ID
func (uc *ProviderUseCase) DeleteProvider(id string) *entitiesDtos.ErrorResponse {
	if err := uc.repo.DeleteProvider(id); err != nil {
		return &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to delete provider: %v", err),
		}
	}
	return nil
}

// GetAllProviders UseCase for getting all providers
func (uc *ProviderUseCase) GetAllProviders() ([]*entities.MarketProvider, *entitiesDtos.ErrorResponse) {
	providers, err := uc.repo.GetAllProviders()
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("failed to retrieve providers: %v", err),
		}
	}
	return providers, nil
}
