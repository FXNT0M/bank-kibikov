package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgresConnection(connectionString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL database")
	return db, nil
}

func RunMigrations(connectionString string) error {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// Create tables if they don't exist
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			cipher VARCHAR(50) UNIQUE NOT NULL,
			password VARCHAR(255) NOT NULL,
			name VARCHAR(100) NOT NULL,
			group_name VARCHAR(50) NOT NULL,
			balance INTEGER DEFAULT 1000,
			joined_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			from_cipher VARCHAR(50) NOT NULL,
			to_cipher VARCHAR(50) NOT NULL,
			recipient_name VARCHAR(100) NOT NULL,
			amount INTEGER NOT NULL,
			date TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
			type VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS tasks (
			id SERIAL PRIMARY KEY,
			title VARCHAR(255) NOT NULL,
			reward INTEGER NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS user_tasks (
			user_id INTEGER REFERENCES users(id),
			task_id INTEGER REFERENCES tasks(id),
			completed BOOLEAN DEFAULT FALSE,
			completed_at TIMESTAMP,
			PRIMARY KEY (user_id, task_id)
		);

		-- Insert default tasks
		INSERT INTO tasks (id, title, reward) VALUES
		(1, 'Пройти обучение по финансовой грамотности', 50),
		(2, 'Пригласить друга в систему', 100),
		(3, 'Участвовать в опросе сообщества', 25),
		(4, 'Провести код-ревью', 75),
		(5, 'Завершить учебный проект', 150)
		ON CONFLICT (id) DO NOTHING;
	`)

	return err
}
