package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/fullstorydev/grpcui/standalone"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Подключаемся к gRPC серверу
	grpcAddr := "localhost:50051"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Используем NewClient вместо Dial
	conn, err := grpc.NewClient(grpcAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Правильный вызов HandlerViaReflection
	handler, err := standalone.HandlerViaReflection(
		ctx,
		conn,
		grpcAddr,
	)
	if err != nil {
		log.Fatalf("Failed to create grpc-ui handler: %v", err)
	}

	// Запускаем HTTP сервер
	uiPort := ":8080"
	log.Printf("gRPC UI available at: http://localhost%s", uiPort)

	http.Handle("/", handler)
	if err := http.ListenAndServe(uiPort, nil); err != nil {
		log.Fatalf("HTTP server failed: %v", err)
	}
}
