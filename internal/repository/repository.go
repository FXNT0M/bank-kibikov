package repository

import "bank-kibikov/internal/models"

type UserRepository interface {
	Create(user *models.User) error
	FindByCipher(cipher string) (*models.User, error)
	UpdateBalance(userID int, newBalance int) error
}

type TransactionRepository interface {
	Create(transaction *models.Transaction) error
	FindByUserCipher(cipher string, limit int) ([]models.Transaction, error)
}

type TaskRepository interface {
	FindAll() ([]models.Task, error)
	FindByID(id int) (*models.Task, error)
	IsTaskCompleted(userID, taskID int) (bool, error)
	CompleteTask(userID, taskID int) error
}
