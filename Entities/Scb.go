package entities

import "time"

type Scb struct {
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
