package models

import (
	"database/sql"
	"time"
)

type Task struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	Priority    string    `json:"priority"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DueDate     *time.Time `json:"due_date,omitempty"`
}

type TaskRepository struct {
	DB *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{DB: db}
}

func (r *TaskRepository) GetAll(status, priority string) ([]Task, error) {
	query := "SELECT id, title, description, status, priority, created_at, updated_at, due_date FROM tasks WHERE 1=1"
	args := []interface{}{}

	if status != "" {
		query += " AND status = ?"
		args = append(args, status)
	}
	if priority != "" {
		query += " AND priority = ?"
		args = append(args, priority)
	}

	rows, err := r.DB.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDate sql.NullTime
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt, &dueDate)
		if err != nil {
			return nil, err
		}
		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) GetByID(id int) (*Task, error) {
	var task Task
	var dueDate sql.NullTime
	err := r.DB.QueryRow(
		"SELECT id, title, description, status, priority, created_at, updated_at, due_date FROM tasks WHERE id = ?",
		id,
	).Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt, &dueDate)
	if err != nil {
		return nil, err
	}
	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}
	return &task, nil
}

func (r *TaskRepository) Create(task *Task) error {
	var dueDate interface{}
	if task.DueDate != nil {
		dueDate = task.DueDate.Format("2006-01-02 15:04:05")
	}
	result, err := r.DB.Exec(
		"INSERT INTO tasks (title, description, status, priority, due_date) VALUES (?, ?, ?, ?, ?)",
		task.Title, task.Description, task.Status, task.Priority, dueDate,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	task.ID = int(id)
	return nil
}

func (r *TaskRepository) Update(task *Task) error {
	var dueDate interface{}
	if task.DueDate != nil {
		dueDate = task.DueDate.Format("2006-01-02 15:04:05")
	}
	_, err := r.DB.Exec(
		"UPDATE tasks SET title = ?, description = ?, status = ?, priority = ?, due_date = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?",
		task.Title, task.Description, task.Status, task.Priority, dueDate, task.ID,
	)
	return err
}

func (r *TaskRepository) Delete(id int) error {
	_, err := r.DB.Exec("DELETE FROM tasks WHERE id = ?", id)
	return err
}

func (r *TaskRepository) Search(query string) ([]Task, error) {
	searchPattern := "%" + query + "%"
	rows, err := r.DB.Query(
		"SELECT id, title, description, status, priority, created_at, updated_at, due_date FROM tasks WHERE title LIKE ? OR description LIKE ?",
		searchPattern, searchPattern,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		var dueDate sql.NullTime
		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Status, &task.Priority, &task.CreatedAt, &task.UpdatedAt, &dueDate)
		if err != nil {
			return nil, err
		}
		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

func (r *TaskRepository) GetStats() (map[string]int, error) {
	stats := make(map[string]int)
	
	rows, err := r.DB.Query("SELECT status, COUNT(*) FROM tasks GROUP BY status")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var status string
		var count int
		if err := rows.Scan(&status, &count); err != nil {
			return nil, err
		}
		stats[status] = count
	}

	return stats, nil
}