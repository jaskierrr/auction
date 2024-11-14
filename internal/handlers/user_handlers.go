package handlers

import (
	"context"
	pb "main/pkg/grpc"

	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *UserHandlers) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := s.service.CreateUser(ctx, in)

	if err != nil {
		s.logger.Error("failed create user: " + err.Error())
		return &pb.UserResponse{}, status.Errorf(codes.Unknown, "failed create user: %v\n", err)
	}

	return &pb.UserResponse{User: &user}, nil
}

func (s *UserHandlers) GetUser(ctx context.Context, in *pb.GetUserRequest) (*pb.UserResponse, error) {
	user, err := s.service.GetUser(ctx, in)

	if err != nil {
		s.logger.Error("failed get user: " + err.Error())
		return &pb.UserResponse{
			User: &user,
		}, status.Errorf(codes.Unknown, "failed get user: %v", err)
	}

	return &pb.UserResponse{User: &user}, nil
}

func (s *UserHandlers) DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error) {
	balance, err := s.service.DepositBalance(ctx, in)

	if err != nil {
		s.logger.Error("failed deposite balance: " + err.Error())
		return &pb.BalanceResponse{}, status.Errorf(codes.Unknown, "failed deposite balance: %v\n", err)
	}

	return balance, nil
}
