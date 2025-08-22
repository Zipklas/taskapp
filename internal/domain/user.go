package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Username  string         `gorm:"uniqueIndex;not null"`
	Password  string         `gorm:"not null"` // Хранится хэш (bcrypt)
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"` // soft delete
}
