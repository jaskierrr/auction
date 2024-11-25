package user_service

import (
	"context"
	pb "main/pkg/grpc"
)

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
