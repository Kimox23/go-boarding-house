package controllers

import (
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/Kimox23/boarding-house-app/internal/models"
	"github.com/Kimox23/boarding-house-app/internal/services"
	"github.com/Kimox23/boarding-house-app/internal/utils"

	"github.com/gofiber/fiber/v3"
)

type DocumentController struct {
	documentService *services.DocumentService
	uploadDir       string
}

func NewDocumentController(documentService *services.DocumentService, uploadDir string) *DocumentController {
	return &DocumentController{
		documentService: documentService,
		uploadDir:       uploadDir,
	}
}

func (c *DocumentController) UploadDocument(ctx fiber.Ctx) error {
	tenantID := ctx.Params("tenantId")
	file, err := ctx.FormFile("document")
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Document file is required"})
	}

	documentType := ctx.FormValue("document_type", "")
	if documentType == "" {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Document type is required"})
	}

	// Save the uploaded file
	filename, err := utils.SaveUploadedFile(ctx, file, c.uploadDir)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save document"})
	}

	tenantIDInt, err := strconv.Atoi(tenantID)
	if err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid tenant ID"})
	}

	document := &models.Document{
		TenantID:     tenantIDInt,
		DocumentType: documentType,
		FilePath:     filename,
	}

	if err := c.documentService.UploadDocument(document); err != nil {
		// Clean up the uploaded file if database operation fails
		os.Remove(filepath.Join(c.uploadDir, filename))
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save document info"})
	}

	return ctx.Status(http.StatusCreated).JSON(document)
}

func (c *DocumentController) GetTenantDocuments(ctx fiber.Ctx) error {
	tenantId := ctx.Params("tenantId")
	documents, err := c.documentService.GetTenantDocuments(tenantId)
	if err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(documents)
}

func (c *DocumentController) VerifyDocument(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	var verification struct {
		Verified   bool   `json:"verified"`
		Notes      string `json:"notes"`
		VerifiedBy int    `json:"verified_by"`
	}
	if err := ctx.Bind().Body(&verification); err != nil {
		return ctx.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
	}

	if err := c.documentService.VerifyDocument(id, verification.Verified,
		verification.Notes, verification.VerifiedBy); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return ctx.SendStatus(http.StatusOK)
}

func (c *DocumentController) DeleteDocument(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	if err := c.documentService.DeleteDocument(id); err != nil {
		return ctx.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.SendStatus(http.StatusNoContent)
}
