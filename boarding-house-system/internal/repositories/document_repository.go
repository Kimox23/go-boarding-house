package repositories

import (
	"database/sql"
	"time"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type DocumentRepository struct {
	db *sql.DB
}

func NewDocumentRepository(db *sql.DB) *DocumentRepository {
	return &DocumentRepository{db: db}
}

func (r *DocumentRepository) UploadDocument(document *models.Document) error {
	query := `INSERT INTO documents 
	          (tenant_id, document_type, file_path, verified_by, notes)
	          VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, document.TenantID, document.DocumentType,
		document.FilePath, document.VerifiedBy, document.Notes)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	document.ID = int(id)
	document.UploadDate = time.Now()
	document.Verified = false
	return nil
}

func (r *DocumentRepository) GetTenantDocuments(tenantId int) ([]models.Document, error) {
	query := `SELECT document_id, tenant_id, document_type, file_path, upload_date,
	          verified, verified_by, notes
	          FROM documents WHERE tenant_id = ?`

	rows, err := r.db.Query(query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var documents []models.Document
	for rows.Next() {
		var document models.Document
		err := rows.Scan(&document.ID, &document.TenantID, &document.DocumentType,
			&document.FilePath, &document.UploadDate, &document.Verified,
			&document.VerifiedBy, &document.Notes)
		if err != nil {
			return nil, err
		}
		documents = append(documents, document)
	}

	return documents, nil
}

func (r *DocumentRepository) VerifyDocument(id int, verified bool, notes string, verifiedBy int) error {
	query := `UPDATE documents SET 
	          verified = ?, verified_by = ?, notes = ?
	          WHERE document_id = ?`

	_, err := r.db.Exec(query, verified, verifiedBy, notes, id)
	return err
}

func (r *DocumentRepository) DeleteDocument(id int) error {
	query := `DELETE FROM documents WHERE document_id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
