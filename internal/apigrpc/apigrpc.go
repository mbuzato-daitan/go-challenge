package apigrpc

import (
	context "context"
	"fmt"
	"os"

	empty "github.com/golang/protobuf/ptypes/empty"
	"github.com/mtbuzato/go-challenge/internal/model"
	"google.golang.org/grpc/metadata"
)

type TaskRepository interface {
	ListAll() ([]model.Task, error)
	ListByCompletion(completed bool) ([]model.Task, error)
	Create(name string) (model.Task, error)
	GetByID(id string) (model.Task, error)
	Update(task model.Task) error
}

type grpcServer struct {
	repo TaskRepository
	UnimplementedTaskServiceServer
}

func NewGRPCServer(repo TaskRepository) *grpcServer {
	server := new(grpcServer)

	server.repo = repo

	return server
}

func (s *grpcServer) authenticate(ctx context.Context) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return fmt.Errorf("grpc.authenticate: unauthorized")
	}

	tokens := md.Get("Authorization")
	if len(tokens) == 0 {
		return fmt.Errorf("grpc.authenticate: unauthorized")
	}

	token := tokens[0]
	if token != "Bearer "+os.Getenv("API_KEY") {
		return fmt.Errorf("grpc.authenticate: unauthorized")
	}

	return nil
}

func (s *grpcServer) ListTasks(_ *empty.Empty, stream TaskService_ListTasksServer) error {
	if err := s.authenticate(stream.Context()); err != nil {
		return fmt.Errorf("grpc.ListTasks: %v", err)
	}

	tasks, err := s.repo.ListAll()
	if err != nil {
		return fmt.Errorf("grpc.ListTasks: %v", err)
	}

	for _, task := range tasks {
		if err := stream.Send(taskAtob(task)); err != nil {
			return fmt.Errorf("grpc.ListTasks: %v", err)
		}
	}

	return nil
}

func (s *grpcServer) ListTasksByCompletion(req *ListTasksByCompletionRequest, stream TaskService_ListTasksByCompletionServer) error {
	if err := s.authenticate(stream.Context()); err != nil {
		return fmt.Errorf("grpc.ListTasksByCompletion: %v", err)
	}

	tasks, err := s.repo.ListByCompletion(req.GetCompleted())
	if err != nil {
		return fmt.Errorf("grpc.ListTasksByCompletion: %v", err)
	}

	for _, task := range tasks {
		if err := stream.Send(taskAtob(task)); err != nil {
			return fmt.Errorf("grpc.ListTasksByCompletion: %v", err)
		}
	}

	return nil
}

func (s *grpcServer) GetTaskByID(ctx context.Context, req *GetTaskByIDRequest) (*Task, error) {
	if err := s.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("grpc.GetTaskByID: %v", err)
	}

	task, err := s.repo.GetByID(req.GetId())
	if err != nil {
		return nil, fmt.Errorf("grpc.GetTaskByID: %v", err)
	}

	return taskAtob(task), nil
}

func (s *grpcServer) CreateTask(ctx context.Context, req *CreateTaskRequest) (*Task, error) {
	if err := s.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("grpc.CreateTask: %v", err)
	}

	task, err := s.repo.Create(req.Name)
	if err != nil {
		return nil, fmt.Errorf("grpc.CreateTask: %v", err)
	}

	return taskAtob(task), nil
}

func (s *grpcServer) UpdateTask(ctx context.Context, task *Task) (*Task, error) {
	if err := s.authenticate(ctx); err != nil {
		return nil, fmt.Errorf("grpc.UpdateTask: %v", err)
	}

	t := taskBtoa(task)
	err := s.repo.Update(t)
	if err != nil {
		return nil, fmt.Errorf("grpc.UpdateTask: %v", err)
	}

	return taskAtob(t), nil
}
