package handlers

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type TaskHandler struct {
	userRepo  *repository.UserRepository
	taskRepo  *repository.TaskRepository
	transRepo *repository.TransactionRepository
}

func NewTaskHandler(userRepo *repository.UserRepository, taskRepo *repository.TaskRepository, transRepo *repository.TransactionRepository) *TaskHandler {
	return &TaskHandler{
		userRepo:  userRepo,
		taskRepo:  taskRepo,
		transRepo: transRepo,
	}
}

func (h *TaskHandler) GetTasks(c *gin.Context) {
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при загрузке заданий",
		})
		return
	}

	c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) CompleteTask(c *gin.Context) {
	cipher := c.Param("cipher")
	taskID, err := strconv.ParseInt(c.Param("taskId"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Неверный ID задания",
		})
		return
	}

	user, err := h.userRepo.GetByCipher(cipher)
	if err != nil || user == nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Пользователь не найден",
		})
		return
	}

	// Проверяем, выполнено ли уже задание
	completed, err := h.taskRepo.IsTaskCompleted(user.ID, taskID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при проверке задания",
		})
		return
	}
	if completed {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Задание уже выполнено",
		})
		return
	}

	// Получаем информацию о задании
	tasks, err := h.taskRepo.GetAllTasks()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при получении информации о задании",
		})
		return
	}

	var taskReward float64
	var taskTitle string
	for _, t := range tasks {
		if t.ID == taskID {
			taskReward = t.Reward
			taskTitle = t.Title
			break
		}
	}

	if taskReward == 0 {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Задание не найдено",
		})
		return
	}

	// Начисляем награду
	newBalance := user.Balance + taskReward
	if err := h.userRepo.UpdateBalance(user.Cipher, newBalance); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при начислении награды",
		})
		return
	}

	// Отмечаем задание как выполненное
	if err := h.taskRepo.CompleteTask(user.ID, taskID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при отметке задания",
		})
		return
	}

	// Создаем транзакцию для награды
	transaction := &models.Transaction{
		FromCipher:    "system",
		ToCipher:      user.Cipher,
		RecipientName: taskTitle,
		Amount:        taskReward,
		Type:          "task_reward",
	}

	if err := h.transRepo.Create(transaction); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Ошибка при создании транзакции",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":     true,
		"message":     "Задание выполнено!",
		"reward":      taskReward,
		"new_balance": newBalance,
	})
}
