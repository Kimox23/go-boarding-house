package repositories

import (
	"database/sql"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type TenantRepository struct {
	db *sql.DB
}

func NewTenantRepository(db *sql.DB) *TenantRepository {
	return &TenantRepository{db: db}
}

func (r *TenantRepository) CreateTenant(tenant *models.Tenant) error {
	query := `INSERT INTO tenants 
	          (user_id, room_id, move_in_date, move_out_date, 
	           deposit_amount, deposit_paid, contract_document, status)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, tenant.UserID, tenant.RoomID, tenant.MoveInDate,
		tenant.MoveOutDate, tenant.DepositAmount, tenant.DepositPaid, tenant.ContractDocument,
		tenant.Status)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	tenant.ID = int(id)
	return nil
}

func (r *TenantRepository) GetTenant(id int) (*models.Tenant, error) {
	query := `SELECT tenant_id, user_id, room_id, move_in_date, move_out_date,
	          deposit_amount, deposit_paid, contract_document, status
	          FROM tenants WHERE id = ?`

	row := r.db.QueryRow(query, id)

	tenant := &models.Tenant{}
	err := row.Scan(&tenant.ID, &tenant.UserID, &tenant.RoomID, &tenant.MoveInDate,
		&tenant.MoveOutDate, &tenant.DepositAmount, &tenant.DepositPaid,
		&tenant.ContractDocument, &tenant.Status)
	if err != nil {
		return nil, err
	}

	return tenant, nil
}

func (r *TenantRepository) GetTenantsByHouse(houseId int) ([]models.Tenant, error) {
	query := `SELECT t.tenant_id, t.user_id, t.room_id, t.move_in_date, t.move_out_date,
	          t.deposit_amount, t.deposit_paid, t.contract_document, t.status
	          FROM tenants t
	          JOIN rooms r ON t.room_id = r.tenant_id
	          WHERE r.house_id = ?`

	rows, err := r.db.Query(query, houseId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tenants []models.Tenant
	for rows.Next() {
		var tenant models.Tenant
		err := rows.Scan(&tenant.ID, &tenant.UserID, &tenant.RoomID, &tenant.MoveInDate,
			&tenant.MoveOutDate, &tenant.DepositAmount, &tenant.DepositPaid,
			&tenant.ContractDocument, &tenant.Status)
		if err != nil {
			return nil, err
		}
		tenants = append(tenants, tenant)
	}

	return tenants, nil
}

func (r *TenantRepository) UpdateTenant(id int, tenant *models.Tenant) error {
	query := `UPDATE tenants SET 
	          room_id = ?, move_in_date = ?, move_out_date = ?,
	          deposit_amount = ?, deposit_paid = ?, contract_document = ?, status = ?
	          WHERE tenant_id = ?`

	_, err := r.db.Exec(query, tenant.RoomID, tenant.MoveInDate, tenant.MoveOutDate,
		tenant.DepositAmount, tenant.DepositPaid, tenant.ContractDocument,
		tenant.Status, id)
	return err
}

func (r *TenantRepository) DeleteTenant(id int) error {
	query := `DELETE FROM tenants WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
