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

	// Подключение к базе данных
	db, err := database.NewPostgresDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Инициализация репозитория и обработчиков
	accountRepo := repository.NewAccountRepository(db)
	accountHandler := handlers.NewAccountHandler(accountRepo)

	// Настройка маршрутов
	router := gin.Default()

	// Группа маршрутов для работы со счетами
	accounts := router.Group("/accounts")
	{
		accounts.POST("", accountHandler.CreateAccount)
		accounts.GET("", accountHandler.GetAccounts)
		accounts.GET("/:id", accountHandler.GetAccount)
		accounts.PUT("/:id", accountHandler.UpdateAccount)
		accounts.DELETE("/:id", accountHandler.DeleteAccount)
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "OK",
			"service": "Bank Kibikov API",
		})
	})

	log.Printf("Server starting on port %s", cfg.Server.Port)
	if err := router.Run(cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
