package entities

type PromptPay struct {
	TransactionID string `json:"transaction_id"` // Unique Transaction ID
	Amount        string `json:"amount"`         // Amount as string (for precise formatting)
	Ref1          string `json:"ref1"`           // Reference 1
	Ref2          string `json:"ref2"`           // Reference 2
	Ref3          string `json:"ref3"`           // Reference 3
	Status        string `json:"status"`         // Status of the PromptPay transaction (e.g., "Pending", "Success")
}

type PaymentConfirmation struct {
	Amount                 string `json:"amount"`
	BillPaymentRef1        string `json:"billPaymentRef1"`
	BillPaymentRef2        string `json:"billPaymentRef2"`
	BillPaymentRef3        string `json:"billPaymentRef3"`
	ChannelCode            string `json:"channelCode"`
	CurrencyCode           string `json:"currencyCode"`
	PayeeAccountNumber     string `json:"payeeAccountNumber"`
	PayeeName              string `json:"payeeName"`
	PayeeProxyId           string `json:"payeeProxyId"`
	PayeeProxyType         string `json:"payeeProxyType"`
	PayerAccountNumber     string `json:"payerAccountNumber"`
	PayerName              string `json:"payerName"`
	PayerProxyId           string `json:"payerProxyId"`
	PayerProxyType         string `json:"payerProxyType"`
	ReceivingBankCode      string `json:"receivingBankCode"`
	SendingBankCode        string `json:"sendingBankCode"`
	TransactionDateAndTime string `json:"transactionDateandTime"`
	TransactionId          string `json:"transactionId"`
	TransactionType        string `json:"transactionType"`
}

type BillPayment struct {
	Status Status      `json:"status"`
	Data   PaymentData `json:"data"`
}

type Status struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
}

type PaymentData struct {
	TransRef          string `json:"transRef"`
	SendingBank       string `json:"sendingBank"`
	ReceivingBank     string `json:"receivingBank"`
	TransDate         string `json:"transDate"`
	TransTime         string `json:"transTime"`
	Sender            User   `json:"sender"`
	Receiver          User   `json:"receiver"`
	Amount            string `json:"amount"`
	PaidLocalAmount   string `json:"paidLocalAmount"`
	PaidLocalCurrency string `json:"paidLocalCurrency"`
	CountryCode       string `json:"countryCode"`
	Ref1              string `json:"ref1"`
	Ref2              string `json:"ref2"`
	Ref3              string `json:"ref3"`
}

type User struct {
	DisplayName string  `json:"displayName"`
	Name        string  `json:"name"`
	Proxy       Proxy   `json:"proxy"`
	Account     Account `json:"account"`
}

type Proxy struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Account struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}
