package Repository

import (
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"os"
	"time"
	"tln-backend/Entities"
	entitiesDtos "tln-backend/Entities/dtos"
	"tln-backend/Interfaces"
)

type AuthRepository struct {
	db          *gorm.DB
	hashService Interfaces.IHashService
}

func NewAuthRepository(db *gorm.DB, hashService Interfaces.IHashService) *AuthRepository {
	return &AuthRepository{db: db, hashService: hashService}
}

// Login checks the user's credentials and returns a JWT token if successful.
func (repo *AuthRepository) Login(usernameOrEmail, password string) (entities.LoginResponse, error) {
	var user entities.Vendor

	// Find the user by either username or email
	if err := repo.db.Where("username = ? OR email = ?", usernameOrEmail, usernameOrEmail).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entities.LoginResponse{}, errors.New("user not found")
		}
		return entities.LoginResponse{}, err
	}

	// Check if the password matches
	if !checkPasswordHash(password, user.Password) {
		return entities.LoginResponse{}, errors.New("invalid password")
	}

	// Generate a JWT token
	token, err := generateToken(user)
	if err != nil {
		return entities.LoginResponse{}, err
	}

	return entities.LoginResponse{
		AccessToken: token.AccessToken,
		VendorID:    user.ID,
	}, nil

}

// CheckPasswordHash checks if the provided password matches the hashed password.
func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// GenerateToken creates a JWT token for the given user.
func generateToken(user entities.Vendor) (entities.LoginResponse, error) {
	secretKey := []byte(os.Getenv("JWT_SECRET_KEY"))
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub":   user.ID,
		"email": user.Email,
		"role":  "vendor", // Add the role claim
		"exp":   time.Now().Add(24 * time.Hour).Unix(),
		"iat":   time.Now().Unix(),
	})
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return entities.LoginResponse{}, err
	}
	return entities.LoginResponse{AccessToken: tokenString}, nil
}

// Register creates a new vendor in the database.
func (repo *AuthRepository) Register(user *entities.Vendor) error {
	if user.Username == "" || user.Email == "" {
		return errors.New("username and Email are required")
	}
	if err := repo.db.Create(user).Error; err != nil {
		return err
	}
	return nil
}

// IsUsernameAndEmailExists checks if either the username or email is already taken.
func (repo *AuthRepository) IsUsernameAndEmailExists(username, email string) (bool, *entitiesDtos.ErrorResponse) {
	var count int64

	// Query to check if the username or email already exists
	if err := repo.db.Model(&entities.Vendor{}).
		Where("username = ? OR email = ?", username, email).
		Count(&count).Error; err != nil {
		// Handle database error
		return false, &entitiesDtos.ErrorResponse{
			Code:    500,
			Message: "Database error while checking username and email: " + err.Error(),
		}
	}

	// If count is greater than 0, it means either the username or email exists
	if count > 0 {
		return true, nil
	}

	// Neither username nor email exists
	return false, nil
}

func (repo *AuthRepository) ProviderLogin(username, password string) (*entities.MarketProvider, error) {
	var provider entities.MarketProvider

	// Find the provider by username
	if err := repo.db.Where("username = ?", username).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("provider not found")
		}
		return nil, err
	}
	// Check if the password matches
	if !checkPasswordHash(password, provider.Password) {
		return nil, errors.New("invalid password")
	}

	return &provider, nil
}
func (repo *AuthRepository) ProviderRegister(provider *entities.MarketProvider) error {
	if provider.Username == "" || provider.Email == "" {
		return errors.New("username and Email are required")
	}
	if err := repo.db.Create(provider).Error; err != nil {
		return err
	}
	return nil
}
