package repositories

import (
	"database/sql"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(payment *models.Payment) error {
	query := `INSERT INTO payments 
	          (tenant_id, amount, payment_date, payment_method, 
	           payment_for_month, receipt_number, status, notes, recorded_by)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, payment.TenantID, payment.Amount, payment.PaymentDate,
		payment.PaymentMethod, payment.PaymentForMonth, payment.ReceiptNumber,
		payment.Status, payment.Notes, payment.RecordedBy)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	payment.ID = int(id)
	return nil
}

func (r *PaymentRepository) GetPayment(id int) (*models.Payment, error) {
	query := `SELECT id, tenant_id, amount, payment_date, payment_method,
	          payment_for_month, receipt_number, status, notes, recorded_by
	          FROM payments WHERE id = ?`

	row := r.db.QueryRow(query, id)

	payment := &models.Payment{}
	err := row.Scan(&payment.ID, &payment.TenantID, &payment.Amount, &payment.PaymentDate,
		&payment.PaymentMethod, &payment.PaymentForMonth, &payment.ReceiptNumber,
		&payment.Status, &payment.Notes, &payment.RecordedBy)
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) GetPaymentsByTenant(tenantId int) ([]models.Payment, error) {
	query := `SELECT id, tenant_id, amount, payment_date, payment_method,
	          payment_for_month, receipt_number, status, notes, recorded_by
	          FROM payments WHERE tenant_id = ?`

	rows, err := r.db.Query(query, tenantId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.Payment
	for rows.Next() {
		var payment models.Payment
		err := rows.Scan(&payment.ID, &payment.TenantID, &payment.Amount, &payment.PaymentDate,
			&payment.PaymentMethod, &payment.PaymentForMonth, &payment.ReceiptNumber,
			&payment.Status, &payment.Notes, &payment.RecordedBy)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *PaymentRepository) UpdatePayment(id int, payment *models.Payment) error {
	query := `UPDATE payments SET 
	          tenant_id = ?, amount = ?, payment_date = ?, payment_method = ?,
	          payment_for_month = ?, receipt_number = ?, status = ?, notes = ?, recorded_by = ?
	          WHERE id = ?`

	_, err := r.db.Exec(query, payment.TenantID, payment.Amount, payment.PaymentDate,
		payment.PaymentMethod, payment.PaymentForMonth, payment.ReceiptNumber,
		payment.Status, payment.Notes, payment.RecordedBy, id)
	return err
}

func (r *PaymentRepository) DeletePayment(id int) error {
	query := `DELETE FROM payments WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
