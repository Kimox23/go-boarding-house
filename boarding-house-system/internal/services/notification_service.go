package services

import (
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
)

type NotificationService struct {
	notificationRepo *repositories.NotificationRepository
}

func NewNotificationService(notificationRepo *repositories.NotificationRepository) *NotificationService {
	return &NotificationService{notificationRepo: notificationRepo}
}

func (s *NotificationService) CreateNotification(notification *models.Notification) error {
	return s.notificationRepo.CreateNotification(notification)
}

func (s *NotificationService) GetUserNotifications(userId string) ([]models.Notification, error) {
	userID, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}
	return s.notificationRepo.GetUserNotifications(userID)
}

func (s *NotificationService) MarkAsRead(id string) error {
	notificationID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *NotificationService) DeleteNotification(id string) error {
	notificationID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.notificationRepo.DeleteNotification(notificationID)
}
