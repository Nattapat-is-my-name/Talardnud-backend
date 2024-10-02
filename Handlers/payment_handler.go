package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	entities2 "tln-backend/Entities"
	"tln-backend/Usecase"
)

type PaymentHandler struct {
	useCase *Usecase.PaymentUseCase
}

func NewPaymentHandler(useCase *Usecase.PaymentUseCase) *PaymentHandler {
	return &PaymentHandler{useCase: useCase}
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

func (ph *PaymentHandler) ScbConfirmation(c *fiber.Ctx) error {
	var request entities2.PaymentConfirmation

	log.Print("ScbConfirmation")

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input", "details": err.Error()})
	}

	response, err := ph.useCase.PaymentConfirmation(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error", "details": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
