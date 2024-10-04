package Usecase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
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

func (uc *PaymentUseCase) PromptPay(request entities2.Payment, paymentID string) (*entitiesDtos.PromptPayResult, *entitiesDtos.ErrorResponse) {
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

	qrResp, promptPay, err := uc.paymentService.CreateQRCode(token, UUID, request.Price)
	if err != nil {
		log.Printf("Failed to create QR code: %v", err)
		return nil, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: fmt.Sprintf("Failed to create QR code: %v", err),
		}
	}
	return &entitiesDtos.PromptPayResult{
		QRResponse:      qrResp,
		PromptPayDetail: promptPay,
	}, nil
}

//
//func (uc *PaymentUseCase) PaymentConfirmation(request *entities2.ConfirmPayment) (*entities2.PaymentConfirmation, *entitiesDtos.ErrorResponse) {
//	// Get the OAuth token
//	oauthResp, err := uc.paymentService.GetOAuthToken()
//	if err != nil {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    500,
//			Message: fmt.Sprintf("Failed to get OAuth token: %v", err),
//		}
//	}
//
//	// Construct the URL with transaction reference and sending bank
//	transRef := request.TransRef
//	sendingBank := request.SendingBank
//	url := fmt.Sprintf("https://api-sandbox.partners.scb/partners/sandbox/v1/payment/billpayment/transactions/%s?sendingBank=%s", transRef, sendingBank)
//
//	// Create and send the HTTP request
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    500,
//			Message: fmt.Sprintf("Failed to create request: %v", err),
//		}
//	}
//
//	req.Header.Set("Content-Type", "application/json")
//	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", oauthResp.Data.AccessToken))
//	req.Header.Set("requestUId", uuid.New().String())
//	req.Header.Set("resourceOwnerId", os.Getenv("API_KEY"))
//
//	resp, err := http.DefaultClient.Do(req)
//	if err != nil {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    500,
//			Message: fmt.Sprintf("Failed to complete payment confirmation: %v", err),
//		}
//	}
//
//	defer resp.Body.Close()
//
//	// Check for non-OK response status
//	if resp.StatusCode != http.StatusOK {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    resp.StatusCode,
//			Message: fmt.Sprintf("Failed to complete payment confirmation. Status: %d", resp.StatusCode),
//		}
//	}
//
//	// Decode the response body into the response struct
//	var confirmationResp entities2.PaymentConfirmation
//	if err := json.NewDecoder(resp.Body).Decode(&confirmationResp); err != nil {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    500,
//			Message: fmt.Sprintf("Failed to parse response: %v", err),
//		}
//	}
//
//	// Verify ref1, ref2, ref3 match records in the database
//	transaction, err := uc.verifyRefs(confirmationResp.Data.Ref1, confirmationResp.Data.Ref2, confirmationResp.Data.Ref3)
//	if err != nil {
//		return nil, err.(*entitiesDtos.ErrorResponse)
//	}
//
//	// Update the payment status to CONFIRMED
//	if _, err := uc.repo.UpdateTransaction(transaction.PaymentID, entities2.TransactionCompleted); err != nil {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    500,
//			Message: fmt.Sprintf("Failed to update payment: %v", err),
//		}
//	}
//
//	return &confirmationResp, nil
//}
//
//func (uc *PaymentUseCase) verifyRefs(ref1, ref2, ref3 string) (*entities2.Transaction, error) {
//	// Get the transaction from the database
//	transaction, err := uc.repo.GetTransaction(ref1, ref2, ref3)
//	if err != nil {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    500,
//			Message: fmt.Sprintf("Failed to get transaction: %v", err),
//		}
//	}
//
//	// Update the transaction status to CONFIRMED
//	update, err := uc.repo.UpdateTransaction(transaction.ID, "CONFIRMED")
//	if err != nil {
//		return nil, &entitiesDtos.ErrorResponse{
//			Code:    500,
//			Message: fmt.Sprintf("Failed to update transaction: %v", err),
//		}
//	}
//
//	return update, nil
//}

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
