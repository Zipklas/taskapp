package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/Zipklas/task-tracker/internal/service"
	authpb "github.com/Zipklas/task-tracker/pkg/protobuf/auth"
	taskpb "github.com/Zipklas/task-tracker/pkg/protobuf/task"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	authService service.AuthService
	taskService service.TaskService
	grpcServer  *grpc.Server
}

func NewServer(authService service.AuthService, taskService service.TaskService) *Server {
	return &Server{
		authService: authService,
		taskService: taskService,
		grpcServer:  grpc.NewServer(),
	}
}

func (s *Server) Run(port string) error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}

	// Регистрируем сервисы
	authpb.RegisterAuthServiceServer(s.grpcServer, &authHandler{authService: s.authService})
	taskpb.RegisterTaskServiceServer(s.grpcServer, &taskHandler{taskService: s.taskService})

	// Для разработки - рефлексия
	reflection.Register(s.grpcServer)

	log.Printf("gRPC server listening on port %s", port)
	return s.grpcServer.Serve(lis)
}

func (s *Server) Stop(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-ctx.Done():
		s.grpcServer.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}
