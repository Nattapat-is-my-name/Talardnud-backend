package Interfaces

import (
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type IBooking interface {
	CreateBooking(booking *entities.Booking) error
	//IsBookingExists(bookingReq *entitiesDtos.BookingRequest) (bool, error)
	GetBooking(bookingID string) (*entities.Booking, error)
	UpdateBookingStatus(bookingID string, status entities.BookingStatus) (*entities.Booking, error)
	IsSlotAvailable(bookingReq *entitiesDtos.BookingRequest) error
}
