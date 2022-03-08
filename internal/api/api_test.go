package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/mtbuzato/go-challenge/internal/model"
	"github.com/stretchr/testify/assert"
)

type StubTaskRepository struct {
	tasks        []model.Task
	createdTasks []model.Task
}

func (r *StubTaskRepository) ListAll() ([]model.Task, error) {
	return r.tasks, nil
}

func (r *StubTaskRepository) ListByCompletion(completed bool) ([]model.Task, error) {
	return r.tasks, nil
}

func (r *StubTaskRepository) Create(name string) (model.Task, error) {
	task := model.Task{ID: "4", Name: name, Completed: false}
	r.createdTasks = append(r.createdTasks, task)
	return task, nil
}

func TestGETTasks(t *testing.T) {
	tasks := []model.Task{
		{ID: "1", Name: "Task 1", Completed: false},
		{ID: "2", Name: "Task 2", Completed: true},
		{ID: "3", Name: "Task 3", Completed: false},
	}

	server := NewAPIServer(&StubTaskRepository{tasks: tasks})

	tests := map[string]struct {
		query          string
		expectedStatus int
		expectedTasks  []model.Task
	}{
		"List all tasks": {
			query:          "",
			expectedStatus: http.StatusOK,
			expectedTasks:  tasks,
		},
		"List completed tasks": {
			query:          "?completed=true",
			expectedStatus: http.StatusOK,
			expectedTasks:  []model.Task{tasks[1]},
		},
		"List incomplete tasks": {
			query:          "?completed=false",
			expectedStatus: http.StatusOK,
			expectedTasks:  []model.Task{tasks[0], tasks[2]},
		},
		"List all tasks with invalid query": {
			query:          "?completed=invalid",
			expectedStatus: http.StatusOK,
			expectedTasks:  tasks,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			req, err := http.NewRequest("GET", "/tasks"+test.query, nil)
			assert.NoError(err)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(test.expectedStatus, w.Code)
		})
	}
}

func TestPOSTTask(t *testing.T) {
	server := NewAPIServer(&StubTaskRepository{})

	tests := map[string]struct {
		body           string
		expectedStatus int
		expectedTask   model.Task
	}{
		"Create a task": {
			body:           `{"name":"Task 4"}`,
			expectedStatus: http.StatusCreated,
			expectedTask:   model.Task{ID: "4", Name: "Task 4", Completed: false},
		},
		"Create a task with invalid JSON": {
			body:           `{"name":"Task 4`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			req, err := http.NewRequest("POST", "/tasks", strings.NewReader(test.body))
			assert.NoError(err)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(test.expectedStatus, w.Code)

			fmt.Println(w.Body.String())

			if test.expectedStatus == http.StatusCreated {
				var task model.Task
				err = json.Unmarshal(w.Body.Bytes(), &task)
				assert.NoError(err)
				assert.Equal(test.expectedTask, task)
			}
		})
	}
}
