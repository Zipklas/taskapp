package domain

import (
	"time"

	"gorm.io/gorm"
)

type Task struct {
	ID          string         `gorm:"type:uuid;default:gen_random_uuid();primaryKey"`
	Title       string         `gorm:"not null"`
	Description string         `gorm:"type:text"`
	UserID      string         `gorm:"type:uuid;not null;index"` // Внешний ключ на users
	User        User           `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}
