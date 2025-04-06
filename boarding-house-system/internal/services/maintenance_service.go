package services

import (
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
)

type MaintenanceService struct {
	maintenanceRepo *repositories.MaintenanceRepository
}

func NewMaintenanceService(maintenanceRepo *repositories.MaintenanceRepository) *MaintenanceService {
	return &MaintenanceService{maintenanceRepo: maintenanceRepo}
}

func (s *MaintenanceService) CreateRequest(request *models.MaintenanceRequest) error {
	return s.maintenanceRepo.CreateRequest(request)
}

func (s *MaintenanceService) GetRequest(id string) (*models.MaintenanceRequest, error) {
	requestID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return s.maintenanceRepo.GetRequest(requestID)
}

func (s *MaintenanceService) GetRequestsByRoom(roomId string) ([]models.MaintenanceRequest, error) {
	roomID, err := strconv.Atoi(roomId)
	if err != nil {
		return nil, err
	}
	return s.maintenanceRepo.GetRequestsByRoom(roomID)
}

func (s *MaintenanceService) UpdateRequest(id string, request *models.MaintenanceRequest) error {
	requestID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.maintenanceRepo.UpdateRequest(requestID, request)
}

func (s *MaintenanceService) UpdateRequestStatus(id string, status string, assignedTo *int) error {
	requestID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.maintenanceRepo.UpdateRequestStatus(requestID, status, assignedTo)
}

func (s *MaintenanceService) DeleteRequest(id string) error {
	requestID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.maintenanceRepo.DeleteRequest(requestID)
}
