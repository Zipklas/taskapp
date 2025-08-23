package service

import (
	"errors"
	"testing"

	"github.com/Zipklas/task-tracker/internal/domain"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockUserRepository для тестирования
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *domain.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByID(id string) (*domain.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*domain.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.User), args.Error(1)
}

func TestAuthService_Register_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Ожидаем, что репозиторий будет вызван с любым пользователем
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).Return(nil)

	user, err := service.Register("testuser", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.NotEqual(t, "password123", user.Password) // Пароль должен быть захэширован

	// Проверяем что пароль действительно захэширован
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte("password123"))
	assert.NoError(t, err)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Register_DuplicateUsername(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Ожидаем ошибку при создании пользователя
	mockRepo.On("Create", mock.AnythingOfType("*domain.User")).
		Return(errors.New("duplicate key value violates unique constraint"))

	user, err := service.Register("testuser", "password123")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "duplicate")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_Success(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Создаем захэшированный пароль
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)

	expectedUser := &domain.User{
		ID:       "user-123",
		Username: "testuser",
		Password: string(hashedPassword),
	}

	// Ожидаем поиск пользователя по username
	mockRepo.On("FindByUsername", "testuser").Return(expectedUser, nil)

	user, err := service.Login("testuser", "correctpassword")

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "user-123", user.ID)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Ожидаем что пользователь не найден
	mockRepo.On("FindByUsername", "nonexistent").Return(nil, errors.New("user not found"))

	user, err := service.Login("nonexistent", "password")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Создаем пользователя с правильным паролем
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("correctpassword"), bcrypt.DefaultCost)
	expectedUser := &domain.User{
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.On("FindByUsername", "testuser").Return(expectedUser, nil)

	user, err := service.Login("testuser", "wrongpassword")

	assert.Error(t, err)
	assert.Nil(t, user)
	assert.Contains(t, err.Error(), "invalid credentials")

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login_EmptyCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo)

	// Не ожидаем вызовов репозитория при пустых credentials
	user, err := service.Login("", "password")
	assert.Error(t, err)
	assert.Nil(t, user)

	user, err = service.Login("testuser", "")
	assert.Error(t, err)
	assert.Nil(t, user)

	mockRepo.AssertNotCalled(t, "FindByUsername")
}
