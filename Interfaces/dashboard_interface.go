package Interfaces

import entities "tln-backend/Entities"

type IDashboard interface {
	GetDashboardData(marketID string) (*entities.DashboardResponse, error)
	GetAllMarketsDashboardStats() ([]entities.MarketDashboardStats, error)
	GetWeeklyStats(marketID string) ([]entities.MarketDashboardStats, error)
}
