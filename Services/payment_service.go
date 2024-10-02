package Services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
	entities2 "tln-backend/Entities"
	entities "tln-backend/Entities/dtos"
)

type PaymentService struct {
}

func NewPaymentService() *PaymentService {
	return &PaymentService{}
}

// GetOAuthToken retrieves an OAuth token by making a request to the external service.
func (uc *PaymentService) GetOAuthToken() (*entities.OAuthResponse, error) {
	OAuthURL := os.Getenv("OAUTH_URL")

	reqBody := entities.OAuthRequest{

		ApplicationKey:    os.Getenv("API_KEY"),
		ApplicationSecret: os.Getenv("APPLICATION_KEY"),
	}

	// Marshal the request body into JSON
	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal OAuth request body: %w", err)
	}

	// Generate a UUID for request
	requestID, err := uuid.NewUUID()

	// Create a new HTTP POST request
	req, err := http.NewRequest("POST", OAuthURL, bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create OAuth request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept-language", "EN")
	req.Header.Set("requestUId", requestID.String())
	req.Header.Set("resourceOwnerId", os.Getenv("API_KEY"))

	// Make the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make OAuth request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OAuth response body: %w", err)
	}

	// Check if the status code is 200 (OK)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("received non-200 OAuth response code: %d, body: %s", resp.StatusCode, string(body))
	}

	// Unmarshal the response body into the OAuthResponse struct
	var oauthResp entities.OAuthResponse
	if err := json.Unmarshal(body, &oauthResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OAuth response body: %w", err)
	}

	log.Print("oauthResp", oauthResp)

	oauthResp.Data.UUID = requestID.String()

	// Return the OAuth token response
	return &oauthResp, nil
}
func (uc *PaymentService) CreateQRCode(accessToken string, uuid string, amount float64) (*entities.PromptPayResponse, *entities2.PromptPay, error) {
	amountStr := fmt.Sprintf("%.2f", amount)
	reqBody := entities.PromptPayRequest{
		QRType: "PP",
		PPType: "BILLERID",
		PPId:   os.Getenv("PP_ID"),
		Amount: amountStr,
		Ref1:   generateRef(""),
		Ref2:   generateRef(""),
		Ref3:   generateRef("SCB"),
	}

	jsonReqBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://api-sandbox.partners.scb/partners/sandbox/v1/payment/qrcode/create", bytes.NewBuffer(jsonReqBody))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("accept-language", "EN")
	req.Header.Set("authorization", fmt.Sprintf("Bearer %s", accessToken))
	req.Header.Set("requestUId", uuid)
	req.Header.Set("resourceOwnerId", os.Getenv("API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	log.Print("resp", resp)

	if err != nil {
		return nil, nil, fmt.Errorf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read response body: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("received non-200 response code: %d, body: %s", resp.StatusCode, body)
	}

	var PromptPay entities2.PromptPay
	PromptPay.TransactionID = generateTransactionID()
	PromptPay.Amount = amountStr
	PromptPay.Ref1 = reqBody.Ref1
	PromptPay.Ref2 = reqBody.Ref2
	PromptPay.Ref3 = reqBody.Ref3
	PromptPay.Status = "Pending"

	var qrResp entities.PromptPayResponse
	if err := json.Unmarshal(body, &qrResp); err != nil {
		return nil, nil, fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	return &qrResp, &PromptPay, nil
}

func generateRef(prefix string) string {
	if prefix != "" {
		return prefix + generateReferenceNumber(20-len(prefix))
	}

	length := 20
	return generateReferenceNumber(length)
}
func generateReferenceNumber(length int) string {
	// Available characters for the reference number
	charset := "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Generate the random reference number
	referenceNumber := make([]byte, length)
	for i := range referenceNumber {
		referenceNumber[i] = charset[rand.Intn(len(charset))]
	}

	return string(referenceNumber)
}

func generateTransactionID() string {
	return uuid.New().String()
}
