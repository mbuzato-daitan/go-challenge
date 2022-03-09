package model

import (
	"github.com/lucsky/cuid"
	"github.com/mtbuzato/go-challenge/internal/errors"
)

type Task struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func ValidateID(id string) error {
	if cuid.IsCuid(id) != nil {
		return errors.NewExternalError("Invalid task ID.")
	}

	return nil
}

func ValidateName(name string) error {
	if name == "" {
		return errors.NewExternalError("Invalid task name.")
	}

	if len(name) > 128 {
		return errors.NewExternalError("Task name is too long.")
	}

	return nil
}

func (t *Task) Validate() error {
	err := ValidateID(t.ID)
	if err != nil {
		return err
	}

	err = ValidateName(t.Name)
	if err != nil {
		return err
	}

	return nil
}
