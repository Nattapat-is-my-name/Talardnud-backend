package Repository

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type UserRepository struct {
	db *gorm.DB
}

func (r *UserRepository) GetUserByID(id string) (*entitiesDtos.GetUserResponse, error) {

	id = strings.TrimSpace(id)

	var user entities.Vendor
	result := r.db.Where("id = ?", id).First(&user)
	if result.Error != nil {
		fmt.Println("no user found with ID: ", id)
		return nil, result.Error
	}

	var bookings []entities.Booking
	bookingResult := r.db.Where("vendor_id = ?", id).Find(&bookings)
	if bookingResult.Error != nil {
		fmt.Println("no bookings found for user with ID: ", id)
		return nil, bookingResult.Error
	}

	var bookingIDs []string
	for _, booking := range bookings {
		bookingIDs = append(bookingIDs, booking.ID)
	}

	getUserResponse := entitiesDtos.GetUserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Bookings: entitiesDtos.BookingDtos{
			IDs: bookingIDs,
		},
	}

	return &getUserResponse, nil
}

func (r *UserRepository) UpdateUser(user *entities.Vendor) error {
	//TODO implement me
	panic("implement me")
}

func (r *UserRepository) DeleteUser(id string) error {

	id = strings.TrimSpace(id)

	fmt.Println("Deleting user with ID: ", id)

	result := r.db.Where("id = ?", id).Delete(&entities.Vendor{})

	if result.Error != nil {
		fmt.Println(
			"Error deleting user with ID: ", id)
		return result.Error
	}

	// Optional: Check if any rows were affected
	if result.RowsAffected == 0 {
		return fmt.Errorf("user with ID %s not found", id)
	}

	return nil
}
func (r *UserRepository) ListUsers() ([]*entities.Vendor, error) {
	//TODO implement me
	panic("implement me")
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *entities.Vendor) error {

	var existingUser entities.Vendor
	if err := r.db.Where("username = ?", user.Username).First(&existingUser).Error; err == nil {
		return errors.New("username already exists")
	}

	return r.db.Create(user).Error
}
