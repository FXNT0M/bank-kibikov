package main

import (
	"bank-kibikov/internal/config"
	"bank-kibikov/internal/handler"
	"bank-kibikov/internal/middleware"
	"bank-kibikov/internal/repository/postgres"
	"bank-kibikov/internal/service"
	"bank-kibikov/pkg/database"
	"bank-kibikov/pkg/jwt"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.Load()

	// Initialize database
	db, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize JWT
	jwtManager := jwt.NewJWTManager(cfg.JWTSecret)

	// Initialize repositories
	userRepo := postgres.NewUserRepository(db)
	transactionRepo := postgres.NewTransactionRepository(db)
	taskRepo := postgres.NewTaskRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, jwtManager)
	userService := service.NewUserService(userRepo)
	transactionService := service.NewTransactionService(transactionRepo, userRepo)
	taskService := service.NewTaskService(taskRepo, userRepo, transactionRepo)

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authService)
	userHandler := handler.NewUserHandler(userService)
	transactionHandler := handler.NewTransactionHandler(transactionService)
	taskHandler := handler.NewTaskHandler(taskService)

	// Setup router
	router := gin.Default()

	// CORS middleware - ИСПРАВЛЕННАЯ КОНФИГУРАЦИЯ
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Accept"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// ИЛИ используйте упрощенную версию:
	// router.Use(cors.Default())

	// API routes
	api := router.Group("/api/v1")
	{
		// Public routes
		auth := api.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
		}

		// Protected routes
		protected := api.Group("/")
		protected.Use(middleware.AuthMiddleware(jwtManager))
		{
			protected.GET("/user/profile", userHandler.GetProfile)
			protected.GET("/user/transactions", userHandler.GetTransactions)

			protected.POST("/transactions/transfer", transactionHandler.MakeTransfer)

			protected.GET("/tasks", taskHandler.GetTasks)
			protected.POST("/tasks/:id/complete", taskHandler.CompleteTask)
		}
	}

	log.Printf("Server starting on port %s", cfg.Port)
	log.Fatal(router.Run(":" + cfg.Port))
}
