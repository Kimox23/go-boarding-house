package services

import (
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
)

type HouseService struct {
	houseRepo *repositories.HouseRepository
}

func NewHouseService(houseRepo *repositories.HouseRepository) *HouseService {
	return &HouseService{houseRepo: houseRepo}
}

func (s *HouseService) CreateHouse(house *models.BoardingHouse) error {
	return s.houseRepo.CreateHouse(house)
}

func (s *HouseService) GetHouse(id string) (*models.BoardingHouse, error) {
	houseID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return s.houseRepo.GetHouse(houseID)
}

func (s *HouseService) GetAllHouses() ([]models.BoardingHouse, error) {
	return s.houseRepo.GetAllHouses()
}

func (s *HouseService) UpdateHouse(id string, house *models.BoardingHouse) error {
	houseID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.houseRepo.UpdateHouse(houseID, house)
}

func (s *HouseService) DeleteHouse(id string) error {
	houseID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.houseRepo.DeleteHouse(houseID)
}
