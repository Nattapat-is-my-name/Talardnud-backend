package dtos

type PromptPayRequest struct {
	QRType string `json:"qrType"`
	PPType string `json:"ppType"`
	PPId   string `json:"ppId"`
	Amount string `json:"amount"`
	Ref1   string `json:"ref1"`
	Ref2   string `json:"ref2"`
	Ref3   string `json:"ref3"`
}

type PromptPayResponse struct {
	Status struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	} `json:"status"`
	Data struct {
		QrRawData string `json:"qrRawData"`
		QRImage   string `json:"qrImage"`
	} `json:"data"`
}
