package handler

import (
	"bank-kibikov/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	cipher, exists := c.Get("userCipher")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	user, err := h.userService.GetUserByCipher(cipher.(string))
	if err != nil {
		if serviceErr, ok := err.(*service.ServiceError); ok {
			c.JSON(http.StatusNotFound, gin.H{"error": serviceErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func (h *UserHandler) GetTransactions(c *gin.Context) {
	cipher, exists := c.Get("userCipher")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	transactions, err := h.userService.GetUserTransactions(cipher.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"transactions": transactions})
}
