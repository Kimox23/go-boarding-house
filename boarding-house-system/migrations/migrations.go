package migrations

import (
	"database/sql"
	"fmt"
	"log"
)

func RunMigrations(db *sql.DB) error {
	log.Println("Starting database migrations...")

	// Verify connection first
	if err := db.Ping(); err != nil {
		return fmt.Errorf("database connection failed: %w", err)
	}

	// Check if users table already exists
	var tableExists int
	err := db.QueryRow(`
		 SELECT COUNT(*) 
		 FROM information_schema.tables 
		 WHERE table_schema = DATABASE() 
		 AND table_name = 'users'
	 `).Scan(&tableExists)

	if err != nil {
		return fmt.Errorf("failed to check for existing tables: %w", err)
	}

	if tableExists > 0 {
		log.Println("Database already initialized, skipping migrations")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("could not begin transaction: %w", err)
	}

	defer func() {
		if err != nil {
			log.Printf("Migration failed, rolling back: %v", err)
			if rbErr := tx.Rollback(); rbErr != nil {
				log.Printf("Rollback failed: %v", rbErr)
			}
		}
	}()

	tables := []struct {
		name  string
		query string
	}{
		{
			"users",
			`CREATE TABLE IF NOT EXISTS users (
				user_id INT PRIMARY KEY AUTO_INCREMENT,
				username VARCHAR(50) UNIQUE NOT NULL,
				password_hash VARCHAR(255) NOT NULL,
				email VARCHAR(100) UNIQUE NOT NULL,
				phone VARCHAR(20),
				role ENUM('admin', 'manager', 'staff', 'tenant') NOT NULL,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				is_active BOOLEAN DEFAULT TRUE
			)`,
		},
		{
			"user_profiles",
			`CREATE TABLE IF NOT EXISTS user_profiles (
				profile_id INT PRIMARY KEY AUTO_INCREMENT,
				user_id INT UNIQUE NOT NULL,
				first_name VARCHAR(50) NOT NULL,
				last_name VARCHAR(50) NOT NULL,
				date_of_birth DATE,
				gender ENUM('male', 'female', 'other'),
				address TEXT,
				id_number VARCHAR(50),
				id_type VARCHAR(50),
				emergency_contact_name VARCHAR(100),
				emergency_contact_phone VARCHAR(20),
				profile_picture VARCHAR(255),
				FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
			)`,
		},
		{
			"boarding_houses",
			`CREATE TABLE IF NOT EXISTS boarding_houses (
				house_id INT PRIMARY KEY AUTO_INCREMENT,
				name VARCHAR(100) NOT NULL,
				address TEXT NOT NULL,
				description TEXT,
				total_rooms INT NOT NULL,
				available_rooms INT NOT NULL,
				manager_id INT,
				amenities TEXT,
				rules TEXT,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
				FOREIGN KEY (manager_id) REFERENCES users(user_id)
			)`,
		},
		{
			"rooms",
			`CREATE TABLE IF NOT EXISTS rooms (
				room_id INT PRIMARY KEY AUTO_INCREMENT,
				house_id INT NOT NULL,
				room_number VARCHAR(20) NOT NULL,
				room_type ENUM('single', 'double', 'dormitory', 'suite') NOT NULL,
				capacity INT NOT NULL,
				current_occupancy INT DEFAULT 0,
				price_per_month DECIMAL(10,2) NOT NULL,
				status ENUM('available', 'occupied', 'maintenance') DEFAULT 'available',
				description TEXT,
				FOREIGN KEY (house_id) REFERENCES boarding_houses(house_id) ON DELETE CASCADE
			)`,
		},
		{
			"tenants",
			`CREATE TABLE IF NOT EXISTS tenants (
				tenant_id INT PRIMARY KEY AUTO_INCREMENT,
				user_id INT UNIQUE NOT NULL,
				room_id INT,
				move_in_date DATE NOT NULL,
				move_out_date DATE,
				deposit_amount DECIMAL(10,2),
				deposit_paid BOOLEAN DEFAULT FALSE,
				contract_document VARCHAR(255),
				status ENUM('active', 'inactive', 'pending') DEFAULT 'pending',
				FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE,
				FOREIGN KEY (room_id) REFERENCES rooms(room_id)
			)`,
		},
		{
			"payments",
			`CREATE TABLE IF NOT EXISTS payments (
				payment_id INT PRIMARY KEY AUTO_INCREMENT,
				tenant_id INT NOT NULL,
				amount DECIMAL(10,2) NOT NULL,
				payment_date DATE NOT NULL,
				payment_method ENUM('cash', 'bank_transfer', 'credit_card', 'mobile_payment') NOT NULL,
				payment_for_month DATE NOT NULL,
				receipt_number VARCHAR(50) UNIQUE,
				status ENUM('paid', 'pending', 'overdue', 'partial') NOT NULL,
				notes TEXT,
				recorded_by INT,
				FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
				FOREIGN KEY (recorded_by) REFERENCES users(user_id)
			)`,
		},
		{
			"maintenance_requests",
			`CREATE TABLE IF NOT EXISTS maintenance_requests (
				request_id INT PRIMARY KEY AUTO_INCREMENT,
				room_id INT NOT NULL,
				reported_by INT NOT NULL,
				issue_type VARCHAR(100) NOT NULL,
				description TEXT NOT NULL,
				priority ENUM('low', 'medium', 'high', 'emergency') NOT NULL,
				status ENUM('pending', 'in_progress', 'completed', 'cancelled') DEFAULT 'pending',
				reported_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				completed_date TIMESTAMP,
				assigned_to INT,
				cost DECIMAL(10,2),
				FOREIGN KEY (room_id) REFERENCES rooms(room_id),
				FOREIGN KEY (reported_by) REFERENCES users(user_id),
				FOREIGN KEY (assigned_to) REFERENCES users(user_id)
			)`,
		},
		{
			"notifications",
			`CREATE TABLE IF NOT EXISTS notifications (
				notification_id INT PRIMARY KEY AUTO_INCREMENT,
				user_id INT NOT NULL,
				title VARCHAR(100) NOT NULL,
				message TEXT NOT NULL,
				is_read BOOLEAN DEFAULT FALSE,
				created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				link VARCHAR(255),
				FOREIGN KEY (user_id) REFERENCES users(user_id) ON DELETE CASCADE
			)`,
		},
		{
			"documents",
			`CREATE TABLE IF NOT EXISTS documents (
				document_id INT PRIMARY KEY AUTO_INCREMENT,
				tenant_id INT NOT NULL,
				document_type VARCHAR(100) NOT NULL,
				file_path VARCHAR(255) NOT NULL,
				upload_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
				verified BOOLEAN DEFAULT FALSE,
				verified_by INT,
				notes TEXT,
				FOREIGN KEY (tenant_id) REFERENCES tenants(tenant_id),
				FOREIGN KEY (verified_by) REFERENCES users(user_id)
			)`,
		},
		// Add other tables here in proper foreign key dependency order
	}

	for _, table := range tables {
		log.Printf("Creating table %s...", table.name)
		if _, err := tx.Exec(table.query); err != nil {
			return fmt.Errorf("failed to create %s: %w", table.name, err)
		}
	}

	// Create admin user
	log.Println("Creating admin user...")
	if _, err := tx.Exec(`
		INSERT IGNORE INTO users 
		(username, email, password_hash, role, phone, is_active)
		VALUES (?, ?, ?, ?, ?, ?)`,
		"admin",
		"admin@example.com",
		"$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi",
		"admin",
		"1234567890",
		true,
	); err != nil {
		return fmt.Errorf("failed to create admin user: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	log.Println("Migrations completed successfully")
	return nil
}
