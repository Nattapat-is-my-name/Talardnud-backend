package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"os"
	"strings"
	"tln-backend/Repository"
)

func JWTAuthMiddleware(userRepo *Repository.UserRepository, providerRepo *Repository.ProviderRepository) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token format"})
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.NewError(fiber.StatusUnauthorized, "unexpected signing method")
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			userID := claims["sub"].(string)
			email := claims["email"].(string)
			role := claims["role"].(string)

			var _ error
			if role == "provider" {
				// Check if the provider exists in the database
				provider, err := providerRepo.GetProviderByID(userID)
				if err != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Provider not found"})
				}
				// Verify the email
				if provider.Email != email {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid provider information"})
				}
			} else {
				// Check if the user exists in the database
				user, err := userRepo.GetUserByID(userID)
				if err != nil {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
				}
				// Verify the email
				if user.Email != email {
					return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid user information"})
				}
			}

			c.Locals("userID", userID)
			c.Locals("email", email)
			c.Locals("role", role)
			c.Locals("exp", claims["exp"])
			return c.Next()
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}
}

func ProviderAuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		role := c.Locals("role")
		if role != "provider" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Access denied. Provider role required.",
			})
		}
		return c.Next()
	}
}
