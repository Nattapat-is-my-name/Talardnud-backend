package dtos

import (
	"time"
)

type GetUserResponse struct {
	ID        string      `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	Bookings  BookingDtos `json:"bookings"`
}
