package apivanilla

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/mtbuzato/go-challenge/internal/errors"
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

func (r *StubTaskRepository) GetByID(id string) (model.Task, error) {
	var task model.Task
	for _, t := range r.tasks {
		if t.ID == id {
			task = t
			break
		}
	}

	if task.Name == "" {
		return task, errors.NewHTTPError("Task not found.", http.StatusNotFound)
	}

	return task, nil
}

func (r *StubTaskRepository) Create(name string) (model.Task, error) {
	task := model.Task{ID: "4", Name: name, Completed: false}
	r.createdTasks = append(r.createdTasks, task)
	return task, nil
}

func (r *StubTaskRepository) Update(task model.Task) error {
	found := false

	for i, t := range r.tasks {
		if t.ID == task.ID {
			r.tasks[i] = task
			found = true
			break
		}
	}

	if !found {
		return errors.NewHTTPError("Task not found.", http.StatusNotFound)
	}

	return nil
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
			req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))
			assert.NoError(err)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(test.expectedStatus, w.Code)
		})
	}
}

func BenchmarkGetTasks(b *testing.B) {
	tasks := []model.Task{
		{ID: "1", Name: "Task 1", Completed: false},
		{ID: "2", Name: "Task 2", Completed: true},
		{ID: "3", Name: "Task 3", Completed: false},
	}

	server := NewAPIServer(&StubTaskRepository{tasks: tasks})

	req, _ := http.NewRequest("GET", "/tasks", nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))

	w := httptest.NewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		server.ServeHTTP(w, req)
	}
}

func TestGETTask(t *testing.T) {
	tasks := []model.Task{
		{ID: "1", Name: "Task 1", Completed: false},
		{ID: "2", Name: "Task 2", Completed: true},
		{ID: "3", Name: "Task 3", Completed: false},
	}

	server := NewAPIServer(&StubTaskRepository{tasks: tasks})

	tests := map[string]struct {
		id             string
		expectedStatus int
		expectedTask   model.Task
	}{
		"Get a task": {
			id:             "1",
			expectedStatus: http.StatusOK,
			expectedTask:   tasks[0],
		},
		"Get a task that does not exist": {
			id:             "4",
			expectedStatus: http.StatusNotFound,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			req, err := http.NewRequest("GET", "/tasks/"+test.id, nil)
			req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))
			assert.NoError(err)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(test.expectedStatus, w.Code)

			if test.expectedStatus == http.StatusOK {
				var task model.Task
				err = json.Unmarshal(w.Body.Bytes(), &task)
				assert.NoError(err)
				assert.Equal(test.expectedTask, task)
			}
		})
	}
}

func BenchmarkGETTask(b *testing.B) {
	tasks := []model.Task{
		{ID: "1", Name: "Task 1", Completed: false},
	}

	server := NewAPIServer(&StubTaskRepository{tasks: tasks})

	req, _ := http.NewRequest("GET", "/tasks/1", nil)
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))

	w := httptest.NewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		server.ServeHTTP(w, req)
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
			req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))
			assert.NoError(err)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(test.expectedStatus, w.Code)

			if test.expectedStatus == http.StatusCreated {
				var task model.Task
				err = json.Unmarshal(w.Body.Bytes(), &task)
				assert.NoError(err)
				assert.Equal(test.expectedTask, task)
			}
		})
	}
}

func BenchmarkPOSTTask(b *testing.B) {
	server := NewAPIServer(&StubTaskRepository{})

	req, _ := http.NewRequest("POST", "/tasks", strings.NewReader(`{"name":"Task 4"}`))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))

	w := httptest.NewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		server.ServeHTTP(w, req)
	}
}

func TestPUTTask(t *testing.T) {
	tasks := []model.Task{
		{ID: "1", Name: "Task 1", Completed: false},
		{ID: "2", Name: "Task 2", Completed: true},
		{ID: "3", Name: "Task 3", Completed: false},
	}

	server := NewAPIServer(&StubTaskRepository{tasks: tasks})

	tests := map[string]struct {
		id             string
		body           string
		expectedStatus int
		expectedTask   model.Task
	}{
		"Update a task": {
			id:             "1",
			body:           `{"name":"Task 1 Updated"}`,
			expectedStatus: http.StatusOK,
			expectedTask:   model.Task{ID: "1", Name: "Task 1 Updated", Completed: false},
		},
		"Update a task that does not exist": {
			id:             "4",
			body:           `{"name":"Task 4 Updated"}`,
			expectedStatus: http.StatusNotFound,
		},
		"Update a task with invalid JSON": {
			id:             "1",
			body:           `{"name":"Task 1 Updated`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert := assert.New(t)
			req, err := http.NewRequest("PUT", "/tasks/"+test.id, strings.NewReader(test.body))
			req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))
			assert.NoError(err)

			w := httptest.NewRecorder()
			server.ServeHTTP(w, req)

			assert.Equal(test.expectedStatus, w.Code)

			if test.expectedStatus == http.StatusOK {
				var task model.Task
				err = json.Unmarshal(w.Body.Bytes(), &task)
				assert.NoError(err)
				assert.Equal(test.expectedTask, task)
			}
		})
	}
}

func BenchmarkPUTTask(b *testing.B) {
	tasks := []model.Task{
		{ID: "1", Name: "Task 1", Completed: false},
	}

	server := NewAPIServer(&StubTaskRepository{tasks: tasks})

	req, _ := http.NewRequest("PUT", "/tasks/1", strings.NewReader(`{"id": "1", "name":"Task 1 Updated", "completed": false}`))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("API_KEY"))

	w := httptest.NewRecorder()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		server.ServeHTTP(w, req)
	}
}
