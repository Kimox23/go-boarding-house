package repositories

import (
	"database/sql"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type HouseRepository struct {
	db *sql.DB
}

func NewHouseRepository(db *sql.DB) *HouseRepository {
	return &HouseRepository{db: db}
}

func (r *HouseRepository) CreateHouse(house *models.BoardingHouse) error {
	query := `INSERT INTO boarding_houses 
	          (name, address, description, total_rooms, available_rooms, 
	           manager_id, amenities, rules)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, house.Name, house.Address, house.Description,
		house.TotalRooms, house.AvailableRooms, house.ManagerID, house.Amenities,
		house.Rules)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	house.ID = int(id)
	return nil
}

func (r *HouseRepository) GetHouse(id int) (*models.BoardingHouse, error) {
	query := `SELECT id, name, address, description, total_rooms, available_rooms,
	          manager_id, amenities, rules, created_at, updated_at
	          FROM boarding_houses WHERE id = ?`

	row := r.db.QueryRow(query, id)

	house := &models.BoardingHouse{}
	err := row.Scan(&house.ID, &house.Name, &house.Address, &house.Description,
		&house.TotalRooms, &house.AvailableRooms, &house.ManagerID, &house.Amenities,
		&house.Rules, &house.CreatedAt, &house.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return house, nil
}

func (r *HouseRepository) GetAllHouses() ([]models.BoardingHouse, error) {
	query := `SELECT id, name, address, description, total_rooms, available_rooms,
	          manager_id, amenities, rules, created_at, updated_at
	          FROM boarding_houses`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var houses []models.BoardingHouse
	for rows.Next() {
		var house models.BoardingHouse
		err := rows.Scan(&house.ID, &house.Name, &house.Address, &house.Description,
			&house.TotalRooms, &house.AvailableRooms, &house.ManagerID, &house.Amenities,
			&house.Rules, &house.CreatedAt, &house.UpdatedAt)
		if err != nil {
			return nil, err
		}
		houses = append(houses, house)
	}

	return houses, nil
}

func (r *HouseRepository) UpdateHouse(id int, house *models.BoardingHouse) error {
	query := `UPDATE boarding_houses SET 
	          name = ?, address = ?, description = ?, total_rooms = ?,
	          available_rooms = ?, manager_id = ?, amenities = ?, rules = ?
	          WHERE id = ?`

	_, err := r.db.Exec(query, house.Name, house.Address, house.Description,
		house.TotalRooms, house.AvailableRooms, house.ManagerID, house.Amenities,
		house.Rules, id)
	return err
}

func (r *HouseRepository) DeleteHouse(id int) error {
	query := `DELETE FROM boarding_houses WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
