// Services/dashboard_service.go
package Services

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"
	entities "tln-backend/Entities"
	"tln-backend/Repository"
)

type DashboardService struct {
	repo      *Repository.DashboardRepository
	scheduler *gocron.Scheduler
}

func NewDashboardService(repo *Repository.DashboardRepository) *DashboardService {
	scheduler := gocron.NewScheduler(time.UTC)
	service := &DashboardService{
		repo:      repo,
		scheduler: scheduler,
	}

	service.startScheduler()
	return service
}

func (s *DashboardService) startScheduler() {
	// Update stats every hour
	_, err := s.scheduler.Every(1).Hour().Do(func() {
		log.Println("Running scheduled dashboard update")
		if err := s.updateAllMarketStats(); err != nil {
			log.Printf("Scheduled update failed: %v", err)
		}
	})

	if err != nil {
		log.Printf("Failed to schedule dashboard updates: %v", err)
	}

	s.scheduler.StartAsync()
}

func (s *DashboardService) updateAllMarketStats() error {
	var marketIDs []string
	err := s.repo.GetDB().Model(&entities.Booking{}).
		Select("DISTINCT market_id").
		Pluck("market_id", &marketIDs).Error

	if err != nil {
		return fmt.Errorf("failed to get market IDs: %v", err)
	}

	for _, marketID := range marketIDs {
		if err := s.repo.UpdateDashboardStats(marketID); err != nil {
			log.Printf("Failed to update stats for market %s: %v", marketID, err)
		}
	}

	return nil
}
