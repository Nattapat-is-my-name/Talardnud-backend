package Interfaces

import (
	entities "tln-backend/Entities"
)

type IPayment interface {
	CreateTransaction(transaction *entities.Transaction) error
	GetTransaction(ref1, ref2, ref3 string) (*entities.Transaction, error)
	GetTransactionByID(transactionID string) (*entities.Transaction, error)
	UpdatePayment(BookingID string, Status entities.PaymentStatus) (*entities.Payment, error)
	CreatePayment(payment *entities.Payment) error
	UpdateTransaction(TransactionID string, Status entities.TransactionStatus) (*entities.Transaction, error)
}
