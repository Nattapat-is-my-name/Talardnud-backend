package Services

import (
	"golang.org/x/crypto/bcrypt"
)

type HashService struct{}

func NewHashService() *HashService {
	return &HashService{}
}

func (h *HashService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func (h *HashService) CompareHashAndPassword(hashedPassword string, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
