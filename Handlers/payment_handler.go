package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	entities "tln-backend/Entities"
	"tln-backend/Usecase"
)

type PaymentHandler struct {
	useCase *Usecase.PaymentUseCase
}

func NewPaymentHandler(useCase *Usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
}

func (ph *PaymentHandler) ScbConfirmation(c *fiber.Ctx) error {
	var request entities.PaymentConfirmation

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input", "details": err.Error()})
	}

	response, err := ph.useCase.PaymentConfirmation(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

//func (ph *PaymentHandler) PromptPay(c *fiber.Ctx) error {
//	var request entities.BookingRequest
//
//	if err := c.BodyParser(&request); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
//	}
//
//	response, err := ph.useCase.PromptPay(&request )
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error"})
//	}
//
//	fmt.Print(response)
//
//	return c.Status(fiber.StatusOK).JSON(response)
//}

//func (ph *PaymentHandler) ScbConfirmation(c *fiber.Ctx) error {
//	var request entities2.ConfirmPayment
//
//	if err := c.BodyParser(&request); err != nil {
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input", "details": err.Error()})
//	}
//	log.Print("cdsc", request)
//	response, err := ph.useCase.PaymentConfirmation(&request)
//	if err != nil {
//		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error", "details": err.Error()})
//	}
//
//	return c.Status(fiber.StatusOK).JSON(response)
//}

// GetPayment godoc
// @Summary Get payment by ID
// @Description Get payment by the provided ID
// @Tags payments
// @Accept  json
// @Produce  json
// @Param id path string true "Payment ID"
// @Success 200 {object} dtos.BookingResponse
// @Failure 404 {object} string "Payment not found"
// @Failure 500 {object} string "Internal server error"
// @Router /payments/get/{id} [get]
func (ph *PaymentHandler) GetPayment(c *fiber.Ctx) error {
	paymentID := c.Params("id")
	payment, errResponse := ph.useCase.GetPayment(paymentID)
	if errResponse != nil {
		log.Printf("Failed to get payment with ID %s: %v", paymentID, errResponse) // Log the error details
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get payment",
			"details": errResponse,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Payment retrieved successfully",
		"data":    payment,
	})

}
