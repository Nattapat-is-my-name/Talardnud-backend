// entities/dashboard.go
package entities

import "time"

// MarketDashboardStats represents the top stats panel
type MarketDashboardStats struct {
	MarketID             string    `gorm:"primaryKey;type:varchar(36)" json:"market_id"`
	Date                 time.Time `gorm:"primaryKey;type:date" json:"date"`
	TotalBookings        int       `gorm:"not null;default:0" json:"total_bookings"`
	TotalConfirmBookings int       `gorm:"not null;default:0" json:"total_confirm_bookings"`
	TotalCancelBookings  int       `gorm:"not null;default:0" json:"total_cancel_bookings"`
	TotalPendingBookings int       `gorm:"not null;default:0" json:"total_pending_bookings"`
	BookingGrowth        float64   `gorm:"type:decimal(5,2);default:0" json:"booking_growth"`
	TotalRevenue         float64   `gorm:"type:decimal(10,2);not null;default:0" json:"total_revenue"`
	RevenueGrowth        float64   `gorm:"type:decimal(5,2);default:0" json:"revenue_growth"`
	OccupancyRate        float64   `gorm:"type:decimal(5,2);default:0" json:"occupancy_rate"`
	TopZone              string    `gorm:"type:varchar(50)" json:"top_zone"`
	TopZoneOccupancy     float64   `gorm:"type:decimal(5,2);default:0" json:"top_zone_occupancy"`
	CreatedAt            time.Time `gorm:"autoCreateTime" json:"created_at"`
}
type SlotNameSummary struct {
	Name  string `json:"name"`
	Count int    `json:"count"`
}

type DashboardResponse struct {
	Stats []MarketDashboardStats `json:"stats"` // Changed to slice
}
