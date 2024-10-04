package Services

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"
	entities "tln-backend/Entities"
	"tln-backend/Interfaces"
)

type BookingService struct {
	scheduler *gocron.Scheduler
	repo      Interfaces.IBooking
	payment   Interfaces.IPayment
}

func NewBookingService(repo Interfaces.IBooking, payment Interfaces.IPayment) *BookingService {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()
	return &BookingService{
		scheduler: scheduler,
		repo:      repo,
		payment:   payment,
	}
}

func (s *BookingService) ScheduleBookingCancellation(transactionID, bookingID string, expiresAt time.Time) {
	job, err := s.scheduler.Every(30).Second().StartAt(time.Now()).Do(s.checkBookingStatus, transactionID, bookingID, expiresAt)
	if err != nil {
		log.Printf("Error scheduling booking status check: %v", err)
		return
	}
	job.Tag(fmt.Sprintf("check-booking-%s", bookingID))
}

func (s *BookingService) RemoveScheduled(bookingID string) {
	err := s.scheduler.RemoveByTag(fmt.Sprintf("check-booking-%s", bookingID))
	if err != nil {
		return
	}
}

func (s *BookingService) checkBookingStatus(transactionID, bookingID string, expiresAt time.Time) {
	transaction, err := s.payment.GetTransactionByID(transactionID)
	if err != nil {
		log.Printf("Error fetching transaction for cancellation: %v", err)
		return
	}

	switch transaction.Status {
	case entities.TransactionPending:
		if time.Now().After(expiresAt) {
			err := s.cancelExpiredBooking(transactionID, transaction.PaymentID, bookingID)
			if err != nil {
				log.Printf("Error cancelling expired booking: %v", err)
				return
			}
			log.Printf("Successfully cancelled expired booking %s", bookingID)
		} else {
			log.Printf("Booking %s is still pending", bookingID)
		}

	case entities.TransactionCompleted:
		if time.Now().Before(expiresAt) {
			err := s.completeBooking(transactionID, transaction.PaymentID, bookingID)
			if err != nil {
				log.Printf("Error completing booking: %v", err)
				return
			}
		}
		log.Printf("Booking %s has already been processed", bookingID)
		s.RemoveScheduled(bookingID)

	case entities.TransactionFailed:
		log.Printf("Booking %s has already been cancelled", bookingID)
		s.RemoveScheduled(bookingID)
	}
}

func (s *BookingService) cancelExpiredBooking(transactionID, paymentID, bookingID string) error {
	if _, err := s.payment.UpdateTransaction(transactionID, entities.TransactionFailed); err != nil {
		return fmt.Errorf("error updating transaction status: %v", err)
	}

	if _, err := s.payment.UpdatePayment(paymentID, entities.PaymentFailed); err != nil {
		return fmt.Errorf("error updating payment status: %v", err)
	}

	if _, err := s.repo.UpdateBookingStatus(bookingID, entities.StatusCancelled); err != nil {
		return fmt.Errorf("error updating booking status: %v", err)
	}

	s.RemoveScheduled(bookingID)
	return nil
}

func (s *BookingService) completeBooking(transactionID, paymentID, bookingID string) error {
	if _, err := s.payment.UpdateTransaction(transactionID, entities.TransactionCompleted); err != nil {
		return fmt.Errorf("error updating transaction status: %v", err)
	}

	if _, err := s.repo.UpdateBookingStatus(bookingID, entities.StatusCompleted); err != nil {
		return fmt.Errorf("error updating booking status: %v", err)
	}

	if _, err := s.payment.UpdatePayment(paymentID, entities.PaymentCompleted); err != nil {
		return fmt.Errorf("error updating payment status: %v", err)
	}

	return nil
}
