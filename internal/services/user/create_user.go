package user_service

import (
	"context"
	"errors"
	pb "main/pkg/grpc"
)

func (s *userService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (pb.User, error) {
	if in.Name == "" {
		err := errors.New("name user is empty")
		return pb.User{}, err
	}

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
