package Server

import (
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"tln-backend/Handlers"
	middleware "tln-backend/Middlewares"
	"tln-backend/Repository"
)

type Server struct {
	App      *fiber.App
	UserRepo *Repository.UserRepository
}

func NewServer(userRepo *Repository.UserRepository) *Server {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	return &Server{App: app, UserRepo: userRepo}
}

func (s *Server) MapHandlers(allHandlers *Handlers.AllHandlers) {
	authMiddleware := middleware.JWTAuthMiddleware(s.UserRepo)

	//create endpoint

	v1 := s.App.Group("/api/v1")

	userGroup := v1.Group("/Users", authMiddleware)
	userGroup.Delete("/:id", allHandlers.UserHandler.DeleteUser)
	userGroup.Get("/:id", allHandlers.UserHandler.GetUserByID)

	providerGroup := v1.Group("/Providers")
	providerGroup.Post("/create", allHandlers.MarketProvider.CreateProvider)
	providerGroup.Put("/update", allHandlers.MarketProvider.UpdateProvider)

	marketGroup := v1.Group("/Markets")
	marketGroup.Post("/create", allHandlers.MarketHandler.CreateMarket)
	marketGroup.Get("/get", allHandlers.MarketHandler.GetMarket)
	marketGroup.Get("/get/:id", allHandlers.MarketHandler.GetMarketByID)

	authGroup := v1.Group("/Auth")
	authGroup.Post("/register", allHandlers.AuthHandler.Register)
	authGroup.Post("/login", allHandlers.AuthHandler.Login)

	bookingGroup := v1.Group("/Bookings", authMiddleware)
	bookingGroup.Post("/create", allHandlers.BookingHandler.CreateBooking)
	bookingGroup.Delete("/cancel", allHandlers.BookingHandler.CancelBooking)

	slotGroup := v1.Group("/Slots")
	slotGroup.Post("/create", allHandlers.SlotHandler.CreateSlot)
	slotGroup.Get("/get/:id", allHandlers.SlotHandler.GetSlot)

	//paymentGroup := v1.Group("/Payments", authMiddleware)
	//paymentGroup.Post("/promptPay", allHandlers.PaymentHandler.PromptPay)

	ScbResponseGroup := v1.Group("/Scb")
	ScbResponseGroup.Post("/confirm", allHandlers.PaymentHandler.ScbConfirmation)

	testGroup := v1.Group("/test")
	testGroup.Get("/info", s.getTestInfo)

	// Serve Swagger documentation
	s.App.Get("/swagger/*", swagger.HandlerDefault)
}

func (s *Server) getTestInfo(c *fiber.Ctx) error {
	//userID := c.Locals("userID")
	//userEmail := c.Locals("email")
	//exp := c.Locals("exp")
	//
	//timestamp := exp.(float64)
	//
	//// Convert timestamp to time.Time
	//t := time.Unix(int64(timestamp), 0)
	//
	//format := "2006-01-02 15:04:05"
	//
	//if userID == nil {
	//	return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID not found in token"})
	//}
	//return c.JSON(fiber.Map{
	//	"userID": userID,
	//	"email":  userEmail,
	//	"exp":    t.Format(format),
	//})

	return c.JSON(fiber.Map{
		"message": "Hello, World!",
	})
}
