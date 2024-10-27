package Database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	entities "tln-backend/Entities"
	"tln-backend/Entities/dtos"
)

func NewDB() (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Bangkok",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_SSLMODE"),
	)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(
		&entities.Vendor{},
		&entities.MarketDashboardStats{},
		&entities.LoginRequest{},


		&dtos.RegisterRequest{},
		&entities.MarketProvider{},
		&entities.Booking{},
		&entities.Market{},
		&entities.Slot{},
		&entities.Payment{},
		&entities.MarketProvider{},
		&entities.Transaction{},
	); err != nil {
		return nil, err
	}

	return db, nil
}
