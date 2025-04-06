package controllers

import (
	"net/http"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
)

type NotificationController struct {
	notificationService *services.NotificationService
}

func NewNotificationController(notificationService *services.NotificationService) *NotificationController {
	return &NotificationController{notificationService: notificationService}
}

func (c *NotificationController) CreateNotification(ctx fiber.Ctx) error {
	var notification models.Notification
	if err := ctx.Bind().Body(&notification); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.notificationService.CreateNotification(&notification); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.Status(http.StatusCreated).JSON(notification)
}

func (c *NotificationController) GetUserNotifications(ctx fiber.Ctx) error {
	userId := ctx.Params("userId")
	notifications, err := c.notificationService.GetUserNotifications(userId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(notifications)
}

func (c *NotificationController) MarkAsRead(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.notificationService.MarkAsRead(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusOK)
}

func (c *NotificationController) DeleteNotification(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.notificationService.DeleteNotification(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusNoContent)
}
