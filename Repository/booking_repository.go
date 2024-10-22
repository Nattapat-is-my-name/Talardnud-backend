package Repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type BookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (repo *BookingRepository) CreateBooking(booking *entities.Booking) error {
	return repo.db.Create(booking).Error
}

func (repo *BookingRepository) GetBooking(bookingID string) (*entities.Booking, error) {
	var booking entities.Booking

	result := repo.db.Where("id = ?", bookingID).First(&booking)

	if result.Error != nil {
		return nil, result.Error
	}

	result = repo.db.Preload("Payment").Where("id = ?", bookingID).First(&booking)
	if result.Error != nil {
		return nil, result.Error
	}

	return &booking, nil
}

func (repo *BookingRepository) GetBookingsByUser(userID string) ([]entities.Booking, error) {
	var bookings []entities.Booking

	result := repo.db.Where("vendor_id = ?", userID).Find(&bookings)

	if result.Error != nil {
		return nil, result.Error
	}

	//create preload for slot entitle and market
	result = repo.db.Preload("Slot").Where("slot_id = ?", userID).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error

	}

	//create preload for payment
	result = repo.db.Preload("Payment").Where("vendor_id = ?", userID).Find(&bookings)
	if result.Error != nil {
		return nil, result.Error
	}

	return bookings, nil
}

func (repo *BookingRepository) UpdateBookingStatus(bookingID string, status entities.BookingStatus) (*entities.Booking, error) {
	var booking entities.Booking
	result := repo.db.Model(&booking).Where("ID = ?", bookingID).Update("status", status)

	if result.Error != nil {
		return nil, result.Error
	}

	return &booking, nil
}

func (repo *BookingRepository) IsSlotAvailable(bookingReq *entitiesDtos.BookingRequest) error {
	// First, check if the slot exists
	var slot entities.Slot
	if err := repo.db.Where("ID = ? AND market_id = ?", bookingReq.SlotID, bookingReq.MarketID).First(&slot).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("slot not found")
		}
		return fmt.Errorf("error checking slot existence: %w", err)
	}

	// Check for existing bookings on the requested date
	var count int64
	err := repo.db.Model(&entities.Booking{}).
		Where(" slot_id = ? AND DATE(booking_date) = ? AND status IN ('pending', 'completed')",
			bookingReq.SlotID, bookingReq.BookingDate).
		Count(&count).Error

	if err != nil {
		return fmt.Errorf("error checking existing bookings: %w", err)
	}

	if count > 0 {
		return fmt.Errorf("you already have a pending or confirmed booking for this slot")
	}

	return nil
}
