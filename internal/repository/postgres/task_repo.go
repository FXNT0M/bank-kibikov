package postgres

import (
	"bank-kibikov/internal/models"
	"database/sql"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

func (r *TaskRepository) FindAll() ([]models.Task, error) {
	query := `
		SELECT id, title, reward
		FROM tasks
		ORDER BY id
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []models.Task
	for rows.Next() {
		var t models.Task
		err := rows.Scan(&t.ID, &t.Title, &t.Reward)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, t)
	}

	return tasks, nil
}

func (r *TaskRepository) FindByID(id int) (*models.Task, error) {
	query := `SELECT id, title, reward FROM tasks WHERE id = $1`

	task := &models.Task{}
	err := r.db.QueryRow(query, id).Scan(&task.ID, &task.Title, &task.Reward)
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return task, err
}

func (r *TaskRepository) IsTaskCompleted(userID, taskID int) (bool, error) {
	query := `SELECT completed FROM user_tasks WHERE user_id = $1 AND task_id = $2`

	var completed bool
	err := r.db.QueryRow(query, userID, taskID).Scan(&completed)
	if err == sql.ErrNoRows {
		return false, nil
	}

	return completed, err
}

func (r *TaskRepository) CompleteTask(userID, taskID int) error {
	query := `
		INSERT INTO user_tasks (user_id, task_id, completed, completed_at)
		VALUES ($1, $2, true, CURRENT_TIMESTAMP)
		ON CONFLICT (user_id, task_id) 
		DO UPDATE SET completed = true, completed_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.Exec(query, userID, taskID)
	return err
}
