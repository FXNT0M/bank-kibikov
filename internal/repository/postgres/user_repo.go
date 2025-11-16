package postgres

import (
	"bank-kibikov/internal/models"
	"database/sql"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
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

func (r *UserRepository) FindByCipher(cipher string) (*models.User, error) {
	query := `
		SELECT id, cipher, password, name, group_name, balance, joined_date, created_at, updated_at
		FROM users WHERE cipher = $1
	`

	user := &models.User{}
	err := r.db.QueryRow(query, cipher).Scan(
		&user.ID,
		&user.Cipher,
		&user.Password,
		&user.Name,
		&user.Group,
		&user.Balance,
		&user.JoinedDate,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	return user, err
}

func (r *UserRepository) UpdateBalance(userID int, newBalance int) error {
	query := `UPDATE users SET balance = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`
	_, err := r.db.Exec(query, newBalance, userID)
	return err
}
