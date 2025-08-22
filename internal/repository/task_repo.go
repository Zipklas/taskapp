package repository

import (
	"github.com/Zipklas/task-tracker/internal/domain"
	"gorm.io/gorm"
)

type TaskRepository interface {
	Create(task *domain.Task) error
	FindByID(id string) (*domain.Task, error)
	FindByUserID(userID string) ([]*domain.Task, error)
	FindAll() ([]*domain.Task, error)
}

type taskRepository struct {
	db *gorm.DB
}

func NewTaskRepository(db *gorm.DB) TaskRepository {
	return &taskRepository{db: db}
}

func (r *taskRepository) Create(task *domain.Task) error {
	return r.db.Create(task).Error
}

func (r *taskRepository) FindByID(id string) (*domain.Task, error) {
	var task domain.Task
	err := r.db.First(&task, "id = ?", id).Error
	return &task, err
}

func (r *taskRepository) FindByUserID(userID string) ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.Where("user_id = ?", userID).Find(&tasks).Error
	return tasks, err
}
func (r *taskRepository) FindAll() ([]*domain.Task, error) {
	var tasks []*domain.Task
	err := r.db.Find(&tasks).Error
	return tasks, err
}
