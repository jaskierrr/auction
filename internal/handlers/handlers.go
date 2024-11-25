package handlers

import (
	"log/slog"
	user_service "main/internal/services/user"
	auction_service "main/internal/services/auction"
	pb "main/pkg/grpc"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
)

type UserHandlers struct {
	pb.UnimplementedUserServiceServer
	logger    *slog.Logger
	validator *validator.Validate
	service   user_service.UserService
}

type AuctionHandlers struct {
	pb.UnimplementedAuctionServiceServer
	logger    *slog.Logger
	validator *validator.Validate
	service   auction_service.AuctionService
}

func NewUserHandlers(service user_service.UserService, logger *slog.Logger, validator *validator.Validate) UserHandlers {
	return UserHandlers{
		logger:    logger,
		validator: validator,
		service:   service,
	}
}
func NewAuctionHandlers(service auction_service.AuctionService, logger *slog.Logger, validator *validator.Validate) AuctionHandlers {
	return AuctionHandlers{
		logger:    logger,
		validator: validator,
		service:   service,
	}
}

func RegisterGRPC(grpc *grpc.Server, userUC UserHandlers, auctionUC AuctionHandlers) {
	pb.RegisterUserServiceServer(grpc, &userUC)
	pb.RegisterAuctionServiceServer(grpc, &auctionUC)
}
