package postgres

import (
	"bank-kibikov/internal/models"
	"database/sql"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
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

func (r *TransactionRepository) FindByUserCipher(cipher string, limit int) ([]models.Transaction, error) {
	query := `
		SELECT id, from_cipher, to_cipher, recipient_name, amount, date, type, created_at
		FROM transactions 
		WHERE from_cipher = $1 OR to_cipher = $1
		ORDER BY created_at DESC
		LIMIT $2
	`

	rows, err := r.db.Query(query, cipher, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []models.Transaction
	for rows.Next() {
		var t models.Transaction
		err := rows.Scan(
			&t.ID,
			&t.FromCipher,
			&t.ToCipher,
			&t.RecipientName,
			&t.Amount,
			&t.Date,
			&t.Type,
			&t.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t)
	}

	return transactions, nil
}
