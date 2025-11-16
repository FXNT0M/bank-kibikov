package models

import (
	"time"
)

type User struct {
	ID         int64     `json:"id" db:"id"`
	Cipher     string    `json:"cipher" db:"cipher"`
	Password   string    `json:"password" db:"password"`
	Name       string    `json:"name" db:"name"`
	Group      string    `json:"group" db:"group_name"`
	Balance    float64   `json:"balance" db:"balance"`
	JoinedDate time.Time `json:"joined_date" db:"joined_date"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
}

type Transaction struct {
	ID            int64     `json:"id" db:"id"`
	FromCipher    string    `json:"from_cipher" db:"from_cipher"`
	ToCipher      string    `json:"to_cipher" db:"to_cipher"`
	RecipientName string    `json:"recipient_name" db:"recipient_name"`
	Amount        float64   `json:"amount" db:"amount"`
	Date          time.Time `json:"date" db:"date"`
	Type          string    `json:"type" db:"type"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
}

type Task struct {
	ID        int64   `json:"id" db:"id"`
	Title     string  `json:"title" db:"title"`
	Reward    float64 `json:"reward" db:"reward"`
	Completed bool    `json:"completed" db:"completed"`
}

type UserTask struct {
	UserID      int64     `json:"user_id" db:"user_id"`
	TaskID      int64     `json:"task_id" db:"task_id"`
	CompletedAt time.Time `json:"completed_at" db:"completed_at"`
}

// Запросы
type RegisterRequest struct {
	Cipher   string `json:"cipher" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Group    string `json:"group" binding:"required"`
}

type LoginRequest struct {
	Cipher   string `json:"cipher" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type TransferRequest struct {
	ToCipher string  `json:"to_cipher" binding:"required"`
	Amount   float64 `json:"amount" binding:"required,min=0.01"`
}
