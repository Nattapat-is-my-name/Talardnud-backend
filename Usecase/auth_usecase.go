package Usecase

import (
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"os"
	"time"
	entities "tln-backend/Entities"
	"tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
)

type AuthUseCase struct {
	repo Interfaces.IAuthRepository
	auth Interfaces.IHashService
}

// NewAuthUseCase initializes an AuthUseCase with the provided repo and auth services
func NewAuthUseCase(repo Interfaces.IAuthRepository, auth Interfaces.IHashService) *AuthUseCase {
	return &AuthUseCase{
		repo: repo,
		auth: auth,
	}
}

// Login handles user login by verifying credentials
func (uc *AuthUseCase) Login(usernameOrEmail, password string) (entities.LoginResponse, error) {
	return uc.repo.Login(usernameOrEmail, password)
}

// Register handles user registration by hashing the password and saving the user data
func (uc *AuthUseCase) Register(username, password, email, phone, firstname, lastname string) (dtos.RegisterResponse, *dtos.ErrorResponse) {
	// Check if the username and email already exist
	exists, errResponse := uc.repo.IsUsernameAndEmailExists(username, email)
	if errResponse != nil {
		return dtos.RegisterResponse{}, errResponse
	}
	if exists {
		return dtos.RegisterResponse{}, &dtos.ErrorResponse{
			Code:    409,
			Message: "Username or email already exists",
		}
	}

	// Hash the password
	hashedPassword, err := uc.auth.HashPassword(password)
	if err != nil {
		return dtos.RegisterResponse{}, &dtos.ErrorResponse{
			Code:    500,
			Message: "Error hashing password: " + err.Error(),
		}
	}

	// Create a new vendor
	newVendor := entities.Vendor{
		ID:        uuid.NewString(),
		Username:  username,
		Password:  hashedPassword,
		Email:     email,
		FirstName: firstname,
		LastName:  lastname,
		Phone:     phone,
	}

	// Register the new vendor
	if err := uc.repo.Register(&newVendor); err != nil {
		return dtos.RegisterResponse{}, &dtos.ErrorResponse{
			Code:    500,
			Message: "Failed to register user: " + err.Error(),
		}
	}

	// Convert to RegisterResponse
	response := dtos.RegisterResponse{
		ID:        newVendor.ID,
		Username:  newVendor.Username,
		Email:     newVendor.Email,
		Phone:     newVendor.Phone,
		Firstname: newVendor.FirstName,
		Lastname:  newVendor.LastName,
	}

	return response, nil
}

func (uc *AuthUseCase) ProviderLogin(username, password string) (dtos.ProviderLoginResponse, error) {
	// Authenticate the provider
	provider, err := uc.repo.ProviderLogin(username, password)
	if err != nil {
		return dtos.ProviderLoginResponse{}, err
	}

	// Generate JWT token
	tokenString, err := uc.generateProviderToken(provider)
	if err != nil {
		return dtos.ProviderLoginResponse{}, err
	}

	return dtos.ProviderLoginResponse{
		AccessToken: tokenString,
		ProviderID:  provider.ID,
	}, nil
}

func (uc *AuthUseCase) generateProviderToken(provider *entities.MarketProvider) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["sub"] = provider.ID
	claims["email"] = provider.Email
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	claims["role"] = "provider"
	claims["iat"] = time.Now().Unix()

	return token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
}

func (uc *AuthUseCase) RegisterProvider(username, phone, email, password string) (entities.MarketProvider, *dtos.ErrorResponse) {
	// Check if the email already exists
	exists, errResponse := uc.repo.IsUsernameAndEmailExists(username, email)
	if errResponse != nil {
		return entities.MarketProvider{}, errResponse
	}
	if exists {
		return entities.MarketProvider{}, &dtos.ErrorResponse{
			Code:    409,
			Message: "Provider email already exists",
		}
	}

	hashedPassword, err := uc.auth.HashPassword(password)
	if err != nil {
		return entities.MarketProvider{}, &dtos.ErrorResponse{
			Code:    500,
			Message: "Error hashing password: " + err.Error(),
		}
	}
	// Create a new market provider
	newProvider := entities.MarketProvider{
		ID:       uuid.NewString(),
		Username: username,
		Password: hashedPassword,
		Phone:    phone,
		Email:    email,
	}

	// Register the new provider
	if err := uc.repo.ProviderRegister(&newProvider); err != nil {
		return entities.MarketProvider{}, &dtos.ErrorResponse{
			Code:    500,
			Message: "Failed to register provider: " + err.Error(),
		}
	}

	// Convert to ProviderRegisterResponse
	response := entities.MarketProvider{
		ID:       newProvider.ID,
		Username: newProvider.Username,
		Phone:    newProvider.Phone,
		Email:    newProvider.Email,
	}

	return response, nil
}
