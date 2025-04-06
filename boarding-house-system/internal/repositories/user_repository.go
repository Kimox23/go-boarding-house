package repositories

import (
	"database/sql"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	query := `INSERT INTO users (username, email, password_hash, role, phone) 
	          VALUES (?, ?, ?, ?, ?)`

	result, err := r.db.Exec(query, user.Username, user.Email, user.Password, user.Role, user.Phone)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	user.ID = int(id)
	return nil
}

func (r *UserRepository) GetUser(id int) (*models.User, error) {
	query := `SELECT user_id, username, email, phone, role, created_at, updated_at, is_active 
	          FROM users WHERE user_id = ?`

	row := r.db.QueryRow(query, id)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Phone,
		&user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	query := `SELECT user_id, username, email, password_hash, phone, role, 
	          created_at, updated_at, is_active 
	          FROM users WHERE email = ?`

	row := r.db.QueryRow(query, email)

	user := &models.User{}
	err := row.Scan(&user.ID, &user.Username, &user.Email, &user.Password,
		&user.Phone, &user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil if user not found
		}
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) GetAllUsers(page, pageSize int) ([]models.User, error) {
	offset := (page - 1) * pageSize
	query := `SELECT user_id, username, email, phone, role, created_at, updated_at, is_active 
	          FROM users LIMIT ? OFFSET ?`

	rows, err := r.db.Query(query, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Phone,
			&user.Role, &user.CreatedAt, &user.UpdatedAt, &user.IsActive)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func (r *UserRepository) CountUsers() (int, error) {
	query := `SELECT COUNT(*) FROM users`
	var count int
	err := r.db.QueryRow(query).Scan(&count)
	return count, err
}

func (r *UserRepository) UpdateUser(id int, user *models.User) error {
	query := `UPDATE users SET 
	          username = ?, email = ?, phone = ?, role = ?, is_active = ?
	          WHERE user_id = ?`

	_, err := r.db.Exec(query, user.Username, user.Email, user.Phone,
		user.Role, user.IsActive, id)
	return err
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `DELETE FROM users WHERE user_id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *UserRepository) CreateProfile(profile *models.UserProfile) error {
	query := `INSERT INTO user_profiles 
	          (user_id, first_name, last_name, date_of_birth, gender, address, 
	           id_number, id_type, emergency_contact_name, emergency_contact_phone, profile_picture)
	          VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, profile.UserID, profile.FirstName, profile.LastName,
		profile.DateOfBirth, profile.Gender, profile.Address, profile.IDNumber,
		profile.IDType, profile.EmergencyContactName, profile.EmergencyContactPhone,
		profile.ProfilePicture)
	return err
}

func (r *UserRepository) GetProfile(userId int) (*models.UserProfile, error) {
	query := `SELECT profile_id, user_id, first_name, last_name, date_of_birth, gender, address,
	          id_number, id_type, emergency_contact_name, emergency_contact_phone, profile_picture
	          FROM user_profiles WHERE user_id = ?`

	row := r.db.QueryRow(query, userId)

	profile := &models.UserProfile{}
	err := row.Scan(&profile.ID, &profile.UserID, &profile.FirstName, &profile.LastName,
		&profile.DateOfBirth, &profile.Gender, &profile.Address, &profile.IDNumber,
		&profile.IDType, &profile.EmergencyContactName, &profile.EmergencyContactPhone,
		&profile.ProfilePicture)
	if err != nil {
		return nil, err
	}

	return profile, nil
}
