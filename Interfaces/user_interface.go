package Interfaces

import (
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type IUserRepository interface {
	// CreateUser saves a new user in the repository
	CreateUser(user *entities.Vendor) error

	// GetUserByID retrieves a user by their ID
	GetUserByID(id string) (*entitiesDtos.GetUserResponse, error)

	// UpdateUser updates the details of an existing user
	UpdateUser(user *entities.Vendor) error

	// DeleteUser removes a user from the repository
	DeleteUser(id string) error

	// ListUsers retrieves all users in the repository
	ListUsers() ([]*entities.Vendor, error)
}
