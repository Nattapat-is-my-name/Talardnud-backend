package Usecase

import (
	"fmt"
	"github.com/google/uuid"
	"log"
	"strconv"
	"time"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Services"
	"tln-backend/contact"
)

type BookingUseCase struct {
	repo           contact.IBooking
	payment        contact.IPayment
	PaymentUseCase *PaymentUseCase
	bookingService *Services.BookingService
	slotUseCase    contact.ISlotUseCase
}

func NewBookingUseCase(repo contact.IBooking, payment contact.IPayment, paymentUseCase *PaymentUseCase, bookingService *Services.BookingService, slotUseCase contact.ISlotUseCase) *BookingUseCase {
	return &BookingUseCase{
		repo:           repo,
		payment:        payment,
		PaymentUseCase: paymentUseCase,
		bookingService: bookingService,
		slotUseCase:    slotUseCase,
	}
}

func (uc *BookingUseCase) CreateBooking(bookingReq *entitiesDtos.BookingRequest) (*entitiesDtos.BookingResponse, *entitiesDtos.ErrorResponse) {
	if err := validateBooking(bookingReq); err != nil {
		log.Printf("Validation failed: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid booking request: " + err.Error(),
		}
	}

	bookingDate, err := time.Parse("2006-01-02", bookingReq.BookingDate)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: "Invalid booking date format: " + err.Error(),
		}
	}

	if err := uc.repo.IsSlotAvailable(bookingReq); err != nil {
		log.Printf("Slot is not available: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    409,
			Message: "Slot is not available: " + err.Error(),
		}
	}

	expirationTime := time.Now().Add(30 * time.Minute)
	bookingEntity := &entities.Booking{
		ID:          uuid.New().String(),
		SlotID:      bookingReq.SlotID,
		VendorID:    bookingReq.VendorID,
		MarketID:    bookingReq.MarketID,
		BookingDate: bookingDate,
		Status:      entities.StatusPending,
		Method:      bookingReq.Method,
		Price:       bookingReq.Price,
		ExpiresAt:   expirationTime,
	}

	if err := uc.repo.CreateBooking(bookingEntity); err != nil {
		log.Printf("Error creating booking: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to create booking: " + err.Error(),
		}
	}

	paymentEntity := entities.Payment{
		ID:          uuid.New().String(),
		BookingID:   bookingEntity.ID,
		Price:       bookingReq.Price,
		Method:      bookingReq.Method,
		Status:      entities.PaymentPending,
		PaymentDate: time.Now(),
		ExpiresAt:   expirationTime,
	}

	if err := uc.payment.CreatePayment(&paymentEntity); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to create payment: " + err.Error(),
		}
	}

	var promptPayResult entitiesDtos.PromptPayResult
	if bookingReq.Method == "PromptPay" {
		promptPayResult, err = uc.handlePayment(paymentEntity, paymentEntity.ID)
		if err != nil {
			return nil, &entitiesDtos.ErrorResponse{
				Code:    500,
				Message: "Failed to handle payment: " + err.Error(),
			}
		}
	}

	Price, err := strconv.ParseFloat(promptPayResult.PromptPayDetail.Amount, 64)
	if err != nil {
		log.Printf("Failed to parse amount: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Failed to parse amount: %v", err),
		}
	}

	transaction := &entities.Transaction{
		ID:              uuid.New().String(),
		PaymentID:       paymentEntity.ID,
		Method:          string(paymentEntity.Method),
		TransactionID:   promptPayResult.PromptPayDetail.TransactionID,
		Ref1:            promptPayResult.PromptPayDetail.Ref1,
		Ref2:            promptPayResult.PromptPayDetail.Ref2,
		Ref3:            promptPayResult.PromptPayDetail.Ref3,
		Price:           Price,
		Image:           promptPayResult.QRResponse.Data.QRImage,
		Status:          entities.TransactionPending,
		TransactionDate: time.Now(),
		ExpiresAt:       expirationTime,
	}

	err = uc.PaymentUseCase.repo.CreateTransaction(transaction)
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create transaction: %v", err),
		}
	}

	bookingResponse := entitiesDtos.BookingResponse{
		ID:            bookingEntity.ID,
		SlotID:        bookingEntity.SlotID,
		VendorID:      bookingEntity.VendorID,
		TransactionID: transaction.ID,
		BookingDate:   bookingEntity.BookingDate,
		Price:         bookingEntity.Price,
		Status:        bookingEntity.Status,
		Method:        bookingEntity.Method,
		Image:         transaction.Image,
		ExpiresAt:     transaction.ExpiresAt,
	}

	// Schedule booking cancellation using BookingService
	uc.bookingService.ScheduleBookingCancellation(transaction.ID, bookingEntity.ID, bookingEntity.SlotID, expirationTime)

	return &bookingResponse, nil
}

// Other methods remain the same...

func (uc *BookingUseCase) GetBooking(bookingID string) (*entities.Booking, error) {
	return uc.repo.GetBooking(bookingID)
}

func (uc *BookingUseCase) GetBookingsByUser(userID string) ([]entities.Booking, error) {
	return uc.repo.GetBookingsByUser(userID)
}

func (uc *BookingUseCase) handlePayment(paymentEntity entities.Payment, paymentID string) (entitiesDtos.PromptPayResult, error) {
	promptPayResult, errResp := uc.PaymentUseCase.PromptPay(paymentEntity, paymentID)
	if errResp != nil {
		return entitiesDtos.PromptPayResult{}, fmt.Errorf("failed to generate PromptPay QR code: %v", errResp)
	}
	return *promptPayResult, nil
}

// validateCancelBooking validates the cancel booking request.
func validateCancelBooking(cancelBookingReq *entitiesDtos.CancelBookingRequest) error {
	if cancelBookingReq.BookingID == "" {
		return fmt.Errorf("booking ID is required")
	}
	return nil
}
func validateBooking(booking *entitiesDtos.BookingRequest) error {
	if booking.VendorID == "" {
		return fmt.Errorf("vendor ID is required")
	}
	if booking.MarketID == "" {
		return fmt.Errorf("market ID is required")
	}
	if booking.SlotID == "" {
		return fmt.Errorf("slot ID is required")
	}
	if booking.BookingDate == "" {
		return fmt.Errorf("booking date is required")
	}
	if booking.Method == "" {
		return fmt.Errorf("payment method is required")
	}
	return nil
}

// CancelBooking cancels an existing booking based on the provided request.//

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
			Code:    404,
			Message: "Failed to get booking: " + err.Error(),
		}
	}

	// Validate that booking exists and has required fields
	if bookingEntity == nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "Booking not found",
		}
	}

	if bookingEntity.Payment == nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Booking payment information not found",
		}
	}

	// Get payment by booking id
	paymentEntity, err := uc.payment.GetPayment(bookingEntity.Payment.ID)
	if err != nil {
		log.Printf("Error getting payment with ID %s: %v", bookingEntity.Payment.ID, err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to get payment: " + err.Error(),
		}
	}

	if paymentEntity == nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "Payment not found",
		}
	}

	// Get transaction by payment id
	transactionEntity, err := uc.payment.GetTransactionByID(paymentEntity.TransactionID)
	if err != nil {
		log.Printf("Error getting transaction with ID %s: %v", paymentEntity.TransactionID, err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Failed to get transaction: " + err.Error(),
		}
	}

	if transactionEntity == nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    404,
			Message: "Transaction not found",
		}
	}

	// Check if slot use case is properly initialized
	if uc.slotUseCase == nil {
		log.Printf("Error: slotUseCase is nil")
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Internal server error: slot service not initialized",
		}
	}

	switch transactionEntity.Status {
	case "completed":
		// Update transaction status
		_, err = uc.payment.UpdateTransaction(transactionEntity.ID, entities.TransactionRefunded)
		if err != nil {
			log.Printf("Error updating transaction status for ID %s: %v", transactionEntity.ID, err)
			return nil, &entitiesDtos.ErrorResponse{
				Code:    500,
				Message: "Failed to update transaction status: " + err.Error(),
			}
		}

		// Update booking status
		_, err = uc.repo.UpdateBookingStatus(cancelBookingReq.BookingID, entities.StatusRefunded)
		if err != nil {
			log.Printf("Error updating booking status for ID %s: %v", cancelBookingReq.BookingID, err)
			return nil, &entitiesDtos.ErrorResponse{
				Code:    500,
				Message: "Failed to update booking status: " + err.Error(),
			}
		}

		// Update payment status
		_, err = uc.payment.UpdatePayment(paymentEntity.ID, entities.PaymentRefunded)
		if err != nil {
			log.Printf("Error updating payment status for ID %s: %v", paymentEntity.ID, err)
			return nil, &entitiesDtos.ErrorResponse{
				Code:    500,
				Message: "Failed to update payment status: " + err.Error(),
			}
		}

		// Update slot status to available
		_, err = uc.slotUseCase.UpdateSlotStatus(bookingEntity.SlotID, entities.StatusAvailable)
		if err != nil {
			// Log the error but don't fail the cancellation process
			log.Printf("Warning: Error updating slot status for slot ID %s: %v", bookingEntity.SlotID, err)
		}

		return &entitiesDtos.BookingResponse{
			ID:          bookingEntity.ID,
			SlotID:      bookingEntity.SlotID,
			VendorID:    bookingEntity.VendorID,
			BookingDate: bookingEntity.BookingDate,
			Price:       bookingEntity.Price,
			Status:      entities.StatusRefunded,
		}, nil

	case "failed":
		return nil, &entitiesDtos.ErrorResponse{
			Code:    409,
			Message: fmt.Sprintf("Transaction with ID %s is already %s", transactionEntity.ID, transactionEntity.Status),
		}

	case "pending":
		// Update transaction status to failed
		_, err = uc.payment.UpdateTransaction(transactionEntity.ID, "failed")
		if err != nil {
			log.Printf("Error updating transaction status for ID %s: %v", transactionEntity.ID, err)
			return nil, &entitiesDtos.ErrorResponse{
				Code:    500,
				Message: "Failed to update transaction status: " + err.Error(),
			}
		}

		return &entitiesDtos.BookingResponse{
			ID:          bookingEntity.ID,
			SlotID:      bookingEntity.SlotID,
			VendorID:    bookingEntity.VendorID,
			BookingDate: bookingEntity.BookingDate,
			Price:       bookingEntity.Price,
			Status:      "CANCELLED",
		}, nil

	default:
		return nil, &entitiesDtos.ErrorResponse{
			Code:    409,
			Message: fmt.Sprintf("Unexpected transaction status: %s", transactionEntity.Status),
		}
	}
}
