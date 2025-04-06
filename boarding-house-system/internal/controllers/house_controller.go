package controllers

import (
	"net/http"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
)

type HouseController struct {
	houseService *services.HouseService
}

func NewHouseController(houseService *services.HouseService) *HouseController {
	return &HouseController{houseService: houseService}
}

func (c *HouseController) CreateHouse(ctx fiber.Ctx) error {
	var house models.BoardingHouse
	if err := ctx.Bind().Body(&house); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.houseService.CreateHouse(&house); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(house)
}

func (c *HouseController) GetHouse(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	house, err := c.houseService.GetHouse(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "House not found"})
	}
	return ctx.JSON(house)
}

func (c *HouseController) GetAllHouses(ctx fiber.Ctx) error {
	houses, err := c.houseService.GetAllHouses()
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(houses)
}

func (c *HouseController) UpdateHouse(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var house models.BoardingHouse
	if err := ctx.Bind().Body(&house); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.houseService.UpdateHouse(id, &house); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(house)
}

func (c *HouseController) DeleteHouse(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.houseService.DeleteHouse(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusNoContent)
}
