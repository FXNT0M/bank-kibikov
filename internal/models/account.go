package models

import (
	"time"
)

type Account struct {
	ID            int64     `json:"id" db:"id"`
	AccountNumber string    `json:"account_number" db:"account_number"`
	OwnerName     string    `json:"owner_name" db:"owner_name"`
	Balance       float64   `json:"balance" db:"balance"`
	Currency      string    `json:"currency" db:"currency"`
	AccountType   string    `json:"account_type" db:"account_type"`
	Status        string    `json:"status" db:"status"`
	CreatedAt     time.Time `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

type CreateAccountRequest struct {
	AccountNumber string  `json:"account_number" binding:"required,min=10,max=20"`
	OwnerName     string  `json:"owner_name" binding:"required,min=2,max=100"`
	Balance       float64 `json:"balance" binding:"min=0"`
	Currency      string  `json:"currency" binding:"required,oneof=RUB USD EUR"`
	AccountType   string  `json:"account_type" binding:"required,oneof=CHECKING SAVINGS CREDIT BUSINESS"`
}

type UpdateAccountRequest struct {
	OwnerName *string  `json:"owner_name,omitempty" binding:"omitempty,min=2,max=100"`
	Balance   *float64 `json:"balance,omitempty" binding:"omitempty,min=0"`
	Status    *string  `json:"status,omitempty" binding:"omitempty,oneof=ACTIVE BLOCKED CLOSED"`
}
