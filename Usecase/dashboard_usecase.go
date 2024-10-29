package Usecase

import (
	"fmt"
	"log"
	entities "tln-backend/Entities"
	"tln-backend/Repository"
	"tln-backend/Services"
)

type DashboardUseCase struct {
	repo    *Repository.DashboardRepository
	service *Services.DashboardService
}

func NewDashboardUseCase(repo *Repository.DashboardRepository, service *Services.DashboardService) *DashboardUseCase {
	return &DashboardUseCase{
		repo:    repo,
		service: service,
	}
}
func (uc *DashboardUseCase) updateAllMarketsStats() error {
	// Get all market IDs
	var marketIDs []string
	err := uc.repo.GetDB().Model(&entities.Market{}).
		Select("id").
		Pluck("id", &marketIDs).Error
	if err != nil {
		return fmt.Errorf("failed to get market IDs: %v", err)
	}

	// Update stats for each market
	for _, marketID := range marketIDs {
		if err := uc.repo.UpdateDashboardStats(marketID); err != nil {
			log.Printf("Warning: Failed to update stats for market %s: %v", marketID, err)
		}
	}

	return nil
}

func (uc *DashboardUseCase) GetWeeklyData(marketID string) (*entities.DashboardResponse, error) {
	// Update stats for the specific market before fetching weekly data
	err := uc.repo.UpdateDashboardStats(marketID)
	if err != nil {
		return nil, fmt.Errorf("failed to update market stats: %v", err)
	}

	// Fetch weekly stats for the given market ID
	weeklyStats, err := uc.repo.GetWeeklyStats(marketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get weekly dashboard data: %v", err)
	}

	return &entities.DashboardResponse{
		Stats: weeklyStats, // Assigns the retrieved weekly stats directly
	}, nil
}

// If you need to get stats for a single market, add this method
func (uc *DashboardUseCase) GetSingleMarketStats(marketID string) (*entities.DashboardResponse, error) {
	// Update stats for this market
	if err := uc.repo.UpdateDashboardStats(marketID); err != nil {
		return nil, fmt.Errorf("failed to update market stats: %v", err)
	}

	Stats, err := uc.repo.GetDashboardData(marketID)
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard data: %v", err)
	}

	return Stats, nil
}
