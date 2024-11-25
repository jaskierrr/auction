package user_service

import (
	"context"
	"log/slog"
	repositories "main/internal/repositories/user"
	pb "main/pkg/grpc"
)

type userService struct {
	repo repositories.UserRepo
	logger *slog.Logger
}

type UserService interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (pb.User, error)
	GetUser(ctx context.Context, in *pb.GetUserRequest) (pb.User, error)
	DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error)
}

func NewUserService(repo repositories.UserRepo, logger *slog.Logger) UserService {
	return &userService{
		repo: repo,
		logger: logger,
	}
}
