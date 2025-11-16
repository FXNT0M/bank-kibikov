package database

import (
	"bank-kibikov/internal/config"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	connStr := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.DBName,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Connect("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Выполняем миграции
	if err := runMigrations(db); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	migrationSQL := `
        CREATE TABLE IF NOT EXISTS accounts (
            id SERIAL PRIMARY KEY,
            account_number VARCHAR(20) UNIQUE NOT NULL,
            owner_name VARCHAR(100) NOT NULL,
            balance DECIMAL(15,2) DEFAULT 0.00 CHECK (balance >= 0),
            currency VARCHAR(3) NOT NULL DEFAULT 'RUB',
            account_type VARCHAR(20) NOT NULL,
            status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE',
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        );
        
        CREATE INDEX IF NOT EXISTS idx_accounts_account_number ON accounts(account_number);
        CREATE INDEX IF NOT EXISTS idx_accounts_owner_name ON accounts(owner_name);
        CREATE INDEX IF NOT EXISTS idx_accounts_status ON accounts(status);
    `

	_, err := db.Exec(migrationSQL)
	return err
}
