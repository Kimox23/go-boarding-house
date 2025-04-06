package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/Kimox23/boarding-house-app/internal/config"
	"github.com/Kimox23/boarding-house-app/internal/routes"
	"github.com/Kimox23/boarding-house-app/migrations"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/gofiber/fiber/v3/middleware/recover"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Initialize database
	if err := config.InitDB(cfg); err != nil {
		log.Fatalf("Database initialization failed: %v", err)
	}
	defer config.CloseDB()
	// Run migrations
	log.Println("Running database migrations...")
	if err := migrations.RunMigrations(config.DB); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
	log.Println("Migrations completed successfully")

	// Create Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: func(ctx fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return ctx.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})

	app.Get("/debug/db", func(c fiber.Ctx) error {
		// Check database connection
		if err := config.DB.Ping(); err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Database connection failed",
				"error":   err.Error(),
			})
		}

		// List all tables
		rows, err := config.DB.Query("SHOW TABLES")
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to query tables",
				"error":   err.Error(),
			})
		}
		defer rows.Close()

		var tables []string
		for rows.Next() {
			var table string
			if err := rows.Scan(&table); err != nil {
				return c.Status(500).JSON(fiber.Map{
					"status":  "error",
					"message": "Failed to scan table name",
					"error":   err.Error(),
				})
			}
			tables = append(tables, table)
		}

		return c.JSON(fiber.Map{
			"status":     "success",
			"tables":     tables,
			"db_version": getDBVersion(config.DB),
		})
	})

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))
	app.Use(logger.New())
	app.Use(recover.New())

	// Setup routes
	routes.SetupRoutes(app, config.DB, cfg)

	// Start server
	port := cfg.Port
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	log.Printf("Server running on port %s", port)
	log.Fatal(app.Listen(":" + port))
}

func getDBVersion(db *sql.DB) string {
	var version string
	err := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return "unknown"
	}
	return version
}
