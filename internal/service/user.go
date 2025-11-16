package service

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
)

type UserService interface {
	GetUserByCipher(cipher string) (*models.UserResponse, error)
	GetUserTransactions(cipher string) ([]models.Transaction, error)
}

type userService struct {
	userRepo        repository.UserRepository
	transactionRepo repository.TransactionRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetUserByCipher(cipher string) (*models.UserResponse, error) {
	user, err := s.userRepo.FindByCipher(cipher)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	return &models.UserResponse{
		ID:         user.ID,
		Cipher:     user.Cipher,
		Name:       user.Name,
		Group:      user.Group,
		Balance:    user.Balance,
		JoinedDate: user.JoinedDate,
	}, nil
}

func (s *userService) GetUserTransactions(cipher string) ([]models.Transaction, error) {
	return s.transactionRepo.FindByUserCipher(cipher, 5)
}

var (
	ErrUserNotFound = NewServiceError("пользователь не найден")
)
