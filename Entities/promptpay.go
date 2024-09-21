package entities

import "time"

type PromptPay struct {
	TransactionID string `json:"transaction_id"` // Unique Transaction ID
	Amount        string `json:"amount"`         // Amount as string (for precise formatting)
	Ref1          string `json:"ref1"`           // Reference 1
	Ref2          string `json:"ref2"`           // Reference 2
	Ref3          string `json:"ref3"`           // Reference 3
	Status        string `json:"status"`         // Status of the PromptPay transaction (e.g., "Pending", "Success")
}

type PaymentConfirmation struct {
	Amount                 float64   `json:"amount"`
	BillPaymentRef1        string    `json:"billPaymentRef1"`
	BillPaymentRef2        string    `json:"billPaymentRef2"`
	BillPaymentRef3        string    `json:"billPaymentRef3"`
	ChannelCode            string    `json:"channelCode"`
	CurrencyCode           int       `json:"currencyCode"`
	PayeeAccountNumber     string    `json:"payeeAccountNumber"`
	PayeeName              string    `json:"payeeName"`
	PayeeProxyId           string    `json:"payeeProxyId"`
	PayeeProxyType         string    `json:"payeeProxyType"`
	PayerAccountNumber     string    `json:"payerAccountNumber"`
	PayerName              string    `json:"payerName"`
	PayerProxyId           string    `json:"payerProxyId"`
	PayerProxyType         string    `json:"payerProxyType"`
	ReceivingBankCode      string    `json:"receivingBankCode"`
	SendingBankCode        string    `json:"sendingBankCode"`
	TransactionDateAndTime time.Time `json:"transactionDateAndTime"`
	TransactionID          string    `json:"transactionId"`
	TransactionType        string    `json:"transactionType"`
}

type ConfirmPayment struct {
	TransRef    string `json:"transRef"`
	SendingBank string `json:"sendingBank"`
}
type PaymentConfirmationResponse struct {
	StatusCode           int       `json:"code"`
	StatusDescription    string    `json:"description"`
	TransRef             string    `json:"transRef"`
	SendingBank          string    `json:"sendingBank"`
	ReceivingBank        string    `json:"receivingBank"`
	TransactionDate      string    `json:"transDate"`
	TransactionTime      string    `json:"transTime"`
	SenderDisplayName    string    `json:"senderDisplayName"`
	SenderName           string    `json:"senderName"`
	SenderProxyType      string    `json:"senderProxyType"`
	SenderProxyValue     string    `json:"senderProxyValue"`
	SenderAccountType    string    `json:"senderAccountType"`
	SenderAccountValue   string    `json:"senderAccountValue"`
	ReceiverDisplayName  string    `json:"receiverDisplayName"`
	ReceiverName         string    `json:"receiverName"`
	ReceiverProxyType    string    `json:"receiverProxyType"`
	ReceiverProxyValue   string    `json:"receiverProxyValue"`
	ReceiverAccountType  string    `json:"receiverAccountType"`
	ReceiverAccountValue string    `json:"receiverAccountValue"`
	Amount               string    `json:"amount"`
	PaidLocalAmount      string    `json:"paidLocalAmount"`
	PaidLocalCurrency    string    `json:"paidLocalCurrency"`
	CountryCode          string    `json:"countryCode"`
	Ref1                 string    `json:"ref1"`
	Ref2                 string    `json:"ref2"`
	Ref3                 string    `json:"ref3"`
	CreatedAt            time.Time `json:"createdAt"` // Optional: Add timestamps for tracking
	UpdatedAt            time.Time `json:"updatedAt"`
}
