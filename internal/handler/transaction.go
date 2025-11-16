package handler

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	transactionService service.TransactionService
}

func NewTransactionHandler(transactionService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{transactionService: transactionService}
}

func (h *TransactionHandler) MakeTransfer(c *gin.Context) {
	cipher, exists := c.Get("userCipher")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Не авторизован"})
		return
	}

	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Неверные данные"})
		return
	}

	err := h.transactionService.MakeTransfer(cipher.(string), req.ToCipher, req.Amount, req.RecipientName)
	if err != nil {
		if serviceErr, ok := err.(*service.ServiceError); ok {
			c.JSON(http.StatusBadRequest, gin.H{"error": serviceErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Ошибка сервера"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Перевод успешно выполнен"})
}
