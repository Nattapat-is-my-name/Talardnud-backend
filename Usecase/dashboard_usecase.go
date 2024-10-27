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

func (uc *DashboardUseCase) GetDashboardData(marketID string) (*entities.DashboardResponse, error) {
	// First ensure stats are updated for all markets
	err := uc.updateAllMarketsStats()
	if err != nil {
		return nil, fmt.Errorf("failed to update market stats: %v", err)
	}

	// Get stats for all markets
	stats, err := uc.repo.GetAllMarketsDashboardStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard data: %v", err)
	}

	return &entities.DashboardResponse{
		Stats: stats, // Now this will work because Stats is a slice
	}, nil
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

// If you need to get stats for a single market, add this method
func (uc *DashboardUseCase) GetSingleMarketStats(marketID string) (*entities.MarketDashboardStats, error) {
	// Update stats for this market
	if err := uc.repo.UpdateDashboardStats(marketID); err != nil {
		return nil, fmt.Errorf("failed to update market stats: %v", err)
	}

	// Get all stats
	allStats, err := uc.repo.GetAllMarketsDashboardStats()
	if err != nil {
		return nil, fmt.Errorf("failed to get dashboard data: %v", err)
	}

	// Find the specific market's stats
	for _, stats := range allStats {
		if stats.MarketID == marketID {
			return &stats, nil
		}
	}

	return nil, fmt.Errorf("market not found")
}
