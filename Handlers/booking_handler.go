package Handlers

import (
	"github.com/gofiber/fiber/v2"
	"log"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Usecase"
)

type BookingHandler struct {
	useCase *Usecase.BookingUseCase
}

func NewBookingHandler(useCase *Usecase.BookingUseCase) *BookingHandler {
	return &BookingHandler{useCase: useCase}
}

func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req entitiesDtos.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err) // Log the detailed error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Call use case method
	booking, errResponse := h.useCase.CreateBooking(&req)
	if errResponse != nil {
		log.Printf("Failed to create booking: %v", errResponse) // Log the error details
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to create booking",
			"details": errResponse,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Booking created successfully",
		"data":    booking,
	})
}

func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
	var req entitiesDtos.CancelBookingRequest
	log.Printf("Raw request body: %s", c.Body())
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err) // Log the detailed error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	// Log the parsed request for debugging
	log.Printf("Received CancelBookingRequest: %+v", req)

	// Call use case method
	booking, errResponse := h.useCase.CancelBooking(&req)
	if errResponse != nil {
		log.Printf("Failed to Cancel booking: %v", errResponse) // Log the error details

		// Determine appropriate status code based on error code
		statusCode := fiber.StatusInternalServerError
		if errResponse.Code == 400 {
			statusCode = fiber.StatusBadRequest
		} else if errResponse.Code == 409 {
			statusCode = fiber.StatusConflict
		}

		return c.Status(statusCode).JSON(fiber.Map{
			"error": "Failed to Cancel booking",
			"details": fiber.Map{
				"code":    errResponse.Code,
				"message": errResponse.Message,
			},
		})
	}

	// Successfully cancelled booking
	log.Printf("Booking cancelled successfully: %+v", booking)
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "Cancel success",
		"message": "Cancel Booking successfully",
		"data":    booking,
	})
}
