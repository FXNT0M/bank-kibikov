package handler

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req models.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	user := &models.User{
		Cipher:   req.Cipher,
		Password: req.Password,
		Name:     req.Name,
		Group:    req.Group,
		Balance:  1000, // Стартовый бонус
	}

	token, err := h.authService.Register(user)
	if err != nil {
		if serviceErr, ok := err.(*service.ServiceError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": serviceErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Регистрация успешна!",
		"token":   token,
		"user": models.UserResponse{
			ID:         user.ID,
			Cipher:     user.Cipher,
			Name:       user.Name,
			Group:      user.Group,
			Balance:    user.Balance,
			JoinedDate: user.JoinedDate,
		},
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req models.UserLoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	user, token, err := h.authService.Login(req.Cipher, req.Password)
	if err != nil {
		if serviceErr, ok := err.(*service.ServiceError); ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": serviceErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Вход выполнен успешно",
		"token":   token,
		"user": models.UserResponse{
			ID:         user.ID,
			Cipher:     user.Cipher,
			Name:       user.Name,
			Group:      user.Group,
			Balance:    user.Balance,
			JoinedDate: user.JoinedDate,
		},
	})
}
