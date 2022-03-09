package orm

import (
	"database/sql"
	"fmt"

	"github.com/lucsky/cuid"
	"github.com/mtbuzato/go-challenge/internal/model"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type TaskRepository struct {
	gormDB *gorm.DB
}

func NewTaskRepository(db *sql.DB) (*TaskRepository, error) {
	gormDB, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("Failed to open GORM: %w", err)
	}

	return &TaskRepository{gormDB}, nil
}

// Lists all tasks.
func (r *TaskRepository) ListAll() ([]model.Task, error) {
	tasks := []model.Task{}
	res := r.gormDB.Find(&tasks)
	if res.Error != nil {
		return nil, fmt.Errorf("Failed to query tasks: %w", res.Error)
	}

	return tasks, nil
}

// Lists all tasks with the matching completion status.
func (r *TaskRepository) ListByCompletion(completed bool) ([]model.Task, error) {
	tasks := []model.Task{}
	res := r.gormDB.Where("completed = ?", completed).Find(&tasks)
	if res.Error != nil {
		return nil, fmt.Errorf("Failed to query tasks: %w", res.Error)
	}

	return tasks, nil
}

// Gets a task by ID and returns it.
func (r *TaskRepository) GetByID(id string) (model.Task, error) {
	if err := model.ValidateID(id); err != nil {
		return model.Task{}, err
	}

	var task model.Task
	res := r.gormDB.First(&task, "id = ?", id)
	if res.Error != nil {
		return model.Task{}, fmt.Errorf("Failed to get task by ID: %w", res.Error)
	}

	return task, nil
}

// Creates a new task with the given name and returns it.
func (r *TaskRepository) Create(name string) (model.Task, error) {
	if err := model.ValidateName(name); err != nil {
		return model.Task{}, err
	}

	task := model.Task{ID: cuid.New(), Name: name, Completed: false}
	res := r.gormDB.Create(&task)
	if res.Error != nil {
		return model.Task{}, fmt.Errorf("Failed to create task: %w", res.Error)
	}

	return task, nil
}

// Updates the given task.
func (r *TaskRepository) Update(task model.Task) error {
	if err := task.Validate(); err != nil {
		return err
	}

	res := r.gormDB.Save(&task)
	if res.Error != nil {
		return fmt.Errorf("Failed to update task: %w", res.Error)
	}

	return nil
}
