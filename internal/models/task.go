package models

type Task struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Reward    int    `json:"reward"`
	Completed bool   `json:"completed"`
}

type UserTask struct {
	UserID    int  `json:"user_id"`
	TaskID    int  `json:"task_id"`
	Completed bool `json:"completed"`
}
