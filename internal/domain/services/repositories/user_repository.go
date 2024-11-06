package repositories

import (
	"context"
	"main/internal/infrastructure/database/repositories"
	pb "main/pkg/grpc"
)

type userService struct {
	// logger *slog.Logger
	repo repositories.UserRepo
}

type UserService interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (pb.User, error)
}

func NewUserService(repo repositories.UserRepo) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (pb.User, error) {
	user, err := s.repo.CreateUser(ctx, in)
	if err != nil {
		return pb.User{}, err
	}

	return pb.User{
		Id:      user.Id,
		Name:    user.Name,
		Balance: user.Balance,
	}, nil
}
