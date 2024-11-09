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
                DATE(created_at) as date
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
                DATE(b.created_at) as booking_date,
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
            GROUP BY b.market_id, DATE(b.created_at)
        ),
        zone_stats AS (
            SELECT 
                b.market_id,
                DATE(b.created_at) as booking_date,
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
            GROUP BY b.market_id, DATE(b.created_at), s.zone
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
	err = repo.db.Where("market_id = ? AND created_at BETWEEN ? AND ?", marketID, startDate, endDate).
		Order("date ASC"). // Order by date in ascending order
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
