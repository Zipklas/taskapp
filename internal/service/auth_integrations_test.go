//go:build integration
// +build integration

package service

import (
	"testing"

	"github.com/Zipklas/task-tracker/internal/db"
	"github.com/Zipklas/task-tracker/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestAuthService_Integration_RegisterAndLogin(t *testing.T) {
	// Пропускаем если запущено с -short
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Используем тестовую БД
	dsn := "postgresql://postgres:1234@localhost:5433/taskapp_test?sslmode=disable"
	database, err := db.NewPostgres(dsn)
	assert.NoError(t, err)

	// Выполняем миграции перед тестом
	err = db.AutoMigrate(database)
	assert.NoError(t, err)

	// Очищаем БД перед тестом
	database.Exec("DELETE FROM users")

	// Создаем реальные зависимости
	userRepo := repository.NewUserRepository(database)
	authService := NewAuthService(userRepo)

	// Test Register
	user, err := authService.Register("integrationuser", "integrationpass")
	assert.NoError(t, err)
	assert.Equal(t, "integrationuser", user.Username)

	// Test Login
	loggedInUser, err := authService.Login("integrationuser", "integrationpass")
	assert.NoError(t, err)
	assert.Equal(t, user.ID, loggedInUser.ID)
	assert.Equal(t, user.Username, loggedInUser.Username)

	// Test Login with wrong password
	_, err = authService.Login("integrationuser", "wrongpassword")
	assert.Error(t, err)
	assert.Equal(t, "invalid credentials", err.Error())
}

func TestAuthService_Integration_DuplicateUsername(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	dsn := "postgresql://postgres:1234@localhost:5433/taskapp_test?sslmode=disable"
	database, err := db.NewPostgres(dsn)
	assert.NoError(t, err)

	// Миграции
	err = db.AutoMigrate(database)
	assert.NoError(t, err)

	database.Exec("DELETE FROM users")

	userRepo := repository.NewUserRepository(database)
	authService := NewAuthService(userRepo)

	// Первый пользователь - успешно
	_, err = authService.Register("duplicateuser", "password")
	assert.NoError(t, err)

	// Второй пользователь с тем же username - ошибка
	_, err = authService.Register("duplicateuser", "password")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "duplicate")
}
