package controllers

import (
	"net/http"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
)

type RoomController struct {
	roomService *services.RoomService
}

func NewRoomController(roomService *services.RoomService) *RoomController {
	return &RoomController{roomService: roomService}
}

// CreateRoom creates a new room
func (c *RoomController) CreateRoom(ctx fiber.Ctx) error {
	var room models.Room
	if err := ctx.Bind().Body(&room); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := c.roomService.CreateRoom(&room); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create room",
		})
	}

	return ctx.Status(http.StatusCreated).JSON(room)
}

// GetRoom retrieves a room by ID
func (c *RoomController) GetRoom(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID is required",
		})
	}

	room, err := c.roomService.GetRoom(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "Room not found",
		})
	}

	return ctx.JSON(room)
}

// GetAllRooms retrieves all rooms
func (c *RoomController) GetAllRooms(ctx fiber.Ctx) error {
	rooms, err := c.roomService.GetAllRooms()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch rooms",
		})
	}

	return ctx.JSON(rooms)
}

// UpdateRoom updates an existing room
func (c *RoomController) UpdateRoom(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID is required",
		})
	}

	var room models.Room
	if err := ctx.Bind().Body(&room); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := c.roomService.UpdateRoom(id, &room); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to update room",
		})
	}

	return ctx.JSON(room)
}

// DeleteRoom deletes a room
func (c *RoomController) DeleteRoom(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Room ID is required",
		})
	}

	if err := c.roomService.DeleteRoom(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete room",
		})
	}

	return ctx.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Room deleted successfully",
	})
}
