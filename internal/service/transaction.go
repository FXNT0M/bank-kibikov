package service

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
)

type TransactionService interface {
	MakeTransfer(fromCipher, toCipher string, amount int, recipientName string) error
}

type transactionService struct {
	transactionRepo repository.TransactionRepository
	userRepo        repository.UserRepository
}

func NewTransactionService(transactionRepo repository.TransactionRepository, userRepo repository.UserRepository) TransactionService {
	return &transactionService{
		transactionRepo: transactionRepo,
		userRepo:        userRepo,
	}
}

func (s *transactionService) MakeTransfer(fromCipher, toCipher string, amount int, recipientName string) error {
	// Validate users
	fromUser, err := s.userRepo.FindByCipher(fromCipher)
	if err != nil {
		return err
	}
	if fromUser == nil {
		return ErrSenderNotFound
	}

	toUser, err := s.userRepo.FindByCipher(toCipher)
	if err != nil {
		return err
	}
	if toUser == nil {
		return ErrRecipientNotFound
	}

	// Check balance
	if fromUser.Balance < amount {
		return ErrInsufficientFunds
	}

	// Update balances
	if err := s.userRepo.UpdateBalance(fromUser.ID, fromUser.Balance-amount); err != nil {
		return err
	}
	if err := s.userRepo.UpdateBalance(toUser.ID, toUser.Balance+amount); err != nil {
		return err
	}

	// Create transaction record
	transaction := &models.Transaction{
		FromCipher:    fromCipher,
		ToCipher:      toCipher,
		RecipientName: recipientName,
		Amount:        amount,
		Type:          "transfer",
	}

	return s.transactionRepo.Create(transaction)
}

var (
	ErrSenderNotFound    = NewServiceError("отправитель не найден")
	ErrRecipientNotFound = NewServiceError("получатель не найден")
	ErrInsufficientFunds = NewServiceError("недостаточно средств")
)
