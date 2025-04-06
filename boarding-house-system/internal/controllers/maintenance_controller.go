package controllers

import (
	"net/http"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
)

type MaintenanceController struct {
	maintenanceService *services.MaintenanceService
}

func NewMaintenanceController(maintenanceService *services.MaintenanceService) *MaintenanceController {
	return &MaintenanceController{maintenanceService: maintenanceService}
}

func (c *MaintenanceController) CreateRequest(ctx fiber.Ctx) error {
	var request models.MaintenanceRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.maintenanceService.CreateRequest(&request); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(request)
}

func (c *MaintenanceController) GetRequest(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	request, err := c.maintenanceService.GetRequest(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Request not found"})
	}
	return ctx.JSON(request)
}

func (c *MaintenanceController) GetRequestsByRoom(ctx fiber.Ctx) error {
	roomId := ctx.Params("roomId")
	requests, err := c.maintenanceService.GetRequestsByRoom(roomId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(requests)
}

func (c *MaintenanceController) UpdateRequest(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var request models.MaintenanceRequest
	if err := ctx.Bind().Body(&request); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.maintenanceService.UpdateRequest(id, &request); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(request)
}

func (c *MaintenanceController) UpdateRequestStatus(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var statusUpdate struct {
		Status     string `json:"status"`
		AssignedTo *int   `json:"assigned_to"`
	}
	if err := ctx.Bind().Body(&statusUpdate); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.maintenanceService.UpdateRequestStatus(id, statusUpdate.Status, statusUpdate.AssignedTo); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (c *MaintenanceController) DeleteRequest(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.maintenanceService.DeleteRequest(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusNoContent)
}
