package service

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
)

type TaskService interface {
	GetAllTasks() ([]models.Task, error)
	CompleteTask(userCipher string, taskID int) error
}

type taskService struct {
	taskRepo        repository.TaskRepository
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewTaskService(taskRepo repository.TaskRepository, userRepo repository.UserRepository, transactionRepo repository.TransactionRepository) TaskService {
	return &taskService{
		taskRepo:        taskRepo,
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

func (s *taskService) GetAllTasks() ([]models.Task, error) {
	return s.taskRepo.FindAll()
}

func (s *taskService) CompleteTask(userCipher string, taskID int) error {
	user, err := s.userRepo.FindByCipher(userCipher)
	if err != nil {
		return err
	}
	if user == nil {
		return ErrUserNotFound
	}

	task, err := s.taskRepo.FindByID(taskID)
	if err != nil {
		return err
	}
	if task == nil {
		return ErrTaskNotFound
	}

	// Check if task already completed
	completed, err := s.taskRepo.IsTaskCompleted(user.ID, taskID)
	if err != nil {
		return err
	}
	if completed {
		return ErrTaskAlreadyCompleted
	}

	// Update user balance
	newBalance := user.Balance + task.Reward
	if err := s.userRepo.UpdateBalance(user.ID, newBalance); err != nil {
		return err
	}

	// Mark task as completed
	if err := s.taskRepo.CompleteTask(user.ID, taskID); err != nil {
		return err
	}

	// Create transaction record
	transaction := &models.Transaction{
		FromCipher:    "system",
		ToCipher:      userCipher,
		RecipientName: task.Title,
		Amount:        task.Reward,
		Type:          "task_reward",
	}

	return s.transactionRepo.Create(transaction)
}

var (
	ErrTaskNotFound         = NewServiceError("задание не найдено")
	ErrTaskAlreadyCompleted = NewServiceError("задание уже выполнено")
)
