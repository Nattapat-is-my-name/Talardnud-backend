package Interfaces

type IHashService interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword string, password string) error
}
