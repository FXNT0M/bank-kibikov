package repository

import (
	"bank-kibikov/internal/models"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
)

type AccountRepository struct {
	db *sqlx.DB
}

func NewAccountRepository(db *sqlx.DB) *AccountRepository {
	return &AccountRepository{db: db}
}

func (r *AccountRepository) Create(account *models.Account) error {
	query := `
        INSERT INTO accounts (account_number, owner_name, balance, currency, account_type, status)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id, created_at, updated_at
    `

	return r.db.QueryRow(
		query,
		account.AccountNumber,
		account.OwnerName,
		account.Balance,
		account.Currency,
		account.AccountType,
		account.Status,
	).Scan(&account.ID, &account.CreatedAt, &account.UpdatedAt)
}

func (r *AccountRepository) GetAll() ([]models.Account, error) {
	var accounts []models.Account
	query := `SELECT * FROM accounts ORDER BY created_at DESC`
	err := r.db.Select(&accounts, query)
	return accounts, err
}

func (r *AccountRepository) GetByID(id int64) (*models.Account, error) {
	var account models.Account
	query := `SELECT * FROM accounts WHERE id = $1`
	err := r.db.Get(&account, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &account, err
}

func (r *AccountRepository) Update(id int64, req *models.UpdateAccountRequest) error {
	query := `UPDATE accounts SET updated_at = NOW()`
	args := []interface{}{}
	argCount := 1

	if req.OwnerName != nil {
		query += fmt.Sprintf(", owner_name = $%d", argCount)
		args = append(args, *req.OwnerName)
		argCount++
	}

	if req.Balance != nil {
		query += fmt.Sprintf(", balance = $%d", argCount)
		args = append(args, *req.Balance)
		argCount++
	}

	if req.Status != nil {
		query += fmt.Sprintf(", status = $%d", argCount)
		args = append(args, *req.Status)
		argCount++
	}

	query += fmt.Sprintf(" WHERE id = $%d", argCount)
	args = append(args, id)

	result, err := r.db.Exec(query, args...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AccountRepository) Delete(id int64) error {
	query := `DELETE FROM accounts WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *AccountRepository) GetByAccountNumber(accountNumber string) (*models.Account, error) {
	var account models.Account
	query := `SELECT * FROM accounts WHERE account_number = $1`
	err := r.db.Get(&account, query, accountNumber)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &account, err
}
