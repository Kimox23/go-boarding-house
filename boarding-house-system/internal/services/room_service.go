package services

import (
	"errors"
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
)

type RoomService struct {
	roomRepo *repositories.RoomRepository
}

func NewRoomService(roomRepo *repositories.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) CreateRoom(room *models.Room) error {
	return s.roomRepo.CreateRoom(room)
}

func (s *RoomService) GetRoom(id string) (*models.Room, error) {
	roomID, err := strconv.Atoi(id)
	if err != nil {
		return nil, errors.New("invalid room ID")
	}
	return s.roomRepo.GetRoom(roomID)
}

func (s *RoomService) GetAllRooms() ([]models.Room, error) {
	return s.roomRepo.GetAllRooms()
}

func (s *RoomService) UpdateRoom(id string, room *models.Room) error {
	roomID, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid room ID")
	}
	return s.roomRepo.UpdateRoom(roomID, room)
}

func (s *RoomService) DeleteRoom(id string) error {
	roomID, err := strconv.Atoi(id)
	if err != nil {
		return errors.New("invalid room ID")
	}
	return s.roomRepo.DeleteRoom(roomID)
}
