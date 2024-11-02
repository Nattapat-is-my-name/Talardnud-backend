package Services

import (
	"fmt"
	"github.com/go-co-op/gocron"
	"log"
	"time"
	entities "tln-backend/Entities"
	"tln-backend/contact"
)

type BookingService struct {
	scheduler   *gocron.Scheduler
	repo        contact.IBooking
	payment     contact.IPayment
	slotUseCase contact.ISlotUseCase
}

func NewBookingService(repo contact.IBooking, payment contact.IPayment, slotUseCase contact.ISlotUseCase) *BookingService {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()
	return &BookingService{
		scheduler:   scheduler,
		repo:        repo,
		payment:     payment,
		slotUseCase: slotUseCase,
	}
}

func (s *BookingService) ScheduleBookingCancellation(transactionID, bookingID, slotID string, expiresAt time.Time) {
	job, err := s.scheduler.Every(30).Second().StartAt(time.Now()).Do(s.checkBookingStatus, transactionID, bookingID, slotID, expiresAt)
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

func (s *BookingService) checkBookingStatus(transactionID, bookingID, slotID string, expiresAt time.Time) {
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
			err := s.completeBooking(transactionID, transaction.PaymentID, bookingID, slotID)
			if err != nil {
				log.Printf("Error completing booking: %v", err)
				return
			}
		}
		log.Printf("Booking %s has already been processed", bookingID)
		s.RemoveScheduled(bookingID)

	case entities.TransactionFailed:
		if time.Now().Before(expiresAt) {
			err := s.cancelExpiredBooking(transactionID, transaction.PaymentID, bookingID)
			if err != nil {
				log.Printf("Error cancelling failed booking: %v", err)
				return
			}

			s.RemoveScheduled(bookingID)
		}

	case entities.TransactionRefunded:
		if time.Now().Before(expiresAt) {
			err := s.RefundBooking(transactionID, transaction.PaymentID, bookingID, slotID)
			if err != nil {
				log.Printf("Error cancelling refunded booking: %v", err)
				return
			}

			s.RemoveScheduled(bookingID)
		}

	}

}

func (s *BookingService) RefundBooking(transactionID, bookingID, paymentID, slotID string) error {

	if _, err := s.payment.UpdateTransaction(transactionID, entities.TransactionRefunded); err != nil {
		return fmt.Errorf("error updating transaction status: %v", err)
	}

	if _, err := s.repo.UpdateBookingStatus(bookingID, entities.StatusRefunded); err != nil {
		return fmt.Errorf("error updating booking status: %v", err)
	}

	if _, err := s.payment.UpdatePayment(paymentID, entities.PaymentRefunded); err != nil {
		return fmt.Errorf("error updating payment status: %v", err)
	}

	if _, err := s.slotUseCase.UpdateSlotStatus(slotID, entities.StatusAvailable); err != nil {
		return fmt.Errorf("error updating slot status: %v", err)
	}

	s.RemoveScheduled(bookingID)
	return nil
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

func (s *BookingService) completeBooking(transactionID, paymentID, bookingID, slotID string) error {
	if _, err := s.payment.UpdateTransaction(transactionID, entities.TransactionCompleted); err != nil {
		return fmt.Errorf("error updating transaction status: %v", err)
	}

	if _, err := s.repo.UpdateBookingStatus(bookingID, entities.StatusCompleted); err != nil {
		return fmt.Errorf("error updating booking status: %v", err)
	}

	if _, err := s.payment.UpdatePayment(paymentID, entities.PaymentCompleted); err != nil {
		return fmt.Errorf("error updating payment status: %v", err)
	}

	if _, err := s.slotUseCase.UpdateSlotStatus(slotID, entities.StatusBooked); err != nil {
		return fmt.Errorf("error updating slot status: %v", err)
	}
	s.RemoveScheduled(bookingID)

	return nil
}
