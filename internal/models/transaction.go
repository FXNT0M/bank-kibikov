package models

import "time"

type Transaction struct {
	ID            int       `json:"id"`
	FromCipher    string    `json:"from_cipher"`
	ToCipher      string    `json:"to_cipher"`
	RecipientName string    `json:"recipient_name"`
	Amount        int       `json:"amount"`
	Date          time.Time `json:"date"`
	Type          string    `json:"type"`
	CreatedAt     time.Time `json:"created_at"`
}

type TransferRequest struct {
	ToCipher      string `json:"to_cipher" binding:"required"`
	Amount        int    `json:"amount" binding:"required,min=1"`
	RecipientName string `json:"recipient_name" binding:"required"`
}
