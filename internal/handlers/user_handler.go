package handlers

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userRepo  *repository.UserRepository
	transRepo *repository.TransactionRepository
	taskRepo  *repository.TaskRepository
}

func NewUserHandler(userRepo *repository.UserRepository, transRepo *repository.TransactionRepository, taskRepo *repository.TaskRepository) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		transRepo: transRepo,
		taskRepo:  taskRepo,
	}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверные данные",
			"details": err.Error(),
		})
		return
	}

	// Проверяем, нет ли пользователя с таким шифром
	existingUser, err := h.userRepo.GetByCipher(req.Cipher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при проверке пользователя",
		})
		return
	}
	if existingUser != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Пользователь с таким шифром уже существует",
		})
		return
	}

	user := &models.User{
		Cipher:   req.Cipher,
		Password: req.Password,
		Name:     req.Name,
		Group:    req.Group,
		Balance:  1000, // Стартовый бонус
	}

	if err := h.userRepo.Create(user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании пользователя",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"success": true,
		"message": "Регистрация успешна!",
		"user":    user,
	})
}

func (h *UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверные данные",
			"details": err.Error(),
		})
		return
	}

	user, err := h.userRepo.GetByCipher(req.Cipher)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при поиске пользователя",
		})
		return
	}
	if user == nil || user.Password != req.Password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Неверный шифр или пароль",
		})
		return
	}

	// Получаем транзакции пользователя
	transactions, err := h.transRepo.GetUserTransactions(user.Cipher, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при загрузке транзакций",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":      true,
		"user":         user,
		"transactions": transactions,
	})
}

func (h *UserHandler) GetUserData(c *gin.Context) {
	cipher := c.Param("cipher")

	user, err := h.userRepo.GetByCipher(cipher)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Пользователь не найден",
		})
		return
	}

	transactions, err := h.transRepo.GetUserTransactions(cipher, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при загрузке транзакций",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":         user,
		"transactions": transactions,
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.userRepo.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при загрузке пользователей",
		})
		return
	}

	c.JSON(http.StatusOK, users)
}
