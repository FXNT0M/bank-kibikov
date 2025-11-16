package handlers

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AccountHandler struct {
	repo *repository.AccountRepository
}

func NewAccountHandler(repo *repository.AccountRepository) *AccountHandler {
	return &AccountHandler{repo: repo}
}

func (h *AccountHandler) CreateAccount(c *gin.Context) {
	var req models.CreateAccountRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	// Проверяем, существует ли счет с таким номером
	existingAccount, err := h.repo.GetByAccountNumber(req.AccountNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to check account existence",
		})
		return
	}

	if existingAccount != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Account with this number already exists",
		})
		return
	}

	account := &models.Account{
		AccountNumber: req.AccountNumber,
		OwnerName:     req.OwnerName,
		Balance:       req.Balance,
		Currency:      req.Currency,
		AccountType:   req.AccountType,
		Status:        "ACTIVE",
	}

	if err := h.repo.Create(account); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create account",
		})
		return
	}

	c.JSON(http.StatusCreated, account)
}

func (h *AccountHandler) GetAccounts(c *gin.Context) {
	accounts, err := h.repo.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch accounts",
		})
		return
	}

	c.JSON(http.StatusOK, accounts)
}

func (h *AccountHandler) GetAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		return
	}

	account, err := h.repo.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch account",
		})
		return
	}

	if account == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Account not found",
		})
		return
	}

	c.JSON(http.StatusOK, account)
}

func (h *AccountHandler) UpdateAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		return
	}

	var req models.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input data",
			"details": err.Error(),
		})
		return
	}

	if err := h.repo.Update(id, &req); err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to update account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account updated successfully",
	})
}

func (h *AccountHandler) DeleteAccount(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid account ID",
		})
		return
	}

	if err := h.repo.Delete(id); err != nil {
		if err.Error() == "sql: no rows in result set" {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "Account not found",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete account",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Account deleted successfully",
	})
}
