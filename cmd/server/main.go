package main

import (
	"bank-kibikov/internal/config"
	"bank-kibikov/internal/handlers"
	"bank-kibikov/internal/repository"
	"bank-kibikov/pkg/database"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Подключение к PostgreSQL
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(db)
	transRepo := repository.NewTransactionRepository(db)
	taskRepo := repository.NewTaskRepository(db)

	// Инициализация обработчиков
	userHandler := handlers.NewUserHandler(userRepo, transRepo, taskRepo)
	transactionHandler := handlers.NewTransactionHandler(userRepo, transRepo)
	taskHandler := handlers.NewTaskHandler(userRepo, taskRepo, transRepo)

	// Настройка маршрутов
	router := gin.Default()

	// CORS middleware
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Обслуживание статических файлов фронтенда
	router.Static("/static", "./static")
	router.StaticFile("/", "./static/index.html")
	router.StaticFile("/index.html", "./static/index.html")

	// API маршруты
	api := router.Group("/api")
	{
		// Пользователи
		users := api.Group("/users")
		{
			users.POST("/register", userHandler.Register)
			users.POST("/login", userHandler.Login)
			users.GET("/:cipher", userHandler.GetUserData)
			users.GET("/", userHandler.GetAllUsers)
		}

		// Переводы
		transfers := api.Group("/transfers")
		{
			transfers.POST("/:cipher", transactionHandler.Transfer)
		}

		// Задания
		tasks := api.Group("/tasks")
		{
			tasks.GET("/", taskHandler.GetTasks)
			tasks.POST("/:cipher/complete/:taskId", taskHandler.CompleteTask)
		}
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":   "OK",
			"service":  "Bank Kibikov API",
			"database": "PostgreSQL",
		})
	})

	// Fallback для SPA
	router.NoRoute(func(c *gin.Context) {
		c.File("./static/index.html")
	})

	log.Printf("🚀 Server starting on port %s", cfg.Server.Port)
	log.Printf("📊 Database: PostgreSQL")
	log.Printf("🌐 Frontend available at: http://localhost%s", cfg.Server.Port)

	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
