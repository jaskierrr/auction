package app

import (
	"context"
	"log"
	"main/internal/domain/services/repositories"
	pb "main/pkg/grpc"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type UserService struct {
	pb.UnimplementedUserServiceServer
	service repositories.UserService
}

func NewUserService(service repositories.UserService) UserService{
	return UserService{
		service: service,
	}
}

func RegisterGRPC(grpc *grpc.Server, service UserService) {
	pb.RegisterUserServiceServer(grpc, &service)
}

func (s *UserService) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (*pb.UserResponse, error) {
	user, err := s.service.CreateUser(ctx, in)

	if err != nil {
		return &pb.UserResponse{}, status.Errorf(codes.Unknown, "failed create user: %v\n", err)
	}

	log.Printf("Received: %v", user.Name)

	return &pb.UserResponse{User: &user}, nil
}

func (s *UserService) GetUser(context.Context, *pb.GetUserRequest) (*pb.UserResponse, error) {
	return nil, nil
}

func (s *UserService) DepositBalance(context.Context, *pb.DepositBalanceRequest) (*pb.BalanceResponse, error) {
	return nil, nil
}
