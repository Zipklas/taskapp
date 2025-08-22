package main

import (
	"log"
	"net/http"
	"os"

	"github.com/Zipklas/task-tracker/internal/kafka"
	"github.com/Zipklas/task-tracker/internal/websocket"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	wsServer := websocket.NewServer()

	kafkaConsumer, err := kafka.NewConsumer(
		[]string{"localhost:9092"},
		"task-events",
	)
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	kafkaConsumer.AddHandler(func(message []byte) {
		wsServer.BroadcastMessage(map[string]interface{}{
			"Type": "task_event",
			"data": string(message),
		})
	})

	go kafkaConsumer.Start()
	go wsServer.Run()

	http.HandleFunc("/ws", wsServer.HandleWebSocket)

	port := os.Getenv("WS_PORT")
	if port == "" {
		port = "8081"
	}

	log.Printf("Websocket server starting on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
