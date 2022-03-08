package repository

import (
	"database/sql"
	"fmt"

	"github.com/lucsky/cuid"
	"github.com/mtbuzato/go-challenge/internal/model"
)

type TaskRepository struct {
	db *sql.DB
}

func NewTaskRepository(db *sql.DB) *TaskRepository {
	return &TaskRepository{db: db}
}

// Lists all tasks.
func (r *TaskRepository) ListAll() ([]model.Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks")
	if err != nil {
		return nil, fmt.Errorf("Failed to query tasks: %s", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %s", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Lists all tasks with the matching completion status.
func (r *TaskRepository) ListByCompletion(completed bool) ([]model.Task, error) {
	rows, err := r.db.Query("SELECT * FROM tasks WHERE completed = ?", completed)
	if err != nil {
		return nil, fmt.Errorf("Failed to query tasks: %s", err)
	}
	defer rows.Close()

	var tasks []model.Task
	for rows.Next() {
		var task model.Task
		if err := rows.Scan(&task.ID, &task.Name, &task.Completed); err != nil {
			return nil, fmt.Errorf("Failed to scan task: %s", err)
		}
		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Gets a task by ID and returns it.
func (r *TaskRepository) GetByID(id string) (model.Task, error) {
	if err := model.ValidateID(id); err != nil {
		return model.Task{}, err
	}

	var task model.Task
	if err := r.db.QueryRow("SELECT * FROM tasks WHERE id = ?", id).Scan(&task.ID, &task.Name, &task.Completed); err != nil {
		return model.Task{}, fmt.Errorf("Failed to get task by ID: %s", err)
	}

	return task, nil
}

// Creates a new task with the given name and returns it.
func (r *TaskRepository) Create(name string) (model.Task, error) {
	if err := model.ValidateName(name); err != nil {
		return model.Task{}, err
	}

	id := cuid.New()

	if _, err := r.db.Exec("INSERT INTO tasks (id, name, completed) VALUES (?, ?, ?)", id, name, false); err != nil {
		return model.Task{}, fmt.Errorf("Failed to create task: %s", err)
	}

	return model.Task{ID: id, Name: name, Completed: false}, nil
}

// Updates the given task.
func (r *TaskRepository) Update(task model.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	if _, err := r.db.Exec("UPDATE tasks SET name = ?, completed = ? WHERE id = ?", task.Name, task.Completed, task.ID); err != nil {
		return fmt.Errorf("Failed to update task: %s", err)
	}

	return nil
}
