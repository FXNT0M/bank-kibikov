package database

import (
	"github.com/jmoiron/sqlx"
)

func RunMigrations(db *sqlx.DB) error {
	migrations := []string{
		// Таблица пользователей
		`CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            cipher VARCHAR(50) UNIQUE NOT NULL,
            password VARCHAR(100) NOT NULL,
            name VARCHAR(100) NOT NULL,
            group_name VARCHAR(50) NOT NULL,
            balance DECIMAL(15,2) DEFAULT 1000.00 CHECK (balance >= 0),
            joined_date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,

		// Таблица транзакций
		`CREATE TABLE IF NOT EXISTS transactions (
            id SERIAL PRIMARY KEY,
            from_cipher VARCHAR(50) NOT NULL,
            to_cipher VARCHAR(50) NOT NULL,
            recipient_name VARCHAR(100) NOT NULL,
            amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
            date TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            type VARCHAR(20) NOT NULL DEFAULT 'transfer',
            created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
        )`,

		// Таблица заданий
		`CREATE TABLE IF NOT EXISTS tasks (
            id SERIAL PRIMARY KEY,
            title VARCHAR(200) NOT NULL,
            reward DECIMAL(15,2) NOT NULL CHECK (reward > 0),
            completed BOOLEAN DEFAULT FALSE
        )`,

		// Связь пользователей и заданий
		`CREATE TABLE IF NOT EXISTS user_tasks (
            user_id INTEGER REFERENCES users(id),
            task_id INTEGER REFERENCES tasks(id),
            completed_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
            PRIMARY KEY (user_id, task_id)
        )`,
	}

	// Создаем таблицы
	for _, migration := range migrations {
		if _, err := db.Exec(migration); err != nil {
			return err
		}
	}

	// Добавляем тестовые задания
	seedTasks := `
        INSERT INTO tasks (title, reward) VALUES 
        ('Пройти обучение по финансовой грамотности', 50),
        ('Пригласить друга в систему', 100),
        ('Участвовать в опросе сообщества', 25),
        ('Провести код-ревью', 75),
        ('Завершить учебный проект', 150)
        ON CONFLICT DO NOTHING
    `

	_, err := db.Exec(seedTasks)
	return err
}
