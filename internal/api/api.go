package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/mtbuzato/go-challenge/internal/errors"
	"github.com/mtbuzato/go-challenge/internal/model"
)

type TaskRepository interface {
	ListAll() ([]model.Task, error)
	ListByCompletion(completed bool) ([]model.Task, error)
	Create(name string) (model.Task, error)
	GetByID(id string) (model.Task, error)
	Update(task model.Task) error
}

type apiServer struct {
	http.Handler
	repo TaskRepository
}

// Creates a new API server with the given repository.
func NewAPIServer(repo TaskRepository) *apiServer {
	server := new(apiServer)

	server.repo = repo

	router := http.NewServeMux()

	router.HandleFunc("/tasks", server.handleTasks)
	router.HandleFunc("/tasks/", server.handleTask)

	router.HandleFunc("/", server.handleNotFound)

	server.Handler = server.mdwHeaders(server.mdwAuthentication(router))

	return server
}

func (s *apiServer) handleError(w http.ResponseWriter, err error) {
	if errors.IsExternal(err) {
		code := errors.GetHTTPStatusCode(err)

		if code != 0 {
			w.WriteHeader(code)
		} else {
			if strings.Contains(strings.ToLower(err.Error()), "not found") {
				w.WriteHeader(http.StatusNotFound)
			} else {
				w.WriteHeader(http.StatusBadRequest)
			}
		}

		fmt.Fprint(w, "{\"error\": \"", err.Error(), "\"}")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "{\"error\": \"An unknown error ocurred.\"}")
		fmt.Printf("API Error: %s\n", err.Error())
	}
}

func (s *apiServer) handleNotFound(w http.ResponseWriter, r *http.Request) {
	s.handleError(w, errors.NewHTTPError("Endpoint not found.", http.StatusNotFound))
}

func (s *apiServer) handleTasks(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.getTasks(w, r)
	case "POST":
		s.postTask(w, r)
	default:
		s.handleNotFound(w, r)
	}
}

func (s *apiServer) handleTask(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "/")

	if len(split) != 3 {
		s.handleNotFound(w, r)
		return
	}

	taskId := split[2]

	switch r.Method {
	case "GET":
		s.getTask(w, r, taskId)
	case "PUT":
		s.putTask(w, r, taskId)
	default:
		s.handleNotFound(w, r)
	}
}

func (s *apiServer) getTasks(w http.ResponseWriter, r *http.Request) {
	completed := r.URL.Query().Get("completed")

	var tasks []model.Task
	var err error

	if completed == "" {
		tasks, err = s.repo.ListAll()
	} else {
		tasks, err = s.repo.ListByCompletion(completed == "true")
	}

	if err != nil {
		s.handleError(w, err)
		return
	}

	str, err := json.Marshal(tasks)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(str)
}

func (s *apiServer) getTask(w http.ResponseWriter, r *http.Request, id string) {
	task, err := s.repo.GetByID(id)
	if err != nil {
		s.handleError(w, err)
		return
	}

	str, err := json.Marshal(task)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(str)
}

type PostTaskBody struct {
	Name string `json:"name"`
}

func (s *apiServer) postTask(w http.ResponseWriter, r *http.Request) {
	var taskBody PostTaskBody

	err := json.NewDecoder(r.Body).Decode(&taskBody)
	if err != nil {
		s.handleError(w, errors.NewExternalError("Invalid body."))
		return
	}

	task, err := s.repo.Create(taskBody.Name)
	if err != nil {
		s.handleError(w, err)
		return
	}

	str, err := json.Marshal(task)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(str)
}

type PutTaskBody struct {
	Name      string `json:"name"`
	Completed bool   `json:"completed"`
}

func (s *apiServer) putTask(w http.ResponseWriter, r *http.Request, id string) {
	var taskBody PutTaskBody

	err := json.NewDecoder(r.Body).Decode(&taskBody)
	if err != nil {
		s.handleError(w, errors.NewExternalError("Invalid body."))
		return
	}

	task := model.Task{
		ID:        id,
		Name:      taskBody.Name,
		Completed: taskBody.Completed,
	}

	err = s.repo.Update(task)
	if err != nil {
		s.handleError(w, err)
		return
	}

	str, err := json.Marshal(task)
	if err != nil {
		s.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(str)
}
