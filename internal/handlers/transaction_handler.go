package handlers

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type TransactionHandler struct {
	userRepo  *repository.UserRepository
	transRepo *repository.TransactionRepository
}

func NewTransactionHandler(userRepo *repository.UserRepository, transRepo *repository.TransactionRepository) *TransactionHandler {
	return &TransactionHandler{
		userRepo:  userRepo,
		transRepo: transRepo,
	}
}

func (h *TransactionHandler) Transfer(c *gin.Context) {
	var req models.TransferRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Неверные данные",
			"details": err.Error(),
		})
		return
	}

	fromCipher := c.Param("cipher")

	// Получаем отправителя
	fromUser, err := h.userRepo.GetByCipher(fromCipher)
	if err != nil || fromUser == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Отправитель не найден",
		})
		return
	}

	// Получаем получателя
	toUser, err := h.userRepo.GetByCipher(req.ToCipher)
	if err != nil || toUser == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Получатель не найден",
		})
		return
	}

	// Проверяем баланс
	if fromUser.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Недостаточно средств",
		})
		return
	}

	// Выполняем перевод
	fromUser.Balance -= req.Amount
	toUser.Balance += req.Amount

	if err := h.userRepo.UpdateBalance(fromUser.Cipher, fromUser.Balance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при обновлении баланса отправителя",
		})
		return
	}

	if err := h.userRepo.UpdateBalance(toUser.Cipher, toUser.Balance); err != nil {
		// Откатываем баланс отправителя в случае ошибки
		h.userRepo.UpdateBalance(fromUser.Cipher, fromUser.Balance+req.Amount)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при обновлении баланса получателя",
		})
		return
	}

	// Создаем запись о транзакции
	transaction := &models.Transaction{
		FromCipher:    fromUser.Cipher,
		ToCipher:      toUser.Cipher,
		RecipientName: toUser.Name,
		Amount:        req.Amount,
		Type:          "transfer",
	}

	if err := h.transRepo.Create(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании транзакции",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Перевод успешно выполнен",
		"new_balance": fromUser.Balance,
	})
}
