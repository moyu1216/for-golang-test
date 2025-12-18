package models

import (
	"database/sql"
	"time"
)

// Task 任务结构体
type Task struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Priority    string    `json:"priority"`
	Status      string    `json:"status"`
	DueDate     *time.Time `json:"due_date,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateTask 创建任务
func CreateTask(db *sql.DB, task *Task) (*Task, error) {
	query := `
		INSERT INTO tasks (title, description, priority, status, due_date)
		VALUES (?, ?, ?, ?, ?)
	`
	
	result, err := db.Exec(query, task.Title, task.Description, task.Priority, task.Status, task.DueDate)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return GetTaskByID(db, id)
}

// GetTaskByID 根据ID获取任务
func GetTaskByID(db *sql.DB, id int64) (*Task, error) {
	query := `
		SELECT id, title, description, priority, status, due_date, created_at, updated_at
		FROM tasks
		WHERE id = ?
	`
	
	var task Task
	var dueDate sql.NullTime
	
	err := db.QueryRow(query, id).Scan(
		&task.ID,
		&task.Title,
		&task.Description,
		&task.Priority,
		&task.Status,
		&dueDate,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	
	if err != nil {
		return nil, err
	}

	if dueDate.Valid {
		task.DueDate = &dueDate.Time
	}

	return &task, nil
}

// GetAllTasks 获取所有任务，支持过滤
func GetAllTasks(db *sql.DB, statusFilter, priorityFilter string) ([]*Task, error) {
	query := `
		SELECT id, title, description, priority, status, due_date, created_at, updated_at
		FROM tasks
		WHERE 1=1
	`
	args := []interface{}{}

	if statusFilter != "" {
		query += " AND status = ?"
		args = append(args, statusFilter)
	}

	if priorityFilter != "" {
		query += " AND priority = ?"
		args = append(args, priorityFilter)
	}

	query += " ORDER BY created_at DESC"

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		var dueDate sql.NullTime

		err := rows.Scan(
			&task.ID,
			&task.Title,
			&task.Description,
			&task.Priority,
			&task.Status,
			&dueDate,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if dueDate.Valid {
			task.DueDate = &dueDate.Time
		}

		tasks = append(tasks, &task)
	}

	return tasks, rows.Err()
}

// UpdateTask 更新任务
func UpdateTask(db *sql.DB, id int64, task *Task) (*Task, error) {
	query := `
		UPDATE tasks
		SET title = ?, description = ?, priority = ?, status = ?, due_date = ?, updated_at = CURRENT_TIMESTAMP
		WHERE id = ?
	`

	_, err := db.Exec(query,
		task.Title,
		task.Description,
		task.Priority,
		task.Status,
		task.DueDate,
		id,
	)
	if err != nil {
		return nil, err
	}

	return GetTaskByID(db, id)
}

// DeleteTask 删除任务
func DeleteTask(db *sql.DB, id int64) error {
	query := `DELETE FROM tasks WHERE id = ?`
	_, err := db.Exec(query, id)
	return err
}

// TaskStats 任务统计信息
type TaskStats struct {
	Total       int `json:"total"`
	Todo        int `json:"todo"`
	InProgress  int `json:"in_progress"`
	Completed   int `json:"completed"`
	LowPriority int `json:"low_priority"`
	MediumPriority int `json:"medium_priority"`
	HighPriority int `json:"high_priority"`
}

// GetTaskStats 获取任务统计信息
func GetTaskStats(db *sql.DB) (*TaskStats, error) {
	stats := &TaskStats{}

	// 总数
	err := db.QueryRow("SELECT COUNT(*) FROM tasks").Scan(&stats.Total)
	if err != nil {
		return nil, err
	}

	// 按状态统计
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'todo'").Scan(&stats.Todo)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'in_progress'").Scan(&stats.InProgress)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE status = 'completed'").Scan(&stats.Completed)
	if err != nil {
		return nil, err
	}

	// 按优先级统计
	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE priority = 'low'").Scan(&stats.LowPriority)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE priority = 'medium'").Scan(&stats.MediumPriority)
	if err != nil {
		return nil, err
	}

	err = db.QueryRow("SELECT COUNT(*) FROM tasks WHERE priority = 'high'").Scan(&stats.HighPriority)
	if err != nil {
		return nil, err
	}

	return stats, nil
}










