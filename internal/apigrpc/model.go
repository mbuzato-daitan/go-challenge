package apigrpc

import (
	"github.com/mtbuzato/go-challenge/internal/model"
)

func taskAtob(task model.Task) *Task {
	return &Task{
		Id:        task.ID,
		Name:      task.Name,
		Completed: task.Completed,
	}
}

func taskBtoa(task *Task) model.Task {
	return model.Task{
		ID:        task.Id,
		Name:      task.Name,
		Completed: task.Completed,
	}
}
