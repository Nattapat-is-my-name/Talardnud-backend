package Repository

import (
	"fmt"
	"gorm.io/gorm"
	"log"
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

func (repo *BookingRepository) IsBookingExists(bookingReq *entitiesDtos.BookingRequest) (bool, error) {
	var count int64

	// Check for overlapping bookings with "pending" status
	err := repo.db.Model(&entities.Booking{}).
		Where("slot_id = ? AND status = ? AND ((start_date < ? AND end_date > ?) OR (start_date < ? AND end_date > ?))",
			bookingReq.SlotID, "pending", bookingReq.EndDate, bookingReq.StartDate, bookingReq.StartDate, bookingReq.EndDate).
		Count(&count).Error

	if err != nil {
		// Log and return error for debugging
		log.Printf("Failed to check booking existence: %v", err)
		return false, fmt.Errorf("failed to check booking existence: %w", err)
	}

	return count > 0, nil
}

func (repo *BookingRepository) GetBooking(bookingID string) (*entities.Booking, error) {
	var booking entities.Booking

	result := repo.db.Where("id = ?", bookingID).First(&booking)

	if result.Error != nil {
		return nil, result.Error
	}

	return &booking, nil
}

func (repo *BookingRepository) UpdateBookingStatus(bookingID string, status string) error {
	return repo.db.Model(&entities.Booking{}).Where("id = ?", bookingID).Update("status", status).Error
}
