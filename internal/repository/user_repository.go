package repository

import (
	"bank-kibikov/internal/models"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	query := `
        INSERT INTO users (cipher, password, name, group_name, balance)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, joined_date, created_at, updated_at
    `
	return r.db.QueryRow(
		query,
		user.Cipher,
		user.Password,
		user.Name,
		user.Group,
		user.Balance,
	).Scan(&user.ID, &user.JoinedDate, &user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) GetByCipher(cipher string) (*models.User, error) {
	var user models.User
	query := `SELECT * FROM users WHERE cipher = $1`
	err := r.db.Get(&user, query, cipher)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return &user, err
}

func (r *UserRepository) UpdateBalance(cipher string, newBalance float64) error {
	query := `UPDATE users SET balance = $1, updated_at = NOW() WHERE cipher = $2`
	_, err := r.db.Exec(query, newBalance, cipher)
	return err
}

func (r *UserRepository) GetAllUsers() ([]models.User, error) {
	var users []models.User
	query := `SELECT * FROM users ORDER BY name`
	err := r.db.Select(&users, query)
	return users, err
}
