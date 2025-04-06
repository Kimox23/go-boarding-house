package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/config"
	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	userService *services.UserService
	cfg         *config.Config
}

func NewAuthController(userService *services.UserService, cfg *config.Config) *AuthController {
	return &AuthController{
		userService: userService,
		cfg:         cfg,
	}
}

func (c *AuthController) Register(ctx fiber.Ctx) error {
	var input struct {
		Username string `json:"username" validate:"required,min=3"`
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required,min=6"`
		Role     string `json:"role" validate:"required,oneof=admin manager staff tenant"`
		Phone    string `json:"phone" validate:"required"`
	}

	if err := ctx.Bind().Body(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	// Check if email already exists
	existingUser, err := c.userService.GetUserByEmail(input.Email)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}
	if existingUser != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{"error": "Email already in use"})
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	user := &models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     input.Role,
		Phone:    input.Phone,
	}

	if err := c.userService.CreateUser(user); err != nil {
		return ctx.Status(http.StatusConflict).JSON(fiber.Map{"error": "User already exists"})
	}

	return ctx.Status(http.StatusCreated).JSON(user)
}

func (c *AuthController) Login(ctx fiber.Ctx) error {
	var input struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}

	if err := ctx.Bind().Body(&input); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	user, err := c.userService.GetUserByEmail(input.Email)
	if err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return ctx.Status(http.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token, err := config.GenerateJWT(user.ID, user.Role, c.cfg.JWTSecret)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Could not generate token"})
	}

	return ctx.JSON(fiber.Map{
		"token": token,
		"user": fiber.Map{
			"id":    user.ID,
			"email": user.Email,
			"role":  user.Role,
		},
	})
}

func (c *AuthController) Me(ctx fiber.Ctx) error {
	// Get the raw value from context first
	rawUserID := ctx.Locals("userID")
	if rawUserID == nil {
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "User ID not found in request context",
		})
	}

	// Type conversion with comprehensive checking
	var userID int
	switch v := rawUserID.(type) {
	case int:
		userID = v
	case float64:
		userID = int(v)
	case json.Number:
		parsed, err := v.Int64()
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID format in token",
			})
		}
		userID = int(parsed)
	case string:
		parsed, err := strconv.Atoi(v)
		if err != nil {
			return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user ID string format",
			})
		}
		userID = parsed
	default:
		return ctx.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error":   "Unsupported user ID type",
			"details": fmt.Sprintf("Expected number, got %T", rawUserID),
		})
	}

	// Debug log the successful extraction
	log.Printf("Successfully extracted userID: %d", userID)

	// Rest of your controller logic...
	user, err := c.userService.GetUser(strconv.Itoa(userID))
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	profile, _ := c.userService.GetProfile(strconv.Itoa(userID))

	return ctx.JSON(fiber.Map{
		"user":    user,
		"profile": profile,
	})
}
