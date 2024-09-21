package Usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
	"strconv"
	entities2 "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
	"tln-backend/Services"
)

type PaymentUseCase struct {
	repo           Interfaces.IPayment
	paymentService *Services.PaymentService
}

func NewPaymentUseCase(repo Interfaces.IPayment, paymentService *Services.PaymentService) *PaymentUseCase {
	return &PaymentUseCase{
		repo:           repo,
		paymentService: paymentService,
	}
}

func (uc *PaymentUseCase) PromptPay(request *entities2.Payment, paymentID string) (*entitiesDtos.PromptPayResult, *entitiesDtos.ErrorResponse) {
	// Get the OAuth token
	oauthResp, err := uc.paymentService.GetOAuthToken()
	if err != nil {
		log.Printf("Failed to get OAuth token: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to get OAuth token: %v", err),
		}
	}

	token := oauthResp.Data.AccessToken
	UUID := oauthResp.Data.UUID

	qrResp, promptPay, err := uc.paymentService.CreateQRCode(token, UUID, request.Amount)
	if err != nil {
		log.Printf("Failed to create QR code: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create QR code: %v", err),
		}
	}

	amount, err := strconv.ParseFloat(promptPay.Amount, 64)
	if err != nil {
		log.Printf("Failed to parse amount: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    400,
			Message: fmt.Sprintf("Failed to parse amount: %v", err),
		}
	}

	transaction := &entities2.Transaction{
		ID: uuid.New().String(),

		PaymentID:     paymentID,
		Method:        "Promptpay",
		TransactionID: promptPay.TransactionID,
		Ref1:          promptPay.Ref1,
		Ref2:          promptPay.Ref2,
		Ref3:          promptPay.Ref3,
		Amount:        amount,
		Status:        "PENDING",
	}
	err = uc.repo.CreateTransaction(transaction)
	if err != nil {
		log.Printf("Failed to create transaction: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create transaction: %v", err),
		}
	}

	log.Printf("PromptPay request processed successfully")
	return &entitiesDtos.PromptPayResult{
		QRResponse:  qrResp,
		Transaction: transaction,
	}, nil
}

func (uc *PaymentUseCase) PaymentConfirmation(request *entities2.PaymentConfirmation) (*entities2.Transaction, *entitiesDtos.ErrorResponse) {
	// Get the OAuth token
	oauthResp, err := uc.paymentService.GetOAuthToken()
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to get OAuth token: %v", err),
		}
	}

	// Convert the request struct to JSON
	jsonReqBody, err := json.Marshal(request)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create request body: %v", err),
		}
	}

	// Create and send the HTTP request
	req, err := http.NewRequest("POST", "https://api-sandbox.partners.scb/partners/sandbox/v1/payment/billpayment/transactions", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create request: %v", err),
		}
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", oauthResp.Data.AccessToken))
	req.Header.Set("requestUId", uuid.New().String())
	req.Header.Set("resourceOwnerID", os.Getenv("API_KEY"))

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to complete payment confirmation. Status: %d", resp.StatusCode),
		}
	}
	defer resp.Body.Close()

	// Decode the response body into the Gorm model struct
	var confirmationResp entities2.PaymentConfirmationResponse
	if err := json.NewDecoder(resp.Body).Decode(&confirmationResp); err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to parse response: %v", err),
		}
	}

	// Verify ref1, ref2, ref3 match records in the database
	transaction, err := uc.verifyRefs(confirmationResp.Ref1, confirmationResp.Ref2, confirmationResp.Ref3)
	if err != nil {
		return nil, err.(*entitiesDtos.ErrorResponse)
	}

	// Update the payment status to CONFIRMED
	_, err = uc.repo.UpdatePayment(transaction.PaymentID, "CONFIRMED")
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to update payment: %v", err),
		}
	}

	return transaction, nil
}

func (uc *PaymentUseCase) verifyRefs(ref1, ref2, ref3 string) (*entities2.Transaction, error) {
	// Get the transaction from the database
	transaction, err := uc.repo.GetTransaction(ref1, ref2, ref3)
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to get transaction: %v", err),
		}
	}

	// Update the transaction status to CONFIRMED
	update, err := uc.repo.UpdateTransaction(transaction.ID, "CONFIRMED")
	if err != nil {
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to update transaction: %v", err),
		}
	}

	return update, nil
}
