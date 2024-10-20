package App

import (
	"github.com/joho/godotenv"
	"gorm.io/gorm"
	"log"
	"os"
	"tln-backend/Config"
	"tln-backend/Database"
	"tln-backend/Handlers"
	"tln-backend/Repository"
	"tln-backend/Server"
	"tln-backend/Services"
	"tln-backend/Usecase"
)

func LoadConfig() (*Config.Configs, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return nil, err
	}

	host := os.Getenv("FIBER_HOST")
	if host == "" {
		host = "localhost"
	}

	port := os.Getenv("FIBER_PORT")
	if port == "" {
		port = "3000"
	}

	return &Config.Configs{
		App: Config.AppConfig{
			Host: host,
			Port: port,
		},
	}, nil
}

func InitializeDatabase() (*gorm.DB, error) {
	return Database.NewDB()
}

func InitializeServer(userRepo *Repository.UserRepository, providerRepo *Repository.ProviderRepository) *Server.Server {
	return Server.NewServer(userRepo, providerRepo)
}

func InitializeHandlers(db *gorm.DB) (*Handlers.AllHandlers, *Repository.UserRepository, *Repository.ProviderRepository, error) {
	hashService := Services.NewHashService()
	paymentService := Services.NewPaymentService()

	userRepo := Repository.NewUserRepository(db)
	userUseCase := Usecase.NewUserUseCase(userRepo)
	userHandler := Handlers.NewUserHandler(userUseCase)

	authRepo := Repository.NewAuthRepository(db, hashService)
	authUseCase := Usecase.NewAuthUseCase(authRepo, hashService)
	authHandler := Handlers.NewAuthHandler(authUseCase)

	paymentRepo := Repository.NewPaymentRepository(db)
	paymentUseCase := Usecase.NewPaymentUseCase(paymentRepo, paymentService)
	paymentHandler := Handlers.NewPaymentHandler(paymentUseCase)

	providerRepo := Repository.NewProviderRepository(db)
	providerUseCase := Usecase.NewProviderUseCase(providerRepo)
	providerHandler := Handlers.NewMarketProvider(providerUseCase)

	marketRepo := Repository.NewMarketRepository(db)
	marketUseCase := Usecase.NewMarketUseCase(marketRepo)
	marketHandler := Handlers.NewMarketHandler(marketUseCase)

	slotRepo := Repository.NewSlotRepository(db)
	slotUseCase := Usecase.NewSlotUseCase(slotRepo)
	slotHandler := Handlers.NewSlotHandler(slotUseCase)

	bookingRepo := Repository.NewBookingRepository(db)
	bookingService := Services.NewBookingService(bookingRepo, paymentRepo, slotUseCase)
	bookingUseCase := Usecase.NewBookingUseCase(bookingRepo, paymentRepo, paymentUseCase, bookingService)
	bookingHandler := Handlers.NewBookingHandler(bookingUseCase)

	allHandlers := &Handlers.AllHandlers{
		UserHandler:    userHandler,
		AuthHandler:    authHandler,
		PaymentHandler: paymentHandler,
		MarketProvider: providerHandler,
		MarketHandler:  marketHandler,
		BookingHandler: bookingHandler,
		SlotHandler:    slotHandler,
	}

	return allHandlers, userRepo, providerRepo, nil
}

func StartServer(server *Server.Server, address string) {
	log.Printf("Server is running on %s ⚡", address)
	if err := server.App.Listen(address); err != nil {
		log.Fatal(err)
	}
}
