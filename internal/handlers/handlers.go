package handlers

import (
	"log/slog"
	service "main/internal/services"
	pb "main/pkg/grpc"

	"google.golang.org/grpc"
)

type UserHandlers struct {
	pb.UnimplementedUserServiceServer
	logger  *slog.Logger
	service service.UserService
}

type AuctionHandlers struct {
	pb.UnimplementedAuctionServiceServer
	logger  *slog.Logger
	service service.AuctionService
}

func NewUserHandlers(service service.UserService, logger *slog.Logger) UserHandlers {
	return UserHandlers{
		logger:  logger,
		service: service,
	}
}
func NewAuctionHandlers(service service.AuctionService, logger *slog.Logger) AuctionHandlers {
	return AuctionHandlers{
		logger:  logger,
		service: service,
	}
}

func RegisterGRPC(grpc *grpc.Server, userUC UserHandlers, auctionUC AuctionHandlers) {
	pb.RegisterUserServiceServer(grpc, &userUC)
	pb.RegisterAuctionServiceServer(grpc, &auctionUC)
}
