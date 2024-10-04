package Repository

import (
	"errors"
	"gorm.io/gorm"
	entities "tln-backend/Entities"
)

type PaymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreateTransaction(transaction *entities.Transaction) error {
	if err := r.db.Create(&transaction).Error; err != nil {
		return err
	}
	return nil
}

func (r *PaymentRepository) GetTransaction(ref1, ref2, ref3 string) (*entities.Transaction, error) {

	var transaction entities.Transaction
	if err := r.db.Where("ref1 = ? AND ref2 = ? AND ref3 = ?", ref1, ref2, ref3).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err

	}
	return &transaction, nil
}

func (r *PaymentRepository) GetTransactionByID(transactionID string) (*entities.Transaction, error) {

	var transaction entities.Transaction
	if err := r.db.Where("id = ?", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("transaction not found")
		}
		return nil, err

	}
	return &transaction, nil
}

func (r *PaymentRepository) UpdatePayment(paymentID string, Status entities.PaymentStatus) (*entities.Payment, error) {

	var payment entities.Payment
	result := r.db.Model(&payment).Where("ID = ?", paymentID).Update("status", Status)
	if result.Error != nil {
		return nil, result.Error
	}

	return &payment, nil
}

func (r *PaymentRepository) CreatePayment(payment *entities.Payment) error {

	if err := r.db.Create(&payment).Error; err != nil {
		return err
	}

	return nil
}

func (r *PaymentRepository) UpdateTransaction(TransactionID string, Status entities.TransactionStatus) (*entities.Transaction, error) {

	var transaction entities.Transaction
	result := r.db.Model(&transaction).Where("id = ?", TransactionID).Update("status", Status)
	if result.Error != nil {
		return nil, result.Error

	}

	return &transaction, nil

}
