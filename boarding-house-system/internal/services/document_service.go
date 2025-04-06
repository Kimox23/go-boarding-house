package services

import (
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
)

type DocumentService struct {
	documentRepo *repositories.DocumentRepository
}

func NewDocumentService(documentRepo *repositories.DocumentRepository) *DocumentService {
	return &DocumentService{documentRepo: documentRepo}
}

func (s *DocumentService) UploadDocument(document *models.Document) error {
	return s.documentRepo.UploadDocument(document)
}

func (s *DocumentService) GetTenantDocuments(tenantId string) ([]models.Document, error) {
	tenantID, err := strconv.Atoi(tenantId)
	if err != nil {
		return nil, err
	}
	return s.documentRepo.GetTenantDocuments(tenantID)
}

func (s *DocumentService) VerifyDocument(id string, verified bool, notes string, verifiedBy int) error {
	documentID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.documentRepo.VerifyDocument(documentID, verified, notes, verifiedBy)
}

func (s *DocumentService) DeleteDocument(id string) error {
	documentID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}
	return s.documentRepo.DeleteDocument(documentID)
}
