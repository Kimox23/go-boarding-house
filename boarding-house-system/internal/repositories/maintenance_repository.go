package repositories

import (
	"database/sql"
	"time"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type MaintenanceRepository struct {
	db *sql.DB
}

func NewMaintenanceRepository(db *sql.DB) *MaintenanceRepository {
	return &MaintenanceRepository{db: db}
}

func (r *MaintenanceRepository) CreateRequest(request *models.MaintenanceRequest) error {
	query := `INSERT INTO maintenance_requests 
	          (room_id, reported_by, issue_type, description, priority, assigned_to)
	          VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, request.RoomID, request.ReportedBy, request.IssueType,
		request.Description, request.Priority, request.AssignedTo)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	request.ID = int(id)
	request.ReportedDate = time.Now()
	request.Status = "pending"
	return nil
}

func (r *MaintenanceRepository) GetRequest(id int) (*models.MaintenanceRequest, error) {
	query := `SELECT id, room_id, reported_by, issue_type, description, priority,
	          status, reported_date, completed_date, assigned_to, cost
	          FROM maintenance_requests WHERE id = ?`

	row := r.db.QueryRow(query, id)

	request := &models.MaintenanceRequest{}
	var completedDate sql.NullTime
	err := row.Scan(&request.ID, &request.RoomID, &request.ReportedBy, &request.IssueType,
		&request.Description, &request.Priority, &request.Status, &request.ReportedDate,
		&completedDate, &request.AssignedTo, &request.Cost)
	if err != nil {
		return nil, err
	}

	if completedDate.Valid {
		request.CompletedDate = &completedDate.Time
	}
	return request, nil
}

func (r *MaintenanceRepository) GetRequestsByRoom(roomId int) ([]models.MaintenanceRequest, error) {
	query := `SELECT id, room_id, reported_by, issue_type, description, priority,
	          status, reported_date, completed_date, assigned_to, cost
	          FROM maintenance_requests WHERE room_id = ?`

	rows, err := r.db.Query(query, roomId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var requests []models.MaintenanceRequest
	for rows.Next() {
		var request models.MaintenanceRequest
		var completedDate sql.NullTime
		err := rows.Scan(&request.ID, &request.RoomID, &request.ReportedBy, &request.IssueType,
			&request.Description, &request.Priority, &request.Status, &request.ReportedDate,
			&completedDate, &request.AssignedTo, &request.Cost)
		if err != nil {
			return nil, err
		}
		if completedDate.Valid {
			request.CompletedDate = &completedDate.Time
		}
		requests = append(requests, request)
	}

	return requests, nil
}

func (r *MaintenanceRepository) UpdateRequest(id int, request *models.MaintenanceRequest) error {
	query := `UPDATE maintenance_requests SET 
	          room_id = ?, issue_type = ?, description = ?, priority = ?, cost = ?
	          WHERE id = ?`

	_, err := r.db.Exec(query, request.RoomID, request.IssueType, request.Description,
		request.Priority, request.Cost, id)
	return err
}

func (r *MaintenanceRepository) UpdateRequestStatus(id int, status string, assignedTo *int) error {
	query := `UPDATE maintenance_requests SET 
	          status = ?, assigned_to = ?, completed_date = ?
	          WHERE id = ?`

	var completedDate interface{}
	if status == "completed" {
		completedDate = time.Now()
	} else {
		completedDate = nil
	}

	_, err := r.db.Exec(query, status, assignedTo, completedDate, id)
	return err
}

func (r *MaintenanceRepository) DeleteRequest(id int) error {
	query := `DELETE FROM maintenance_requests WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
