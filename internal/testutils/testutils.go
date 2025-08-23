package testutils

import (
	"os"
	"testing"

	"github.com/Zipklas/task-tracker/internal/db"
	"gorm.io/gorm"
)

// SkipIfShort пропускает интеграционные тесты при флаге -short
func SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
}

// GetTestDB возвращает подключение к тестовой БД
func GetTestDB(t *testing.T) *gorm.DB {
	SkipIfShort(t)

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		dsn = "postgresql://postgres:1234@localhost:5432/taskapp_test?sslmode=disable"
	}

	database, err := db.NewPostgres(dsn)
	if err != nil {
		t.Fatalf("Failed to connect to test database: %v", err)
	}

	// Очищаем БД перед тестом
	database.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;")

	return database
}
