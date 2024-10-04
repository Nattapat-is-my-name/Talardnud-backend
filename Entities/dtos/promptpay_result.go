package dtos

import entities "tln-backend/Entities"

type PromptPayResult struct {
	QRResponse      *PromptPayResponse
	PromptPayDetail *entities.PromptPay
}
