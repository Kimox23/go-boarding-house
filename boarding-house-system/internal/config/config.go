package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	Port   string
	AppEnv string
	AppURL string

	JWTSecret     string
	JWTExpiration time.Duration

	UploadDir     string
	MaxUploadSize int64
	AllowedTypes  []string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	EmailFrom    string

	AdminEmail    string
	AdminPassword string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found - using default values")
	}

	return &Config{
		// Database Configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "root"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", ""),

		// Server Configuration
		Port:   getEnv("PORT", "8080"),
		AppEnv: getEnv("APP_ENV", "development"),
		AppURL: getEnv("APP_URL", "http://localhost:8080"),

		// Authentication
		JWTSecret:     getEnv("JWT_SECRET", "your_jwt_secret_key"),
		JWTExpiration: parseDuration(getEnv("JWT_EXPIRATION", "24h")),

		// File Uploads
		UploadDir:     getEnv("UPLOAD_DIR", "uploads"),
		MaxUploadSize: parseInt64(getEnv("MAX_UPLOAD_SIZE", "5242880")), // 5MB
		AllowedTypes:  parseAllowedTypes(getEnv("ALLOWED_FILE_TYPES", ".pdf,.jpg,.jpeg,.png")),

		// Email Configuration
		SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
		SMTPPort:     parseInt(getEnv("SMTP_PORT", "587")),
		SMTPUser:     getEnv("SMTP_USER", "your_email@example.com"),
		SMTPPassword: getEnv("SMTP_PASSWORD", "your_email_password"),
		EmailFrom:    getEnv("EMAIL_FROM", "support@example.com"),

		// Admin Defaults
		AdminEmail:    getEnv("ADMIN_EMAIL", "admin@example.com"),
		AdminPassword: getEnv("ADMIN_INITIAL_PASSWORD", "ChangeMe123!"),
	}
}

// Helper functions remain the same as previous example
func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func parseInt(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}

func parseInt64(s string) int64 {
	i, _ := strconv.ParseInt(s, 10, 64)
	return i
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		return 24 * time.Hour // Default to 24h if parsing fails
	}
	return d
}

func parseAllowedTypes(s string) []string {
	return strings.Split(s, ",")
}
