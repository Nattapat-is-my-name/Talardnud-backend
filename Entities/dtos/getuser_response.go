package dtos

type GetUserResponse struct {
	ID       string      `json:"id"`
	Username string      `json:"username"`
	Email    string      `json:"email"`
	Bookings BookingDtos `json:"bookings"`
}
