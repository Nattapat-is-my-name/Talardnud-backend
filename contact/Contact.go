package contact

import (
	"time"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type IBookingService interface {
	ScheduleBookingCancellation(transactionID, bookingID, slotID string, expiresAt time.Time)
	RemoveScheduled(bookingID string)
}

type ISlotUseCase interface {
	UpdateSlotStatus(slotID string, status entities.SlotStatus) (*entities.Slot, *entitiesDtos.ErrorResponse)
}
type IBooking interface {
	CreateBooking(booking *entities.Booking) error
	//IsBookingExists(bookingReq *entitiesDtos.BookingRequest) (bool, error)
	GetBooking(bookingID string) (*entities.Booking, error)
	UpdateBookingStatus(bookingID string, status entities.BookingStatus) (*entities.Booking, error)
	IsSlotAvailable(bookingReq *entitiesDtos.BookingRequest) error
	GetBookingsByUser(userID string) ([]entities.Booking, error)
}

type IPayment interface {
	CreateTransaction(transaction *entities.Transaction) error
	GetPayment(paymentID string) (*entitiesDtos.BookingResponse, error)
	GetTransaction(ref1, ref2, ref3 string) (*entities.Transaction, error)
	GetTransactionByID(transactionID string) (*entities.Transaction, error)
	UpdatePayment(BookingID string, Status entities.PaymentStatus) (*entities.Payment, error)
	CreatePayment(payment *entities.Payment) error
	UpdateTransaction(TransactionID string, Status entities.TransactionStatus) (*entities.Transaction, error)
}
