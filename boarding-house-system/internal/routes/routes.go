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
		authGroup.Get("/me", middleware.AuthRequired(cfg), authController.Me)
	}

	// User routes
	userGroup := app.Group("/api/users", middleware.AuthRequired(cfg))
	{
		userGroup.Get("/", middleware.RoleRequired("admin", cfg), userController.GetAllUsers)
		userGroup.Get("/:id", userController.GetUser)
		userGroup.Put("/:id", userController.UpdateUser)
		userGroup.Delete("/:id", middleware.RoleRequired("admin", cfg), userController.DeleteUser)

		// User profile routes
		userGroup.Post("/:userId/profile", userController.CreateProfile)
		userGroup.Get("/:userId/profile", userController.GetProfile)
	}

	// Boarding house routes
	houseGroup := app.Group("/api/houses", middleware.AuthRequired(cfg))
	{
		houseGroup.Post("/", middleware.RoleRequired("admin", cfg), houseController.CreateHouse)
		houseGroup.Get("/", houseController.GetAllHouses)
		houseGroup.Get("/:id", houseController.GetHouse)
		houseGroup.Put("/:id", middleware.RoleRequired("admin", cfg), houseController.UpdateHouse)
		houseGroup.Delete("/:id", middleware.RoleRequired("admin", cfg), houseController.DeleteHouse)
	}

	// Room routes
	roomGroup := app.Group("/api/rooms", middleware.AuthRequired(cfg))
	{
		roomGroup.Post("/", middleware.RoleRequired("manager", cfg), roomController.CreateRoom)
		roomGroup.Get("/", roomController.GetAllRooms)
		roomGroup.Get("/:id", roomController.GetRoom)
		roomGroup.Put("/:id", middleware.RoleRequired("manager", cfg), roomController.UpdateRoom)
		roomGroup.Delete("/:id", middleware.RoleRequired("manager", cfg), roomController.DeleteRoom)
	}

	// Tenant routes
	tenantGroup := app.Group("/api/tenants", middleware.AuthRequired(cfg))
	{
		tenantGroup.Post("/", middleware.RoleRequired("manager", cfg), tenantController.CreateTenant)
		tenantGroup.Get("/house/:houseId", tenantController.GetTenantsByHouse)
		tenantGroup.Get("/:id", tenantController.GetTenant)
		tenantGroup.Put("/:id", middleware.RoleRequired("manager", cfg), tenantController.UpdateTenant)
		tenantGroup.Delete("/:id", middleware.RoleRequired("manager", cfg), tenantController.DeleteTenant)
	}

	// Payment routes
	paymentGroup := app.Group("/api/payments", middleware.AuthRequired(cfg))
	{
		paymentGroup.Post("/", middleware.RoleRequired("staff", cfg), paymentController.CreatePayment)
		paymentGroup.Get("/tenant/:tenantId", paymentController.GetPaymentsByTenant)
		paymentGroup.Get("/:id", paymentController.GetPayment)
		paymentGroup.Put("/:id", middleware.RoleRequired("staff", cfg), paymentController.UpdatePayment)
		paymentGroup.Delete("/:id", middleware.RoleRequired("staff", cfg), paymentController.DeletePayment)
	}

	// Maintenance routes
	maintenanceGroup := app.Group("/api/maintenance", middleware.AuthRequired(cfg))
	{
		maintenanceGroup.Post("/", maintenanceController.CreateRequest)
		maintenanceGroup.Get("/room/:roomId", maintenanceController.GetRequestsByRoom)
		maintenanceGroup.Get("/:id", maintenanceController.GetRequest)
		maintenanceGroup.Put("/:id", middleware.RoleRequired("staff", cfg), maintenanceController.UpdateRequest)
		maintenanceGroup.Patch("/:id/status", middleware.RoleRequired("staff", cfg), maintenanceController.UpdateRequestStatus)
		maintenanceGroup.Delete("/:id", middleware.RoleRequired("staff", cfg), maintenanceController.DeleteRequest)
	}

	// Notification routes
	notificationGroup := app.Group("/api/notifications", middleware.AuthRequired(cfg))
	{
		notificationGroup.Post("/", middleware.RoleRequired("admin", cfg), notificationController.CreateNotification)
		notificationGroup.Get("/user/:userId", notificationController.GetUserNotifications)
		notificationGroup.Patch("/:id/read", notificationController.MarkAsRead)
		notificationGroup.Delete("/:id", notificationController.DeleteNotification)
	}

	// Document routes
	documentGroup := app.Group("/api/documents", middleware.AuthRequired(cfg))
	{
		documentGroup.Post("/tenant/:tenantId", documentController.UploadDocument)
		documentGroup.Get("/tenant/:tenantId", documentController.GetTenantDocuments)
		documentGroup.Patch("/:id/verify", middleware.RoleRequired("staff", cfg), documentController.VerifyDocument)
		documentGroup.Delete("/:id", documentController.DeleteDocument)
	}
}
