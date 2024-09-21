package dtos

import entities "tln-backend/Entities"

type PromptPayResult struct {
	QRResponse  *PromptPayResponse
	Transaction *entities.Transaction
}
