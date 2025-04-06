package controllers

import (
	"net/http"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
)

type TenantController struct {
	tenantService *services.TenantService
}

func NewTenantController(tenantService *services.TenantService) *TenantController {
	return &TenantController{tenantService: tenantService}
}

func (c *TenantController) CreateTenant(ctx fiber.Ctx) error {
	var tenant models.Tenant
	if err := ctx.Bind().Body(&tenant); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.tenantService.CreateTenant(&tenant); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(tenant)
}

func (c *TenantController) GetTenant(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	tenant, err := c.tenantService.GetTenant(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Tenant not found"})
	}
	return ctx.JSON(tenant)
}

func (c *TenantController) GetTenantsByHouse(ctx fiber.Ctx) error {
	houseId := ctx.Params("houseId")
	tenants, err := c.tenantService.GetTenantsByHouse(houseId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(tenants)
}

func (c *TenantController) UpdateTenant(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var tenant models.Tenant
	if err := ctx.Bind().Body(&tenant); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.tenantService.UpdateTenant(id, &tenant); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(tenant)
}

func (c *TenantController) DeleteTenant(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.tenantService.DeleteTenant(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusNoContent)
}
