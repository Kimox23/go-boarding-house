package services

import (
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
	"github.com/Kimox23/boarding-house-app/internal/utils"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) CreateUser(user *models.User) error {
	return s.userRepo.CreateUser(user)
}

func (s *UserService) GetUser(id string) (*models.User, error) {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetUser(userID)
}

func (s *UserService) GetUserByEmail(email string) (*models.User, error) {
	return s.userRepo.GetUserByEmail(email)
}

func (s *UserService) GetAllUsers(pagination utils.Pagination) ([]models.User, int, error) {
	users, err := s.userRepo.GetAllUsers(pagination.Page, pagination.PageSize)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.userRepo.CountUsers()
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

func (s *UserService) UpdateUser(id string, user *models.User) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.userRepo.UpdateUser(userID, user)
}

func (s *UserService) DeleteUser(id string) error {
	userID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.userRepo.DeleteUser(userID)
}

func (s *UserService) CreateProfile(userId string, profile *models.UserProfile) error {
	userID, err := strconv.Atoi(userId)
	if err != nil {
		return err
	}
	profile.UserID = userID
	return s.userRepo.CreateProfile(profile)
}

func (s *UserService) GetProfile(userId string) (*models.UserProfile, error) {
	userID, err := strconv.Atoi(userId)
	if err != nil {
		return nil, err
	}
	return s.userRepo.GetProfile(userID)
}
