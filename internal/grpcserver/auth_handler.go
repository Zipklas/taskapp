package grpcserver

import (
	"context"

	"github.com/Zipklas/task-tracker/internal/service"
	authpb "github.com/Zipklas/task-tracker/pkg/protobuf/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authHandler struct {
	authpb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func (h *authHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user, err := h.authService.Register(req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to register user: %v", err)
	}

	return &authpb.RegisterResponse{
		Id:       user.ID,
		Username: user.Username,
	}, nil
}

func (h *authHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials: %v", err)
	}

	// Здесь должен быть генератор JWT токенов
	// Пока возвращаем заглушку
	token := "jwt-token-placeholder-" + user.ID

	return &authpb.LoginResponse{
		Token: token,
	}, nil
}
