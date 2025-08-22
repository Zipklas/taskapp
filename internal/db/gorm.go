package db

import (
	"github.com/Zipklas/task-tracker/internal/domain"
	"gorm.io/gorm"
)

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&domain.User{},
		&domain.Task{},
	)
}
