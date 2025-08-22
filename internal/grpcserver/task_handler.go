package grpcserver

import (
	"context"

	"github.com/Zipklas/task-tracker/internal/service"
	taskpb "github.com/Zipklas/task-tracker/pkg/protobuf/task"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type taskHandler struct {
	taskpb.UnimplementedTaskServiceServer
	taskService service.TaskService
}

func (h *taskHandler) CreateTask(ctx context.Context, req *taskpb.CreateTaskRequest) (*taskpb.CreateTaskResponse, error) {
	createdTask, err := h.taskService.CreateTask(req.Title, req.Description, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create task: %v", err)
	}

	return &taskpb.CreateTaskResponse{
		Id:          createdTask.ID,
		Title:       createdTask.Title,
		Description: createdTask.Description,
		UserId:      createdTask.UserID,
		CreatedAt:   createdTask.CreatedAt.Unix(),
	}, nil
}

func (h *taskHandler) ListTasks(ctx context.Context, req *taskpb.ListTasksRequest) (*taskpb.ListTasksResponse, error) {
	tasks, err := h.taskService.ListTasks(req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list tasks: %v", err)
	}

	var taskProtos []*taskpb.Task
	for _, t := range tasks {
		taskProtos = append(taskProtos, &taskpb.Task{
			Id:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			UserId:      t.UserID,
			CreatedAt:   t.CreatedAt.Unix(),
		})
	}

	return &taskpb.ListTasksResponse{
		Tasks: taskProtos,
	}, nil
}

// Дополнительные методы которые можно добавить в proto
func (h *taskHandler) GetTask(ctx context.Context, req *taskpb.GetTaskRequest) (*taskpb.Task, error) {
	t, err := h.taskService.GetTaskByID(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "task not found: %v", err)
	}

	return &taskpb.Task{
		Id:          t.ID,
		Title:       t.Title,
		Description: t.Description,
		UserId:      t.UserID,
		CreatedAt:   t.CreatedAt.Unix(),
	}, nil
}

func (h *taskHandler) UpdateTask(ctx context.Context, req *taskpb.UpdateTaskRequest) (*taskpb.Task, error) {
	updatedTask, err := h.taskService.UpdateTask(req.Id, req.Title, req.Description)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update task: %v", err)
	}

	return &taskpb.Task{
		Id:          updatedTask.ID,
		Title:       updatedTask.Title,
		Description: updatedTask.Description,
		UserId:      updatedTask.UserID,
		CreatedAt:   updatedTask.CreatedAt.Unix(),
	}, nil
}

func (h *taskHandler) DeleteTask(ctx context.Context, req *taskpb.DeleteTaskRequest) (*taskpb.DeleteTaskResponse, error) {
	if err := h.taskService.DeleteTask(req.Id); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete task: %v", err)
	}

	return &taskpb.DeleteTaskResponse{Success: true}, nil
}
func (h *taskHandler) GetAllTasks(ctx context.Context, req *taskpb.GetAllTasksRequest) (*taskpb.GetAllTasksResponse, error) {
	tasks, err := h.taskService.GetAllTasks()
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get all tasks: %v", err)
	}

	var taskProtos []*taskpb.Task
	for _, t := range tasks {
		taskProtos = append(taskProtos, &taskpb.Task{
			Id:          t.ID,
			Title:       t.Title,
			Description: t.Description,
			UserId:      t.UserID,
			CreatedAt:   t.CreatedAt.Unix(),
			UpdatedAt:   t.UpdatedAt.Unix(),
		})
	}

	return &taskpb.GetAllTasksResponse{
		Tasks: taskProtos,
	}, nil
}
