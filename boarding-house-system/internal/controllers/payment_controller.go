package controllers

import (
	"net/http"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController(paymentService *services.PaymentService) *PaymentController {
	return &PaymentController{paymentService: paymentService}
}

func (c *PaymentController) CreatePayment(ctx fiber.Ctx) error {
	var payment models.Payment
	if err := ctx.Bind().Body(&payment); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.paymentService.CreatePayment(&payment); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(payment)
}

func (c *PaymentController) GetPayment(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	payment, err := c.paymentService.GetPayment(id)
	if err != nil {
		return ctx.Status(http.StatusNotFound).JSON(fiber.Map{"error": "Payment not found"})
	}
	return ctx.JSON(payment)
}

func (c *PaymentController) GetPaymentsByTenant(ctx fiber.Ctx) error {
	tenantId := ctx.Params("tenantId")
	payments, err := c.paymentService.GetPaymentsByTenant(tenantId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(payments)
}

func (c *PaymentController) UpdatePayment(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var payment models.Payment
	if err := ctx.Bind().Body(&payment); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.paymentService.UpdatePayment(id, &payment); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.JSON(payment)
}

func (c *PaymentController) DeletePayment(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.paymentService.DeletePayment(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusNoContent)
}
