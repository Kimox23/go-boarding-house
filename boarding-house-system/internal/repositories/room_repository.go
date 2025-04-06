package repositories

import (
	"database/sql"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) CreateRoom(room *models.Room) error {
	query := `INSERT INTO rooms 
		(house_id, room_number, room_type, capacity, price_per_month, status, description) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query,
		room.HouseID, room.RoomNumber, room.RoomType,
		room.Capacity, room.PricePerMonth, room.Status, room.Description)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	room.ID = int(id)
	return nil
}

func (r *RoomRepository) GetRoom(id int) (*models.Room, error) {
	query := `SELECT id, house_id, room_number, room_type, capacity, 
		current_occupancy, price_per_month, status, description 
		FROM rooms WHERE id = ?`

	row := r.db.QueryRow(query, id)

	room := &models.Room{}
	err := row.Scan(&room.ID, &room.HouseID, &room.RoomNumber, &room.RoomType,
		&room.Capacity, &room.CurrentOccupancy, &room.PricePerMonth,
		&room.Status, &room.Description)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) GetAllRooms() ([]models.Room, error) {
	query := `SELECT id, house_id, room_number, room_type, capacity, 
		current_occupancy, price_per_month, status, description 
		FROM rooms`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []models.Room
	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.HouseID, &room.RoomNumber, &room.RoomType,
			&room.Capacity, &room.CurrentOccupancy, &room.PricePerMonth,
			&room.Status, &room.Description)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	return rooms, nil
}

func (r *RoomRepository) UpdateRoom(id int, room *models.Room) error {
	query := `UPDATE rooms SET 
		house_id = ?, room_number = ?, room_type = ?, capacity = ?, 
		price_per_month = ?, status = ?, description = ? 
		WHERE id = ?`

	_, err := r.db.Exec(query,
		room.HouseID, room.RoomNumber, room.RoomType,
		room.Capacity, room.PricePerMonth, room.Status,
		room.Description, id)

	return err
}

func (r *RoomRepository) DeleteRoom(id int) error {
	query := `DELETE FROM rooms WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
