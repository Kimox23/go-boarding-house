package controllers

import (
	"net/http"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"
	"github.com/Kimox23/boarding-house-app/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type UserController struct {
	userService *services.UserService
}

func NewUserController(userService *services.UserService) *UserController {
	return &UserController{userService: userService}
}

// CreateUser handles user registration
func (c *UserController) CreateUser(ctx fiber.Ctx) error {
	var user models.User
	if err := ctx.Bind().Body(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.userService.CreateUser(&user); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(user)
}

// GetUser retrieves a user by ID
func (c *UserController) GetUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	user, err := c.userService.GetUser(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
	}
	return ctx.JSON(user)
}

// GetAllUsers retrieves all users
func (c *UserController) GetAllUsers(ctx fiber.Ctx) error {
	pagination := utils.GetPagination(ctx)
	users, total, err := c.userService.GetAllUsers(pagination)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(fiber.Map{
		"data": users,
		"meta": fiber.Map{
			"page":        pagination.Page,
			"page_size":   pagination.PageSize,
			"total":       total,
			"total_pages": (total + pagination.PageSize - 1) / pagination.PageSize,
		},
	})
}

// UpdateUser updates user information
func (c *UserController) UpdateUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var user models.User
	if err := ctx.Bind().Body(&user); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.userService.UpdateUser(id, &user); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(user)
}

// DeleteUser deletes a user
func (c *UserController) DeleteUser(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.userService.DeleteUser(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusNoContent)
}

// UserProfile CRUD operations
func (c *UserController) CreateProfile(ctx fiber.Ctx) error {
	userId := ctx.Params("userId")
	var profile models.UserProfile
	if err := ctx.Bind().Body(&profile); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.userService.CreateProfile(userId, &profile); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(profile)
}

func (c *UserController) GetProfile(ctx fiber.Ctx) error {
	userId := ctx.Params("userId")
	profile, err := c.userService.GetProfile(userId)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Profile not found"})
	}
	return ctx.JSON(profile)
}
