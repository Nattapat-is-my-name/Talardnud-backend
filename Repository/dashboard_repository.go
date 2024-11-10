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
	// First, let's delete any existing records for this market_id to prevent duplicates
	cleanupQuery := `
        DELETE FROM market_dashboard_stats 
        WHERE market_id = $1 
        AND date IN (
            SELECT DISTINCT DATE(booking_date) 
            FROM bookings 
            WHERE market_id = $1
        )
    `

	err := repo.db.Exec(cleanupQuery, marketID).Error
	if err != nil {
		return fmt.Errorf("failed to cleanup existing stats: %v", err)
	}

	// Now proceed with the main stats update
	query := `
        WITH distinct_dates AS (
            SELECT DISTINCT 
                market_id,
                DATE(booking_date) as date
            FROM bookings
            WHERE market_id = $1
        ),
        booking_stats AS (
            SELECT 
                b.market_id,
                DATE(b.booking_date) as booking_date,
                COUNT(DISTINCT b.id) as total_bookings,
                COUNT(DISTINCT CASE WHEN b.status = 'completed' THEN b.id END) as completed_bookings,
                COUNT(DISTINCT CASE WHEN b.status = 'cancelled' THEN b.id END) as cancelled_bookings,
                COUNT(DISTINCT CASE WHEN b.status = 'pending' THEN b.id END) as pending_bookings,
                COALESCE(SUM(CASE WHEN b.status = 'completed' AND p.status = 'completed' 
                    THEN b.price ELSE 0 END), 0) as total_revenue,
                COALESCE(
                    (COUNT(DISTINCT CASE WHEN b.status = 'completed' THEN b.id END) * 100.0 / 
                    NULLIF(COUNT(DISTINCT b.id), 0)),
                    0
                ) as occupancy_rate
            FROM bookings b
            LEFT JOIN payments p ON b.id = p.booking_id
            WHERE b.market_id = $1
            GROUP BY b.market_id, DATE(b.booking_date)
        ),
        zone_stats AS (
            SELECT 
                b.market_id,
                DATE(b.booking_date) as booking_date,
                s.zone as zone,
                COUNT(*) as zone_bookings,
                COUNT(CASE WHEN b.status = 'completed' THEN 1 END) as completed_in_zone,
                COUNT(CASE WHEN b.status = 'cancelled' THEN 1 END) as cancelled_in_zone,
                COUNT(CASE WHEN b.status = 'pending' THEN 1 END) as pending_in_zone,
                COALESCE(
                    (COUNT(CASE WHEN b.status = 'completed' THEN 1 END) * 100.0 / 
                    NULLIF(COUNT(*), 0)),
                    0
                ) as zone_occupancy
            FROM bookings b
            JOIN slots s ON b.slot_id = s.id
            WHERE b.market_id = $1
            GROUP BY 
                b.market_id, 
                DATE(b.booking_date), 
                s.zone
        ),
        ranked_zones AS (
            SELECT 
                market_id,
                booking_date,
                zone,
                zone_occupancy,
                completed_in_zone,
                MAX(completed_in_zone) OVER (PARTITION BY market_id, booking_date) as max_completed
            FROM zone_stats
        ),
        aggregated_zones AS (
            SELECT 
                market_id,
                booking_date,
                STRING_AGG(zone, ',' ORDER BY zone) as top_zone,
                MAX(zone_occupancy) as top_zone_occupancy
            FROM ranked_zones
            WHERE completed_in_zone = max_completed
            GROUP BY market_id, booking_date
        ),
        prev_day_stats AS (
            SELECT 
                d.market_id,
                d.date,
                LAG(bs.completed_bookings, 1) OVER (
                    PARTITION BY d.market_id ORDER BY d.date
                ) as prev_bookings,
                LAG(bs.total_revenue, 1) OVER (
                    PARTITION BY d.market_id ORDER BY d.date
                ) as prev_revenue
            FROM distinct_dates d
            LEFT JOIN booking_stats bs ON 
                d.market_id = bs.market_id AND 
                d.date = bs.booking_date
        )
        INSERT INTO market_dashboard_stats (
            market_id,
            date,
            total_bookings,
            total_confirm_bookings,
            total_cancel_bookings,
            total_pending_bookings,
            booking_growth,
            total_revenue,
            revenue_growth,
            occupancy_rate,
            top_zone,
            top_zone_occupancy,
            created_at
        )
        SELECT DISTINCT ON (dd.market_id, dd.date)
            dd.market_id,
            dd.date,
            COALESCE(bs.total_bookings, 0),
            COALESCE(bs.completed_bookings, 0),
            COALESCE(bs.cancelled_bookings, 0),
            COALESCE(bs.pending_bookings, 0),
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
            COALESCE(az.top_zone, ''),
            COALESCE(az.top_zone_occupancy, 0),
            NOW()
        FROM distinct_dates dd
        LEFT JOIN booking_stats bs ON 
            dd.market_id = bs.market_id AND 
            dd.date = bs.booking_date
        LEFT JOIN aggregated_zones az ON 
            dd.market_id = az.market_id AND 
            dd.date = az.booking_date
        LEFT JOIN prev_day_stats pds ON 
            dd.market_id = pds.market_id AND 
            dd.date = pds.date
    `

	err = repo.db.Exec(query, marketID).Error
	if err != nil {
		return fmt.Errorf("failed to update stats: %v", err)
	}

	log.Printf("Successfully updated stats for market %s", marketID)
	return nil
}

const mondayOffsetWhenSunday = -6

func (repo *DashboardRepository) GetWeeklyStats(marketID string) ([]entities.MarketDashboardStats, error) {
	// Load Bangkok time zone
	bangkokLocation, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		return nil, fmt.Errorf("failed to load Bangkok time zone: %v", err)
	}

	// Get today's date at midnight in Bangkok time
	now := time.Now().In(bangkokLocation)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, bangkokLocation)
	weekday := today.Weekday()

	// Calculate the start of the current calendar week (Monday) in Bangkok time
	offsetToMonday := int(time.Monday - weekday)
	if offsetToMonday > 0 {
		offsetToMonday = mondayOffsetWhenSunday
	}
	startDate := today.AddDate(0, 0, int(time.Monday-weekday))
	endDate := startDate.AddDate(0, 0, 6) // Set endDate to the end of the current week (Sunday)

	// Log for debugging
	fmt.Printf("Today is: %s\nWeekday is: %s\n", today.Format("2006-01-02"), weekday)
	fmt.Printf("Querying for market %s between %s (Monday) and %s (Sunday)\n", marketID, startDate.Format("2006-01-02"), endDate.Format("2006-01-02"))

	// Query for stats within the date range for the specified market
	var stats []entities.MarketDashboardStats
	err = repo.db.Where("market_id = ? AND date BETWEEN ? AND ?", marketID, startDate, endDate).
		Order("date ASC").
		Find(&stats).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get weekly stats: %v", err)
	}

	// If no data for the current week, create a placeholder
	if len(stats) == 0 {
		stats = []entities.MarketDashboardStats{
			{
				MarketID:             marketID,
				Date:                 startDate,
				TotalBookings:        0,
				TotalConfirmBookings: 0,
				TotalCancelBookings:  0,
				TotalPendingBookings: 0,
				BookingGrowth:        0,
				TotalRevenue:         0,
				RevenueGrowth:        0,
				OccupancyRate:        0,
				TopZone:              "",
				TopZoneOccupancy:     0,
				CreatedAt:            time.Now(),
			},
		}
	}

	return stats, nil
}
