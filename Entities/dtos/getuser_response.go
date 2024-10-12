package dtos

type GetUserResponse struct {
	ID        string      `json:"id"`
	Username  string      `json:"username"`
	Email     string      `json:"email"`
	FirstName string      `json:"firstName"`
	LastName  string      `json:"lastName"`
	Bookings  BookingDtos `json:"bookings"`
}
