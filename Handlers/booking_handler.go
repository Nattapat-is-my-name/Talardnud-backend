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

// CreateBooking godoc
// @Summary Create a booking
// @Description Create a new booking with the provided data
// @Tags bookings
// @Accept  json
// @Produce  json
// @Param booking body dtos.BookingRequest true "Booking data"
// @Success 200 {object} dtos.BookingResponse
// @Failure 400 {object} string "Invalid input"
// @Failure 409 {object} string "Booking already exists"
// @Failure 500 {object} string "Internal server error"
// @Router /bookings/create [post]
func (h *BookingHandler) CreateBooking(c *fiber.Ctx) error {
	var req entitiesDtos.BookingRequest
	if err := c.BodyParser(&req); err != nil {
		log.Printf("Failed to parse request body: %v", err) // Log the detailed error
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

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

// GetBookingsByUser godoc
// @Summary Get bookings by user
// @Description Get bookings by user with the provided ID
// @Tags bookings
// @Accept  json
// @Produce  json
// @Param id path string true "User ID"
// @Success 200 {object} []entities.Booking
// @Failure 404 {object} string "Bookings not found"
// @Failure 500 {object} string "Internal server error"
// @Router /bookings/user/{id} [get]
func (h *BookingHandler) GetBookingsByUser(c *fiber.Ctx) error {
	userId := c.Params("id")
	bookings, errResponse := h.useCase.GetBookingsByUser(userId)
	if errResponse != nil {
		log.Printf("Failed to get bookings for user with ID %s: %v", userId, errResponse) // Log the error details
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get bookings",
			"details": errResponse,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Bookings retrieved successfully",
		"data":    bookings,
	})
}

// CancelBooking godoc
// @Summary Cancel a booking
// @Description Cancel a booking with the provided data
// @Tags bookings
// @Accept  json
// @Produce  json
// @Param booking body entitiesDtos.CancelBookingRequest true "Booking data"
// @Success 200 {object} entitiesDtos.BookingResponse
// @Failure 400 {object} string "Invalid input"
// @Failure 409 {object} string "Booking already exists"
// @Failure 500 {object} string "Internal server error"
// @Router /bookings/cancel [post]
//func (h *BookingHandler) CancelBooking(c *fiber.Ctx) error {
//	var req entitiesDtos.CancelBookingRequest
//	log.Printf("Raw request body: %s", c.Body())
//	if err := c.BodyParser(&req); err != nil {
//		log.Printf("Failed to parse request body: %v", err) // Log the detailed error
//		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
//			"error":   "Invalid request body",
//			"details": err.Error(),
//		})
//	}
//
//	// Log the parsed request for debugging
//	log.Printf("Received CancelBookingRequest: %+v", req)
//
//	// Call use case method
//	booking, errResponse := h.useCase.CancelBooking(&req)
//	if errResponse != nil {
//		log.Printf("Failed to Cancel booking: %v", errResponse) // Log the error details
//
//		// Determine appropriate status code based on error code
//		statusCode := fiber.StatusInternalServerError
//		if errResponse.Code == 400 {
//			statusCode = fiber.StatusBadRequest
//		} else if errResponse.Code == 409 {
//			statusCode = fiber.StatusConflict
//		}
//
//		return c.Status(statusCode).JSON(fiber.Map{
//			"error": "Failed to Cancel booking",
//			"details": fiber.Map{
//				"code":    errResponse.Code,
//				"message": errResponse.Message,
//			},
//		})
//	}
//
//	// Successfully cancelled booking
//	log.Printf("Booking cancelled successfully: %+v", booking)
//	return c.Status(fiber.StatusOK).JSON(fiber.Map{
//		"status":  "Cancel success",
//		"message": "Cancel Booking successfully",
//		"data":    booking,
//	})
//}

// GetBooking godoc
// @Summary Get a booking
// @Description Get a booking with the provided ID
// @Tags bookings
// @Accept  json
// @Produce  json
// @Param id path string true "Booking ID"
// @Success 200 {object} dtos.BookingResponse
// @Failure 404 {object} string "Booking not found"
// @Failure 500 {object} string "Internal server error"
// @Router /bookings/get/{id} [get]
func (h *BookingHandler) GetBooking(c *fiber.Ctx) error {
	bookingID := c.Params("id")
	booking, errResponse := h.useCase.GetBooking(bookingID)
	if errResponse != nil {
		log.Printf("Failed to get booking with ID %s: %v", bookingID, errResponse) // Log the error details
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":   "Failed to get booking",
			"details": errResponse,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Booking retrieved successfully",
		"data":    booking,
	})
}
