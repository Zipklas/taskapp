package main

import (
	"log"
	"os"

	"github.com/Zipklas/task-tracker/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	// Загружаем .env файл
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		log.Fatal("POSTGRES_DSN is not set")
	}

	// Создаем и запускаем приложение
	a, err := app.NewApp(dsn)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	log.Println("Starting application...")
	a.Run()
	log.Println("Application stopped")
}
