package services

import (
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
)

type TenantService struct {
	tenantRepo *repositories.TenantRepository
}

func NewTenantService(tenantRepo *repositories.TenantRepository) *TenantService {
	return &TenantService{tenantRepo: tenantRepo}
}

func (s *TenantService) CreateTenant(tenant *models.Tenant) error {
	return s.tenantRepo.CreateTenant(tenant)
}

func (s *TenantService) GetTenant(id string) (*models.Tenant, error) {
	tenantID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return s.tenantRepo.GetTenant(tenantID)
}

func (s *TenantService) GetTenantsByHouse(houseId string) ([]models.Tenant, error) {
	houseID, err := strconv.Atoi(houseId)
	if err != nil {
		return nil, err
	}
	return s.tenantRepo.GetTenantsByHouse(houseID)
}

func (s *TenantService) UpdateTenant(id string, tenant *models.Tenant) error {
	tenantID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.tenantRepo.UpdateTenant(tenantID, tenant)
}

func (s *TenantService) DeleteTenant(id string) error {
	tenantID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.tenantRepo.DeleteTenant(tenantID)
}
