package dtos

import entities "tln-backend/Entities"

type MarketResponse struct {
	Status  string            `json:"status"`  // success or error
	Message string            `json:"message"` // message to accompany the response
	Data    []entities.Market `json:"data"`    // the data to be returned
}

type GetListMarketResponse struct {
	Market []MarketResponse `json:"market"`
}
