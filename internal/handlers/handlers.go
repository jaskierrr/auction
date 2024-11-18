package handlers

import (
	"log/slog"
	service "main/internal/services"
	pb "main/pkg/grpc"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type UserHandlers struct {
	pb.UnimplementedUserServiceServer
	logger  *slog.Logger
	validator *validator.Validate
	service service.UserService
}

type AuctionHandlers struct {
	pb.UnimplementedAuctionServiceServer
	logger  *slog.Logger
	validator *validator.Validate
	service service.AuctionService
}

func NewUserHandlers(service service.UserService, logger *slog.Logger, validator *validator.Validate) UserHandlers {
	return UserHandlers{
		logger:  logger,
		validator:  validator,
		service: service,
	}
}
func NewAuctionHandlers(service service.AuctionService, logger *slog.Logger, validator *validator.Validate) AuctionHandlers {
	return AuctionHandlers{
		logger:  logger,
		validator:  validator,
		service: service,
	}
}

func RegisterGRPC(grpc *grpc.Server, userUC UserHandlers, auctionUC AuctionHandlers) {
	pb.RegisterUserServiceServer(grpc, &userUC)
	pb.RegisterAuctionServiceServer(grpc, &auctionUC)
}
