package repositories

import (
	"database/sql"
	"time"

	"github.com/Kimox23/boarding-house-app/internal/models"
)

type NotificationRepository struct {
	db *sql.DB
}

func NewNotificationRepository(db *sql.DB) *NotificationRepository {
	return &NotificationRepository{db: db}
}

func (r *NotificationRepository) CreateNotification(notification *models.Notification) error {
	query := `INSERT INTO notifications 
	          (user_id, title, message, link)
	          VALUES (?, ?, ?, ?)`

	result, err := r.db.Exec(query, notification.UserID, notification.Title,
		notification.Message, notification.Link)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	notification.ID = int(id)
	notification.CreatedAt = time.Now()
	notification.IsRead = false
	return nil
}

func (r *NotificationRepository) GetUserNotifications(userId int) ([]models.Notification, error) {
	query := `SELECT id, user_id, title, message, is_read, created_at, link
	          FROM notifications WHERE user_id = ?
	          ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var notifications []models.Notification
	for rows.Next() {
		var notification models.Notification
		err := rows.Scan(&notification.ID, &notification.UserID, &notification.Title,
			&notification.Message, &notification.IsRead, &notification.CreatedAt,
			&notification.Link)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}

func (r *NotificationRepository) MarkAsRead(id int) error {
	query := `UPDATE notifications SET is_read = true WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}

func (r *NotificationRepository) DeleteNotification(id int) error {
	query := `DELETE FROM notifications WHERE id = ?`
	_, err := r.db.Exec(query, id)
	return err
}
