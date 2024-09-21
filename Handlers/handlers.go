package Handlers

type AllHandlers struct {
	UserHandler    *UserHandler
	AuthHandler    *AuthHandler
	PaymentHandler *PaymentHandler
	MarketProvider *MarketProvider
	MarketHandler  *MarketHandler
	BookingHandler *BookingHandler
	SlotHandler    *SlotHandler
}
