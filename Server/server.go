package Server

import (
	"tln-backend/Handlers"
	middleware "tln-backend/Middlewares"
	"tln-backend/Repository"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

type Server struct {
	App          *fiber.App
	UserRepo     *Repository.UserRepository
	ProviderRepo *Repository.ProviderRepository
}

func NewServer(userRepo *Repository.UserRepository, providerRepo *Repository.ProviderRepository) *Server {
	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())
	return &Server{
		App:          app,
		UserRepo:     userRepo,
		ProviderRepo: providerRepo,
	}
}

func (s *Server) MapHandlers(allHandlers *Handlers.AllHandlers) {
	authMiddleware := middleware.JWTAuthMiddleware(s.UserRepo, s.ProviderRepo)
	providerMiddleware := middleware.ProviderAuthMiddleware()

	v1 := s.App.Group("/api/v1")

	userGroup := v1.Group("/Users", authMiddleware)
	userGroup.Delete("/:id", allHandlers.UserHandler.DeleteUser)
	userGroup.Get("/:id", allHandlers.UserHandler.GetUserByID)
	//userGroup.Patch("/:id", allHandlers.UserHandler.UpdateUser)

	providerGroup := v1.Group("/Providers", authMiddleware, providerMiddleware)

	providerGroup.Put("/update", allHandlers.MarketProvider.UpdateProvider)

	marketGroup := v1.Group("/Markets")
	marketGroup.Post("/create", allHandlers.MarketHandler.CreateMarket, providerMiddleware)
	marketGroup.Get("/get", allHandlers.MarketHandler.GetMarket)
	marketGroup.Patch("/edit/:id", allHandlers.MarketHandler.EditMarket, providerMiddleware)
	marketGroup.Get("/get/:id", allHandlers.MarketHandler.GetMarketByID)
	marketGroup.Get("/provider/get/:id", allHandlers.MarketHandler.GetMarketByProviderID, providerMiddleware)

	authGroup := v1.Group("/Auth")
	authGroup.Post("/register", allHandlers.AuthHandler.Register)
	authGroup.Post("/login", allHandlers.AuthHandler.Login)
	authGroup.Post("/provider/login", allHandlers.AuthHandler.ProviderLogin)
	authGroup.Post("/provider/register", allHandlers.AuthHandler.RegisterProvider)

	bookingGroup := v1.Group("/Bookings")
	bookingGroup.Post("/create", allHandlers.BookingHandler.CreateBooking)
	bookingGroup.Get("/get/:id", allHandlers.BookingHandler.GetBooking)
	bookingGroup.Get("/user/:id", allHandlers.BookingHandler.GetBookingsByUser)
	bookingGroup.Patch("/cancel", allHandlers.BookingHandler.CancelBooking)

	slotGroup := v1.Group("/Slots")
	slotGroup.Post("/:marketId/create", allHandlers.SlotHandler.CreateOrUpdateLayout, providerMiddleware)
	slotGroup.Get("/get/:id", allHandlers.SlotHandler.GetSlot)
	slotGroup.Patch("/edit/:id", allHandlers.SlotHandler.EditSlot, providerMiddleware)
	slotGroup.Delete("/delete/:id", allHandlers.SlotHandler.DeleteSlot, providerMiddleware)
	slotGroup.Get("/provider/get/:id", allHandlers.SlotHandler.GetProviderSlots, providerMiddleware)
	slotGroup.Get("/markets/:marketID/date/:date", allHandlers.SlotHandler.GetSlotByDate, providerMiddleware)
	slotGroup.Delete("/delete/:id/zone/:zoneID/date/:date", allHandlers.SlotHandler.DeleteSlotByDateAndZone, providerMiddleware)

	paymentGroup := v1.Group("/Payments", authMiddleware)
	paymentGroup.Get("/get/:id", allHandlers.PaymentHandler.GetPayment)
	//paymentGroup.Post("/promptPay", allHandlers.PaymentHandler.PromptPay)

	dashboardGroup := v1.Group("/Dashboard")
	//dashboardGroup.Get("/", allHandlers.DashboardHandler.GetDashboardData) // Get all markets
	dashboardGroup.Get("/get/:id", allHandlers.DashboardHandler.GetSingleMarketStats)
	dashboardGroup.Get("/weekly/:id", allHandlers.DashboardHandler.GetWeeklyStats)

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
