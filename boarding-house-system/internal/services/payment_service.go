package services

import (
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
)

type PaymentService struct {
	paymentRepo *repositories.PaymentRepository
}

func NewPaymentService(paymentRepo *repositories.PaymentRepository) *PaymentService {
	return &PaymentService{paymentRepo: paymentRepo}
}

func (s *PaymentService) CreatePayment(payment *models.Payment) error {
	return s.paymentRepo.CreatePayment(payment)
}

func (s *PaymentService) GetPayment(id string) (*models.Payment, error) {
	paymentID, err := strconv.Atoi(id)
	if err != nil {
		return nil, err
	}
	return s.paymentRepo.GetPayment(paymentID)
}

func (s *PaymentService) GetPaymentsByTenant(tenantId string) ([]models.Payment, error) {
	tenantID, err := strconv.Atoi(tenantId)
	if err != nil {
		return nil, err
	}
	return s.paymentRepo.GetPaymentsByTenant(tenantID)
}

func (s *PaymentService) UpdatePayment(id string, payment *models.Payment) error {
	paymentID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.paymentRepo.UpdatePayment(paymentID, payment)
}

func (s *PaymentService) DeletePayment(id string) error {
	paymentID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.paymentRepo.DeletePayment(paymentID)
}
