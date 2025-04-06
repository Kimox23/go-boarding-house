package models

import "time"

type User struct {
	ID        int       `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Phone     string    `json:"phone"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsActive  bool      `json:"is_active"`
}

type UserProfile struct {
	ID                    int    `json:"id"`
	UserID                int    `json:"user_id"`
	FirstName             string `json:"first_name"`
	LastName              string `json:"last_name"`
	DateOfBirth           string `json:"date_of_birth"`
	Gender                string `json:"gender"`
	Address               string `json:"address"`
	IDNumber              string `json:"id_number"`
	IDType                string `json:"id_type"`
	EmergencyContactName  string `json:"emergency_contact_name"`
	EmergencyContactPhone string `json:"emergency_contact_phone"`
	ProfilePicture        string `json:"profile_picture"`
}

type BoardingHouse struct {
	ID             int       `json:"id"`
	Name           string    `json:"name"`
	Address        string    `json:"address"`
	Description    string    `json:"description"`
	TotalRooms     int       `json:"total_rooms"`
	AvailableRooms int       `json:"available_rooms"`
	ManagerID      int       `json:"manager_id"`
	Amenities      string    `json:"amenities"`
	Rules          string    `json:"rules"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

type Room struct {
	ID               int     `json:"id"`
	HouseID          int     `json:"house_id"`
	RoomNumber       string  `json:"room_number"`
	RoomType         string  `json:"room_type"`
	Capacity         int     `json:"capacity"`
	CurrentOccupancy int     `json:"current_occupancy"`
	PricePerMonth    float64 `json:"price_per_month"`
	Status           string  `json:"status"`
	Description      string  `json:"description"`
}

type Tenant struct {
	ID               int        `json:"id"`
	UserID           int        `json:"user_id"`
	RoomID           int        `json:"room_id"`
	MoveInDate       time.Time  `json:"move_in_date"`
	MoveOutDate      *time.Time `json:"move_out_date"`
	DepositAmount    float64    `json:"deposit_amount"`
	DepositPaid      bool       `json:"deposit_paid"`
	ContractDocument string     `json:"contract_document"`
	Status           string     `json:"status"`
}

type Payment struct {
	ID              int       `json:"id"`
	TenantID        int       `json:"tenant_id"`
	Amount          float64   `json:"amount"`
	PaymentDate     time.Time `json:"payment_date"`
	PaymentMethod   string    `json:"payment_method"`
	PaymentForMonth time.Time `json:"payment_for_month"`
	ReceiptNumber   string    `json:"receipt_number"`
	Status          string    `json:"status"`
	Notes           string    `json:"notes"`
	RecordedBy      int       `json:"recorded_by"`
}

type MaintenanceRequest struct {
	ID            int        `json:"id"`
	RoomID        int        `json:"room_id"`
	ReportedBy    int        `json:"reported_by"`
	IssueType     string     `json:"issue_type"`
	Description   string     `json:"description"`
	Priority      string     `json:"priority"`
	Status        string     `json:"status"`
	ReportedDate  time.Time  `json:"reported_date"`
	CompletedDate *time.Time `json:"completed_date"`
	AssignedTo    *int       `json:"assigned_to"`
	Cost          *float64   `json:"cost"`
}

type Notification struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	IsRead    bool      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
	Link      string    `json:"link"`
}

type Document struct {
	ID           int       `json:"id"`
	TenantID     int       `json:"tenant_id"`
	DocumentType string    `json:"document_type"`
	FilePath     string    `json:"file_path"`
	UploadDate   time.Time `json:"upload_date"`
	Verified     bool      `json:"verified"`
	VerifiedBy   *int      `json:"verified_by"`
	Notes        string    `json:"notes"`
}
