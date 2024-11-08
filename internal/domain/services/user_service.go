package service

import (
	"context"
	"main/internal/infrastructure/database/repositories/user_repository"
	pb "main/pkg/grpc"
)

type userService struct {
	repo repositories.UserRepo
}

type UserService interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (pb.User, error)
	GetUser(ctx context.Context, in *pb.GetUserRequest) (pb.User, error)
	DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error)
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

func (s *userService) GetUser(ctx context.Context, in *pb.GetUserRequest) (pb.User, error) {
	user, err := s.repo.GetUser(ctx, in)
	if err != nil {
		return pb.User{}, err
	}

	return pb.User{
		Id:      user.Id,
		Name:    user.Name,
		Balance: user.Balance,
	}, nil
}

func (s *userService) DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error) {
	balance, err := s.repo.DepositBalance(ctx, in)
	if err != nil {
		return nil, err
	}

	return balance, nil
}
