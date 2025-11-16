package service

import (
	"bank-kibikov/internal/models"
	"bank-kibikov/internal/repository"
	"bank-kibikov/pkg/jwt"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(user *models.User) (string, error)
	Login(cipher, password string) (*models.User, string, error)
}

type authService struct {
	userRepo repository.UserRepository
	jwt      *jwt.JWTManager
}

func NewAuthService(userRepo repository.UserRepository, jwtManager *jwt.JWTManager) AuthService {
	return &authService{
		userRepo: userRepo,
		jwt:      jwtManager,
	}
}

func (s *authService) Register(user *models.User) (string, error) {
	// Check if user already exists
	existingUser, err := s.userRepo.FindByCipher(user.Cipher)
	if err != nil {
		return "", err
	}
	if existingUser != nil {
		return "", ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := s.userRepo.Create(user); err != nil {
		return "", err
	}

	// Generate JWT token
	token, err := s.jwt.GenerateToken(user.ID, user.Cipher)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *authService) Login(cipher, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByCipher(cipher)
	if err != nil {
		return nil, "", err
	}
	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	// Verify password
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidCredentials
	}

	// Generate JWT token
	token, err := s.jwt.GenerateToken(user.ID, user.Cipher)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

var (
	ErrUserAlreadyExists  = NewServiceError("пользователь с таким шифром уже существует")
	ErrInvalidCredentials = NewServiceError("неверный шифр или пароль")
)

type ServiceError struct {
	Message string
}

func NewServiceError(message string) *ServiceError {
	return &ServiceError{Message: message}
}

func (e *ServiceError) Error() string {
	return e.Message
}
