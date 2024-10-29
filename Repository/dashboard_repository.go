// Repository/dashboard_repository.go
package Repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
	"time"
	entities "tln-backend/Entities"
)

type DashboardRepository struct {
	db *gorm.DB
}

func NewDashboardRepository(db *gorm.DB) *DashboardRepository {
	return &DashboardRepository{db: db}
}

func (repo *DashboardRepository) GetDB() *gorm.DB {
	return repo.db
}
func (repo *DashboardRepository) UpdateDashboardStats(marketID string) error {
	query := `
        WITH distinct_dates AS (
            SELECT DISTINCT 
                market_id,
                DATE(booking_date) as date
            FROM bookings
        ),
        markets AS (
            SELECT id AS market_id FROM markets
        ),
        market_dates AS (
            SELECT DISTINCT d.market_id, d.date
            FROM distinct_dates d
        ),
        booking_stats AS (
            SELECT 
                b.market_id,
                DATE(b.booking_date) as booking_date,
                COUNT(*) as total_bookings,
                COUNT(CASE WHEN b.status = 'completed' THEN 1 END) as completed_bookings,
                COUNT(CASE WHEN b.status = 'cancelled' THEN 1 END) as cancelled_bookings,
                COALESCE(SUM(CASE WHEN b.status = 'completed' AND p.status = 'completed' 
                    THEN p.price ELSE 0 END), 0) as total_revenue,
                COALESCE(
                    (COUNT(CASE WHEN b.status = 'completed' THEN 1 END) * 100.0 / 
                    NULLIF(COUNT(*), 0)),
                    0
                ) as occupancy_rate
            FROM bookings b
            LEFT JOIN payments p ON b.id = p.booking_id
            GROUP BY b.market_id, DATE(b.booking_date)
        ),
        zone_stats AS (
            SELECT 
                b.market_id,
                DATE(b.booking_date) as booking_date,
                s.zone,
                COUNT(*) as zone_bookings,
                COUNT(CASE WHEN b.status = 'completed' THEN 1 END) as completed_in_zone,
                COALESCE(
                    COUNT(CASE WHEN b.status = 'completed' THEN 1 END) * 100.0 / 
                    NULLIF(COUNT(*), 0),
                    0
                ) as zone_occupancy
            FROM bookings b
            JOIN slots s ON b.slot_id = s.id
            GROUP BY b.market_id, DATE(b.booking_date), s.zone
        ),
        top_zones AS (
            SELECT DISTINCT ON (market_id, booking_date)
                market_id,
                booking_date,
                zone as top_zone,
                zone_occupancy as top_zone_occupancy
            FROM zone_stats
            ORDER BY market_id, booking_date, completed_in_zone DESC
        ),
        prev_day_stats AS (
            SELECT 
                md.market_id,
                md.date,
                LAG(bs.completed_bookings, 1) OVER (PARTITION BY md.market_id ORDER BY md.date) as prev_bookings,
                LAG(bs.total_revenue, 1) OVER (PARTITION BY md.market_id ORDER BY md.date) as prev_revenue
            FROM market_dates md
            LEFT JOIN booking_stats bs ON md.market_id = bs.market_id AND md.date = bs.booking_date
        )
        INSERT INTO market_dashboard_stats (
            market_id,
            date,
            total_bookings,
            total_confirm_bookings,
            total_cancel_bookings,
            booking_growth,
            total_revenue,
            revenue_growth,
            occupancy_rate,
            top_zone,
            top_zone_occupancy,
            created_at
        )
        SELECT 
            md.market_id,
            md.date,
            COALESCE(bs.total_bookings, 0),
            COALESCE(bs.completed_bookings, 0),
            COALESCE(bs.cancelled_bookings, 0),
            CASE 
                WHEN pds.prev_bookings > 0 THEN 
                    ((COALESCE(bs.completed_bookings, 0) - pds.prev_bookings) * 100.0 / pds.prev_bookings)
                ELSE 0 
            END as booking_growth,
            COALESCE(bs.total_revenue, 0),
            CASE 
                WHEN pds.prev_revenue > 0 THEN 
                    ((COALESCE(bs.total_revenue, 0) - pds.prev_revenue) * 100.0 / pds.prev_revenue)
                ELSE 0 
            END as revenue_growth,
            COALESCE(bs.occupancy_rate, 0),
            COALESCE(tz.top_zone, ''),
            COALESCE(tz.top_zone_occupancy, 0),
            NOW()
        FROM market_dates md
        LEFT JOIN booking_stats bs ON 
            md.market_id = bs.market_id AND 
            md.date = bs.booking_date
        LEFT JOIN top_zones tz ON 
            md.market_id = tz.market_id AND 
            md.date = tz.booking_date
        LEFT JOIN prev_day_stats pds ON 
            md.market_id = pds.market_id AND 
            md.date = pds.date
        ON CONFLICT (market_id, date) 
        DO UPDATE SET 
            total_bookings = EXCLUDED.total_bookings,
            total_confirm_bookings = EXCLUDED.total_confirm_bookings,
            total_cancel_bookings = EXCLUDED.total_cancel_bookings,
            booking_growth = EXCLUDED.booking_growth,
            total_revenue = EXCLUDED.total_revenue,
            revenue_growth = EXCLUDED.revenue_growth,
            occupancy_rate = EXCLUDED.occupancy_rate,
            top_zone = EXCLUDED.top_zone,
            top_zone_occupancy = EXCLUDED.top_zone_occupancy,
            created_at = EXCLUDED.created_at
    `

	err := repo.db.Exec(query).Error
	if err != nil {
		return fmt.Errorf("failed to update stats: %v", err)
	}

	log.Printf("Successfully updated stats for all markets and their booking dates")
	return nil
}

func (repo *DashboardRepository) GetDashboardData(marketID string) (*entities.DashboardResponse, error) {
	// First try to update stats
	if err := repo.UpdateDashboardStats(marketID); err != nil {
		log.Printf("Warning: Failed to update stats: %v", err)
	}

	var stats entities.MarketDashboardStats
	currentDate := time.Now().Truncate(24 * time.Hour)

	err := repo.db.Where("market_id = ? AND date = ?",
		marketID,
		currentDate,
	).First(&stats).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Initialize with zero values if no record exists
			stats = entities.MarketDashboardStats{
				MarketID:      marketID,
				Date:          currentDate,
				TotalBookings: 0,
				BookingGrowth: 0,
				TotalRevenue:  0,
				RevenueGrowth: 0,
				CreatedAt:     time.Now(),
			}

			if err := repo.db.Create(&stats).Error; err != nil {
				return nil, fmt.Errorf("failed to create initial stats: %v", err)
			}
		} else {
			return nil, fmt.Errorf("failed to get dashboard data: %v", err)
		}
	}

	return &entities.DashboardResponse{
		Stats: []entities.MarketDashboardStats{stats},
	}, nil

}

func (repo *DashboardRepository) GetAllMarketsDashboardStats() ([]entities.MarketDashboardStats, error) {
	// Get all markets first
	var marketIDs []string
	err := repo.db.Model(&entities.Market{}).
		Select("id").
		Pluck("id", &marketIDs).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get market IDs: %v", err)
	}

	// Define date range
	endDate := time.Now().Truncate(24 * time.Hour)
	startDate := endDate.AddDate(0, 0, -30) // Last 30 days

	query := `
        WITH dates AS (
            SELECT generate_series(
                ?::date, 
                ?::date, 
                '1 day'::interval
            )::date AS date
        ),
        markets AS (
            SELECT id AS market_id FROM markets
        ),
        market_dates AS (
            SELECT m.market_id, d.date
            FROM markets m
            CROSS JOIN dates d
        ),
        booking_stats AS (
            SELECT 
                b.market_id,
                DATE(b.booking_date) as booking_date,
                COUNT(*) as total_bookings,
                COUNT(CASE WHEN b.status = 'completed' THEN 1 END) as completed_bookings,
                COUNT(CASE WHEN b.status = 'cancelled' THEN 1 END) as cancelled_bookings,
                COALESCE(SUM(CASE WHEN b.status = 'completed' AND p.status = 'completed' 
                    THEN p.price ELSE 0 END), 0) as total_revenue,
                COALESCE(
                    (COUNT(CASE WHEN b.status = 'completed' THEN 1 END) * 100.0 / 
                    NULLIF(COUNT(*), 0)),
                    0
                ) as occupancy_rate
            FROM bookings b
            LEFT JOIN payments p ON b.id = p.booking_id
            WHERE DATE(b.booking_date) BETWEEN ? AND ?
            GROUP BY b.market_id, DATE(b.booking_date)
        )
        SELECT 
            md.market_id,
            md.date,
            COALESCE(bs.total_bookings, 0) as total_bookings,
            COALESCE(bs.completed_bookings, 0) as total_confirm_bookings,
            COALESCE(bs.cancelled_bookings, 0) as total_cancel_bookings,
            COALESCE(bs.total_revenue, 0) as total_revenue,
            COALESCE(bs.occupancy_rate, 0) as occupancy_rate,
            COALESCE(
                (SELECT zone 
                 FROM slots s
                 JOIN bookings b ON s.id = b.slot_id
                 WHERE b.market_id = md.market_id 
                 AND DATE(b.booking_date) = md.date
                 GROUP BY zone
                 ORDER BY COUNT(*) DESC
                 LIMIT 1), 
                ''
            ) as top_zone,
            COALESCE(
                (SELECT COUNT(CASE WHEN b.status = 'completed' THEN 1 END) * 100.0 / 
                        NULLIF(COUNT(*), 0)
                 FROM slots s
                 JOIN bookings b ON s.id = b.slot_id
                 WHERE b.market_id = md.market_id 
                 AND DATE(b.booking_date) = md.date
                 GROUP BY zone
                 ORDER BY COUNT(*) DESC
                 LIMIT 1), 
                0
            ) as top_zone_occupancy
        FROM market_dates md
        LEFT JOIN booking_stats bs ON 
            md.market_id = bs.market_id AND 
            md.date = bs.booking_date
        ORDER BY md.market_id, md.date DESC
    `

	var stats []entities.MarketDashboardStats
	err = repo.db.Raw(query,
		startDate, endDate, // for date range
		startDate, endDate, // for booking_stats
	).Scan(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get all markets stats: %v", err)
	}

	// Calculate growth rates
	for i := range stats {
		if i > 0 && stats[i].MarketID == stats[i-1].MarketID {
			// Calculate booking growth
			if stats[i-1].TotalConfirmBookings > 0 {
				stats[i].BookingGrowth = ((float64(stats[i].TotalConfirmBookings) - float64(stats[i-1].TotalConfirmBookings)) /
					float64(stats[i-1].TotalConfirmBookings)) * 100
			}

			// Calculate revenue growth
			if stats[i-1].TotalRevenue > 0 {
				stats[i].RevenueGrowth = ((stats[i].TotalRevenue - stats[i-1].TotalRevenue) /
					stats[i-1].TotalRevenue) * 100
			}
		}
	}

	return stats, nil
}
func (repo *DashboardRepository) GetDebugInfo(marketID string) map[string]interface{} {
	currentDate := time.Now().Truncate(24 * time.Hour)
	debug := make(map[string]interface{})

	// Count total bookings
	var totalBookings int64
	repo.db.Model(&entities.Booking{}).
		Where("market_id = ? AND DATE(booking_date) = ?", marketID, currentDate).
		Count(&totalBookings)
	debug["total_bookings"] = totalBookings

	// Count completed bookings
	var completedBookings int64
	repo.db.Model(&entities.Booking{}).
		Where("market_id = ? AND DATE(booking_date) = ? AND status = ?",
			marketID, currentDate, "completed").
		Count(&completedBookings)
	debug["completed_bookings"] = completedBookings

	// Get total revenue from completed bookings
	var totalRevenue float64
	repo.db.Model(&entities.Booking{}).
		Joins("LEFT JOIN payments ON bookings.id = payments.booking_id").
		Where("bookings.market_id = ? AND DATE(bookings.booking_date) = ? AND bookings.status = ? AND payments.status = ?",
			marketID, currentDate, "completed", "completed").
		Select("COALESCE(SUM(payments.price), 0)").
		Scan(&totalRevenue)
	debug["total_revenue"] = totalRevenue

	// Get previous day stats for growth calculation
	previousDate := currentDate.AddDate(0, 0, -1)
	var prevStats entities.MarketDashboardStats
	repo.db.Where("market_id = ? AND date = ?", marketID, previousDate).
		First(&prevStats)
	debug["previous_day_bookings"] = prevStats.TotalBookings
	debug["previous_day_revenue"] = prevStats.TotalRevenue

	return debug
}
func (repo *DashboardRepository) GetWeeklyStats(marketID string) ([]entities.MarketDashboardStats, error) {
	// Load Bangkok time zone
	bangkokLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, fmt.Errorf("failed to load Bangkok location: %v", err)
	}

	// Get today's date at midnight in Bangkok time
	now := time.Now().In(bangkokLocation)
	today := time.Date(now.Year(), now.Month(), now.Day()-2, 0, 0, 0, 0, bangkokLocation)
	weekday := today.Weekday()

	// Calculate the start of the current calendar week (Monday) in Bangkok time
	offsetToMonday := int(time.Monday - weekday)
	if offsetToMonday > 0 {
		offsetToMonday = -6 // Adjust to previous Monday if today is Sunday
	}
	startDate := today.AddDate(0, 0, offsetToMonday)
	endDate := today // Set endDate as today to include data up to the current day of the week

	// Log for debugging
	fmt.Printf("Today is: %s\nWeekday is: %s\n", today, weekday)
	fmt.Printf("Querying for market %s between %s (Monday) and %s (Today)\n", marketID, startDate, endDate)

	// Query for stats within the date range for the specified market
	var stats []entities.MarketDashboardStats
	err = repo.db.Where("market_id = ? AND date BETWEEN ? AND ?", marketID, startDate, endDate).
		Order("date ASC"). // Order by date in ascending order
		Find(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get weekly stats for market %s: %v", marketID, err)
	}

	return stats, nil
}
