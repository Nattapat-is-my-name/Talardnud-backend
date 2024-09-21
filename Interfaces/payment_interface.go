package Interfaces

import (
	entities "tln-backend/Entities"
)

type IPayment interface {
	CreateTransaction(transaction *entities.Transaction) error
	GetTransaction(ref1, ref2, ref3 string) (*entities.Transaction, error)
	UpdatePayment(BookingID string, Status string) (*entities.Payment, error)
	CreatePayment(payment *entities.Payment) error
	UpdateTransaction(TransactionID string, Status string) (*entities.Transaction, error)
}
