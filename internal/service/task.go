package service

import (
	"errors"
	"log"
	"time"

	"github.com/Zipklas/task-tracker/internal/domain"
	"github.com/Zipklas/task-tracker/internal/kafka"
	"github.com/Zipklas/task-tracker/internal/repository"
)

type TaskService interface {
	CreateTask(title, description, userID string) (*domain.Task, error)
	ListTasks(userID string) ([]*domain.Task, error)
	GetTaskByID(id string) (*domain.Task, error)
	UpdateTask(id, title, description string) (*domain.Task, error)
	DeleteTask(id string) error
	GetAllTasks() ([]*domain.Task, error)
}

type taskService struct {
	taskRepo  repository.TaskRepository
	kafkaProd *kafka.Producer
}

func NewTaskService(taskRepo repository.TaskRepository, kafkaProd *kafka.Producer) TaskService {
	return &taskService{
		taskRepo:  taskRepo,
		kafkaProd: kafkaProd,
	}
}

func (s *taskService) CreateTask(title, description, userID string) (*domain.Task, error) {
	task := &domain.Task{
		Title:       title,
		Description: description,
		UserID:      userID,
		CreatedAt:   time.Now(),
	}

	if err := s.taskRepo.Create(task); err != nil {
		return nil, err
	}

	event := map[string]interface{}{
		"type":    "task_created",
		"task_id": task.ID,
		"title":   task.Title,
		"user_id": task.UserID,
	}

	if err := s.kafkaProd.SendTaskCreated(event); err != nil {
		log.Printf("Failed to send Kafka message: %v", err)
	}

	return task, nil
}

func (s *taskService) ListTasks(userID string) ([]*domain.Task, error) {
	return s.taskRepo.FindByUserID(userID)
}

func (s *taskService) GetTaskByID(id string) (*domain.Task, error) {
	if id == "" {
		return nil, errors.New("task ID cannot be empty")
	}

	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (s *taskService) UpdateTask(id, title, description string) (*domain.Task, error) {
	task, err := s.taskRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	task.Title = title
	task.Description = description
	task.UpdatedAt = time.Now()

	// Здесь нужно добавить метод Update в репозиторий
	// Пока просто возвращаем обновленный объект
	return task, nil
}

func (s *taskService) DeleteTask(id string) error {
	// Здесь нужно добавить метод Delete в репозиторий
	return nil
}
func (s *taskService) GetAllTasks() ([]*domain.Task, error) {
	return s.taskRepo.FindAll()
}
