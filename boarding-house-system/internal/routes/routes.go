package routes

import (
	"database/sql"
	"log"
	"os"

	"github.com/Kimox23/boarding-house-app/internal/config"
	"github.com/Kimox23/boarding-house-app/internal/controllers"
	"github.com/Kimox23/boarding-house-app/internal/middleware"
	"github.com/Kimox23/boarding-house-app/internal/repositories"
	"github.com/Kimox23/boarding-house-app/internal/services"

	"github.com/gofiber/fiber/v3"
)

func SetupRoutes(app *fiber.App, db *sql.DB, cfg *config.Config) {
	log.Printf("Database pointer in SetupRoutes: %p", db)

	if db == nil {
		panic("Database connection is nil!")
	}
	// Ensure upload directory exists
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		panic(err)
	}

	// Initialize all repositories
	userRepo := repositories.NewUserRepository(db)
	houseRepo := repositories.NewHouseRepository(db)
	roomRepo := repositories.NewRoomRepository(db)
	tenantRepo := repositories.NewTenantRepository(db)
	paymentRepo := repositories.NewPaymentRepository(db)
	maintenanceRepo := repositories.NewMaintenanceRepository(db)
	notificationRepo := repositories.NewNotificationRepository(db)
	documentRepo := repositories.NewDocumentRepository(db)

	// Initialize all services
	userService := services.NewUserService(userRepo)
	houseService := services.NewHouseService(houseRepo)
	roomService := services.NewRoomService(roomRepo)
	tenantService := services.NewTenantService(tenantRepo)
	paymentService := services.NewPaymentService(paymentRepo)
	maintenanceService := services.NewMaintenanceService(maintenanceRepo)
	notificationService := services.NewNotificationService(notificationRepo)
	documentService := services.NewDocumentService(documentRepo)

	// Initialize all controllers
	authController := controllers.NewAuthController(userService, cfg)
	userController := controllers.NewUserController(userService)
	houseController := controllers.NewHouseController(houseService)
	roomController := controllers.NewRoomController(roomService)
	tenantController := controllers.NewTenantController(tenantService)
	paymentController := controllers.NewPaymentController(paymentService)
	maintenanceController := controllers.NewMaintenanceController(maintenanceService)
	notificationController := controllers.NewNotificationController(notificationService)
	documentController := controllers.NewDocumentController(documentService, uploadDir)

	app.Get("/uploads/*", func(c fiber.Ctx) error {
		file := "./uploads/" + c.Params("*")
		return c.SendFile(file)
	})

	// Auth routes
	authGroup := app.Group("/api/auth")
	{
		authGroup.Post("/register", authController.Register)
		authGroup.Post("/login", authController.Login)
		authGroup.Get("/me", authController.Me, middleware.AuthRequired(cfg))
	}

	// User routes
	userGroup := app.Group("/api/users", middleware.AuthRequired(cfg))
	{
		userGroup.Get("/", userController.GetAllUsers, middleware.RoleRequired("admin", cfg))
		userGroup.Get("/:id", userController.GetUser)
		userGroup.Put("/:id", userController.UpdateUser)
		userGroup.Delete("/:id", userController.DeleteUser, middleware.RoleRequired("admin", cfg))

		// User profile routes
		userGroup.Post("/:userId/profile", userController.CreateProfile)
		userGroup.Get("/:userId/profile", userController.GetProfile)
	}

	// Boarding house routes
	houseGroup := app.Group("/api/houses", middleware.AuthRequired(cfg))
	{
		houseGroup.Post("/", houseController.CreateHouse, middleware.RoleRequired("admin", cfg))
		houseGroup.Get("/", houseController.GetAllHouses)
		houseGroup.Get("/:id", houseController.GetHouse)
		houseGroup.Put("/:id", houseController.UpdateHouse, middleware.RoleRequired("admin", cfg))
		houseGroup.Delete("/:id", houseController.DeleteHouse, middleware.RoleRequired("admin", cfg))
	}

	// Room routes
	roomGroup := app.Group("/api/rooms", middleware.AuthRequired(cfg))
	{
		roomGroup.Post("/", roomController.CreateRoom, middleware.RoleRequired("manager", cfg))
		roomGroup.Get("/", roomController.GetAllRooms)
		roomGroup.Get("/:id", roomController.GetRoom)
		roomGroup.Put("/:id", roomController.UpdateRoom, middleware.RoleRequired("manager", cfg))
		roomGroup.Delete("/:id", roomController.DeleteRoom, middleware.RoleRequired("manager", cfg))
	}

	// Tenant routes
	tenantGroup := app.Group("/api/tenants", middleware.AuthRequired(cfg))
	{
		tenantGroup.Post("/", tenantController.CreateTenant, middleware.RoleRequired("manager", cfg))
		tenantGroup.Get("/house/:houseId", tenantController.GetTenantsByHouse)
		tenantGroup.Get("/:id", tenantController.GetTenant)
		tenantGroup.Put("/:id", tenantController.UpdateTenant, middleware.RoleRequired("manager", cfg))
		tenantGroup.Delete("/:id", tenantController.DeleteTenant, middleware.RoleRequired("manager", cfg))
	}

	// Payment routes
	paymentGroup := app.Group("/api/payments", middleware.AuthRequired(cfg))
	{
		paymentGroup.Post("/", paymentController.CreatePayment, middleware.RoleRequired("staff", cfg))
		paymentGroup.Get("/tenant/:tenantId", paymentController.GetPaymentsByTenant)
		paymentGroup.Get("/:id", paymentController.GetPayment)
		paymentGroup.Put("/:id", paymentController.UpdatePayment, middleware.RoleRequired("staff", cfg))
		paymentGroup.Delete("/:id", paymentController.DeletePayment, middleware.RoleRequired("staff", cfg))
	}

	// Maintenance routes
	maintenanceGroup := app.Group("/api/maintenance", middleware.AuthRequired(cfg))
	{
		maintenanceGroup.Post("/", maintenanceController.CreateRequest)
		maintenanceGroup.Get("/room/:roomId", maintenanceController.GetRequestsByRoom)
		maintenanceGroup.Get("/:id", maintenanceController.GetRequest)
		maintenanceGroup.Put("/:id", maintenanceController.UpdateRequest, middleware.RoleRequired("staff", cfg))
		maintenanceGroup.Patch("/:id/status", maintenanceController.UpdateRequestStatus, middleware.RoleRequired("staff", cfg))
		maintenanceGroup.Delete("/:id", maintenanceController.DeleteRequest, middleware.RoleRequired("staff", cfg))
	}

	// Notification routes
	notificationGroup := app.Group("/api/notifications", middleware.AuthRequired(cfg))
	{
		notificationGroup.Post("/", notificationController.CreateNotification, middleware.RoleRequired("admin", cfg))
		notificationGroup.Get("/user/:userId", notificationController.GetUserNotifications)
		notificationGroup.Patch("/:id/read", notificationController.MarkAsRead)
		notificationGroup.Delete("/:id", notificationController.DeleteNotification)
	}

	// Document routes
	documentGroup := app.Group("/api/documents", middleware.AuthRequired(cfg))
	{
		documentGroup.Post("/tenant/:tenantId", documentController.UploadDocument)
		documentGroup.Get("/tenant/:tenantId", documentController.GetTenantDocuments)
		documentGroup.Patch("/:id/verify", documentController.VerifyDocument, middleware.RoleRequired("staff", cfg))
		documentGroup.Delete("/:id", documentController.DeleteDocument)
	}
}
