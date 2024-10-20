package Interfaces

import (
	entities "tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
)

type IAuthRepository interface {
	Login(username, password string) (entities.LoginResponse, error)
	Register(user *entities.Vendor) error
	IsUsernameAndEmailExists(username, email string) (bool, *entitiesDtos.ErrorResponse)
	ProviderLogin(username, password string) (entities.MarketProvider, error)
	ProviderRegister(provider *entities.MarketProvider) error
}
