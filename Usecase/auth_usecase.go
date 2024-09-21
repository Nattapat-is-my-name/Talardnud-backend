package Usecase

import (
	"github.com/google/uuid"
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
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
func (uc *AuthUseCase) Register(username, password, email string) *entitiesDtos.ErrorResponse {
	// Check if the username already exists
	exists, errResponse := uc.repo.IsUsernameAndEmailExists(username, email)
	if errResponse != nil {
		return errResponse
	}
	if exists {
		return &entitiesDtos.ErrorResponse{
			Code:    409,
			Message: "Username or email already exists",
		}
	}

	// Hash the password
	hashedPassword, err := uc.auth.HashPassword(password)
	if err != nil {
		return &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Error hashing password: " + err.Error(),
		}
	}

	// Create a new vendor
	newVendor := entities.Vendor{
		ID:       uuid.NewString(),
		Username: username,
		Email:    email,
		Password: hashedPassword,
	}

	// Register the new vendor
	if err := uc.repo.Register(&newVendor); err != nil {
		return &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Error registering user: " + err.Error(),
		}
	}

	return nil
}
