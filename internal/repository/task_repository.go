package repository

import (
	"bank-kibikov/internal/models"

	"github.com/jmoiron/sqlx"
)

type TaskRepository struct {
	db *sqlx.DB
}

func NewTaskRepository(db *sqlx.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) GetAllTasks() ([]models.Task, error) {
	var tasks []models.Task
	query := `SELECT * FROM tasks ORDER BY id`
	err := r.db.Select(&tasks, query)
	return tasks, err
}

func (r *TaskRepository) CompleteTask(userID, taskID int64) error {
	query := `INSERT INTO user_tasks (user_id, task_id) VALUES ($1, $2)`
	_, err := r.db.Exec(query, userID, taskID)
	return err
}

func (r *TaskRepository) IsTaskCompleted(userID, taskID int64) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM user_tasks WHERE user_id = $1 AND task_id = $2`
	err := r.db.Get(&count, query, userID, taskID)
	return count > 0, err
}
