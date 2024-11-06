package app

import (
	"log/slog"
	service "main/internal/domain/services"
	pb "main/pkg/grpc"

	"google.golang.org/grpc"
)

type UserUsecase struct {
	pb.UnimplementedUserServiceServer
	logger  *slog.Logger
	service service.UserService
}

type AuctionUsecase struct {
	pb.UnimplementedAuctionServiceServer
	logger  *slog.Logger
	service service.AuctionService
}

func NewUserUsecase(service service.UserService, logger *slog.Logger) UserUsecase {
	return UserUsecase{
		logger:  logger,
		service: service,
	}
}
func NewAuctionUsecase(service service.AuctionService, logger *slog.Logger) AuctionUsecase {
	return AuctionUsecase{
		logger:  logger,
		service: service,
	}
}

func RegisterGRPC(grpc *grpc.Server, userUC UserUsecase, auctionUC AuctionUsecase) {
	pb.RegisterUserServiceServer(grpc, &userUC)
	pb.RegisterAuctionServiceServer(grpc, &auctionUC)
}
