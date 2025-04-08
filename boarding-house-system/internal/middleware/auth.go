package middleware

import (
	"log"
	"strings"

	"github.com/Kimox23/boarding-house-app/internal/config"

	"github.com/gofiber/fiber/v3"
)

func AuthRequired(cfg *config.Config) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		authHeader := ctx.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization header missing")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authorization header is required",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			log.Println("Bearer prefix missing in token")
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Bearer token not found",
			})
		}

		log.Printf("Verifying token: %.20s...", tokenString) // Log first 20 chars

		claims, err := config.VerifyJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			log.Printf("Token verification failed: %v", err)
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   "Invalid token",
				"details": err.Error(),
			})
		}

		log.Printf("Authenticated user %d with role %s",
			claims.UserID, claims.Role)

		// Store with explicit type
		ctx.Locals("userID", claims.UserID) // Force int type
		ctx.Locals("userRole", claims.Role)

		return ctx.Next()
	}
}

func RoleRequired(requiredRole string, cfg *config.Config) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		userRole, ok := ctx.Locals("userRole").(string)
		if !ok || userRole != requiredRole {
			return ctx.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "Insufficient permissions",
			})
		}
		return ctx.Next()
	}
}
