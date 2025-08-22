package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Zipklas/task-tracker/internal/db"
	"github.com/Zipklas/task-tracker/internal/grpcserver"
	"github.com/Zipklas/task-tracker/internal/kafka"
	"github.com/Zipklas/task-tracker/internal/repository"
	"github.com/Zipklas/task-tracker/internal/service"
	"gorm.io/gorm"
)

type App struct {
	db         *gorm.DB
	grpcServer *grpcserver.Server
}

func NewApp(dsn string) (*App, error) {
	// Инициализация БД
	database, err := db.NewPostgres(dsn)
	if err != nil {
		return nil, err
	}

	// Автомиграция
	if err := db.AutoMigrate(database); err != nil {
		return nil, err
	}

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(database)
	taskRepo := repository.NewTaskRepository(database)

	kafkaProducer, err := kafka.NewProducer(
		[]string{"localhost:9092"},
		"task-events",
	)
	if err != nil {
		return nil, err
	}

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo)
	taskService := service.NewTaskService(taskRepo, kafkaProducer)

	// Инициализация gRPC сервера
	grpcSrv := grpcserver.NewServer(authService, taskService)

	return &App{
		db:         database,
		grpcServer: grpcSrv,
	}, nil
}

func (a *App) Run() error {
	// Запуск gRPC сервера
	go func() {
		if err := a.grpcServer.Run("50051"); err != nil {
			log.Fatalf("failed to run gRPC server: %v", err)
		}
	}()

	log.Println("gRPC server started on port 50051")

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down application...")
	return a.GracefulShutdown()
}

func (a *App) GracefulShutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Остановка gRPC сервера
	if err := a.grpcServer.Stop(ctx); err != nil {
		return err
	}

	// Закрытие соединения с БД
	sqlDB, err := a.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
