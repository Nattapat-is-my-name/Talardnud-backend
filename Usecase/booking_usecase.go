package Usecase

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"time"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
)

type BookingUseCase struct {
	repo           Interfaces.IBooking
	payment        Interfaces.IPayment
	PaymentUseCase *PaymentUseCase
}

func NewBookingUseCase(repo Interfaces.IBooking, payment Interfaces.IPayment, paymentUseCase *PaymentUseCase) *BookingUseCase {
	return &BookingUseCase{
		repo:           repo,
		payment:        payment,
		PaymentUseCase: paymentUseCase,
	}
}

func (uc *BookingUseCase) CreateBooking(bookingReq *entitiesDtos.BookingRequest) (*entitiesDtos.BookingResponse, *entitiesDtos.ErrorResponse) {
	// Validate the booking request
	if err := validateBooking(bookingReq); err != nil {
		log.Printf("Validation failed: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid booking request: " + err.Error(),
		}
	}

	//Check slot available or not

	// Check if a conflicting booking exists with 'pending' status
	exists, err := uc.repo.IsBookingExists(bookingReq)
	if err != nil {
		log.Printf("Error checking booking existence: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to check booking existence: " + err.Error(),
		}
	}

	if exists {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    409,
			Message: "A pending booking already exists for the selected slot and time. Please cancel the existing booking to create a new one.",
		}
	}

	// Create the booking entity
	bookingEntity := entities.Booking{
		ID: uuid.New().String(),

		VendorID:    bookingReq.VendorID,
		StartDate:   bookingReq.StartDate,
		EndDate:     bookingReq.EndDate,
		Status:      "pending",
		BookingDate: time.Now(),
	}

	// Save the booking entity to the database
	if err := uc.repo.CreateBooking(&bookingEntity); err != nil {
		log.Printf("Error creating booking: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to create booking: " + err.Error(),
		}
	}

	var paymentEntity entities.Payment

	paymentEntity = entities.Payment{
		ID:          uuid.New().String(),
		BookingID:   bookingEntity.ID,
		Amount:      bookingReq.Amount,
		Method:      bookingReq.Method,
		Status:      "pending",
		PaymentDate: time.Now(),
	}

	if err := uc.payment.CreatePayment(&paymentEntity); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to create payment: " + err.Error(),
		}
	}

	var promptPayResult *entitiesDtos.PromptPayResult
	if bookingReq.Method == "PromptPay" {
		promptPayResult, err = uc.handlePayment(&paymentEntity, paymentEntity.ID)
		if err != nil {
			return nil, &entitiesDtos.ErrorResponse{
				Code:    500,
				Message: "Failed to handle payment: " + err.Error(),
			}
		}
	}

	bookingResponse := &entitiesDtos.BookingResponse{
		ID: bookingEntity.ID,

		VendorID:    bookingEntity.VendorID,
		BookingDate: bookingEntity.BookingDate,
		StartDate:   bookingEntity.StartDate,
		EndDate:     bookingEntity.EndDate,
		Status:      bookingEntity.Status,
		Amount:      paymentEntity.Amount,
		Method:      paymentEntity.Method,
	}
	if promptPayResult != nil {
		bookingResponse.PromptPay = promptPayResult.QRResponse
	}

	return bookingResponse, nil
}

func validateBooking(booking *entitiesDtos.BookingRequest) error {
	if booking.SlotID == "" || booking.VendorID == "" || booking.StartDate.IsZero() || booking.EndDate.IsZero() {
		return fmt.Errorf("missing required booking information")
	}
	if booking.StartDate.After(booking.EndDate) {
		return fmt.Errorf("start date must be before end date")
	}
	return nil
}

func (uc *BookingUseCase) GetBooking(bookingID string) (*entities.Booking, error) {
	return uc.repo.GetBooking(bookingID)
}

func (uc *BookingUseCase) handlePayment(paymentEntity *entities.Payment, paymentID string) (*entitiesDtos.PromptPayResult, error) {

	// Call the PaymentUseCase to generate a PromptPay QR code
	promptPayResult, errResp := uc.PaymentUseCase.PromptPay(paymentEntity, paymentID)
	if errResp != nil {
		return nil, fmt.Errorf("failed to generate PromptPay QR code: %v", errResp)
	}

	return promptPayResult, nil

}

// validateCancelBooking validates the cancel booking request.
func validateCancelBooking(cancelBookingReq *entitiesDtos.CancelBookingRequest) error {
	if cancelBookingReq.BookingID == "" {
		return fmt.Errorf("booking ID is required")
	}
	return nil
}

// CancelBooking cancels an existing booking based on the provided request.
func (uc *BookingUseCase) CancelBooking(cancelBookingReq *entitiesDtos.CancelBookingRequest) (*entitiesDtos.BookingResponse, *entitiesDtos.ErrorResponse) {
	// Validate the cancel booking request
	if err := validateCancelBooking(cancelBookingReq); err != nil {
		log.Printf("Validation failed: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid cancel booking request: " + err.Error(),
		}
	}

	// Get the booking entity based on BookingID
	bookingEntity, err := uc.repo.GetBooking(cancelBookingReq.BookingID)
	if err != nil {
		log.Printf("Error getting booking with ID %s: %v", cancelBookingReq.BookingID, err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to get booking: " + err.Error(),
		}
	}

	// Check if the booking is already cancelled or completed
	if bookingEntity.Status == "CANCELLED" || bookingEntity.Status == "COMPLETED" {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    409,
			Message: fmt.Sprintf("Booking with ID %s is already %s", cancelBookingReq.BookingID, bookingEntity.Status),
		}
	}

	// Update the booking status to "cancelled"
	if err := uc.repo.UpdateBookingStatus(cancelBookingReq.BookingID, "cancelled"); err != nil {
		log.Printf("Error updating booking status for ID %s: %v", cancelBookingReq.BookingID, err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to update booking status: " + err.Error(),
		}
	}

	// Update the payment status to "cancelled"
	paymentRes, err := uc.payment.UpdatePayment(cancelBookingReq.BookingID, "CANCELLED")
	if err != nil {
		log.Printf("Error updating payment status for booking ID %s: %v", cancelBookingReq.BookingID, err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to update payment status: " + err.Error(),
		}
	}

	//update transaction status to "cancelled"
	_, err = uc.payment.UpdateTransaction(paymentRes.ID, "cancelled")
	if err != nil {
		log.Printf("Error updating transaction status for payment ID %s: %v", paymentRes.ID, err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to update transaction status: " + err.Error(),
		}
	}

	// Create a response object with updated booking data
	bookingResponse := &entitiesDtos.BookingResponse{
		ID:          bookingEntity.ID,
		SlotID:      bookingEntity.SlotID,
		VendorID:    bookingEntity.VendorID,
		BookingDate: bookingEntity.BookingDate,
		StartDate:   bookingEntity.StartDate,
		EndDate:     bookingEntity.EndDate,
		Status:      "CANCELLED",
	}

	log.Printf("Successfully cancelled booking with ID %s", cancelBookingReq.BookingID)
	return bookingResponse, nil
}
