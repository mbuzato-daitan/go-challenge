package repository

import (
	"database/sql"
	"fmt"

	"github.com/lucsky/cuid"
)

type Task struct {
	ID        string
	Name      string
	Completed bool
}

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Lists all tasks.
func (r *TaskRepository) ListAll() ([]Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("Failed to query tasks: %s", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %s", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Lists all tasks with the matching completion status.
func (r *TaskRepository) ListByCompletion(completed bool) ([]Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks WHERE completed = ?", completed)
	if err != nil {
		return nil, fmt.Errorf("Failed to query tasks: %s", err)
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %s", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Gets a task by ID and returns it.
func (r *TaskRepository) GetByID(id string) (Task, error) {
	if err := cuid.IsCuid(id); err != nil {
		return Task{}, fmt.Errorf("Invalid task ID: %s", err)
	}

	var task Task
	if err := r.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Name, &task.Completed); err != nil {
		return Task{}, fmt.Errorf("Failed to get task by ID: %s", err)
	}

	return task, nil
}

// Creates a new task with the given name and returns it.
func (r *TaskRepository) Create(name string) (Task, error) {
	if len(name) < 1 {
		return Task{}, fmt.Errorf("Invalid task name: %s", name)
	}

	id := cuid.New()

	if _, err := r.db.Exec("INSERT INTO tasks (id, name, completed) VALUES (?, ?, ?)", id, name, false); err != nil {
		return Task{}, fmt.Errorf("Failed to create task: %s", err)
	}

	return Task{ID: id, Name: name, Completed: false}, nil
}

// Updates the given task.
func (r *TaskRepository) Update(task Task) error {
	if err := cuid.IsCuid(task.ID); err != nil {
		return fmt.Errorf("Invalid task ID: %s", err)
	}

	if len(task.Name) < 1 {
		return fmt.Errorf("Invalid task name: %s", task.Name)
	}

	if _, err := r.db.Exec("UPDATE tasks SET name = ?, completed = ? WHERE id = ?", task.Name, task.Completed, task.ID); err != nil {
		return fmt.Errorf("Failed to update task: %s", err)
	}

	return nil
}
