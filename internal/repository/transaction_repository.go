package repository

import (
	"bank-kibikov/internal/models"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository struct {
	db *sqlx.DB
}

func NewTransactionRepository(db *sqlx.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(transaction *models.Transaction) error {
	query := `
        INSERT INTO transactions (from_cipher, to_cipher, recipient_name, amount, type)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, date, created_at
    `
	return r.db.QueryRow(
		query,
		transaction.FromCipher,
		transaction.ToCipher,
		transaction.RecipientName,
		transaction.Amount,
		transaction.Type,
	).Scan(&transaction.ID, &transaction.Date, &transaction.CreatedAt)
}

func (r *TransactionRepository) GetUserTransactions(cipher string, limit int) ([]models.Transaction, error) {
	var transactions []models.Transaction
	query := `
        SELECT * FROM transactions 
        WHERE from_cipher = $1 OR to_cipher = $1 
        ORDER BY date DESC 
        LIMIT $2
    `
	err := r.db.Select(&transactions, query, cipher, limit)
	return transactions, err
}
