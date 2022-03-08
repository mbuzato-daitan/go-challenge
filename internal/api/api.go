package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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

type APIServer struct {
	http.Handler
	repo TaskRepository
}

// Creates a new API server with the given repository.
func NewAPIServer(repo TaskRepository) *APIServer {
	server := new(APIServer)

	server.repo = repo

	router := mux.NewRouter()

	tasksRouter := router.PathPrefix("/tasks").Subrouter()

	tasksRouter.HandleFunc("", server.getTasks).Methods("GET").Queries("{completed:(true)|(false)}")
	tasksRouter.HandleFunc("", server.getTasks).Methods("GET")
	tasksRouter.HandleFunc("/{id}", server.getTask).Methods("GET")
	tasksRouter.HandleFunc("", server.postTask).Methods("POST")
	tasksRouter.HandleFunc("/{id}", server.putTask).Methods("PUT")

	server.Handler = router

	return server
}

func (s *APIServer) handleError(w http.ResponseWriter, err error) {
	if errors.IsExternal(err) {
		if errors.IsNotFound(err) {
			w.WriteHeader(http.StatusNotFound)
		} else {
			w.WriteHeader(http.StatusBadRequest)
		}

		fmt.Fprint(w, "{\"error\": \"", err.Error(), "\"}")
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "{\"error\": \"An unknown error ocurred.\"}")
		fmt.Printf("API Error: %s\n", err.Error())
	}
}

func (s *APIServer) getTasks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	completed, ok := vars["completed"]

	var tasks []model.Task
	var err error

	if !ok {
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

func (s *APIServer) getTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

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

func (s *APIServer) postTask(w http.ResponseWriter, r *http.Request) {
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

func (s *APIServer) putTask(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	id := vars["id"]

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
