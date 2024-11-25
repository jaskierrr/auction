package service

import (
	"context"
	"log/slog"
	repositories "main/internal/repositories/auction"
	pb "main/pkg/grpc"
)

const (
	activeAuctionStatus = "Active"
	endedAuctionStatus  = "Ended"
)

type auctionService struct {
	repo   repositories.AuctionRepo
	logger *slog.Logger
}

type AuctionService interface {
	CreateLot(ctx context.Context, in *pb.CreateLotRequest) (pb.Lot, error)
	GetLot(ctx context.Context, in *pb.GetLotRequest) (pb.Lot, error)

	CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (pb.Auction, error)
	CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (pb.Auction, error)

	PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (pb.Bid, error)
	GetBid(ctx context.Context, in *pb.GetBidRequest) (pb.Bid, error)
}

func NewAuctionService(repo repositories.AuctionRepo, logger *slog.Logger) AuctionService {
	return &auctionService{
		repo:   repo,
		logger: logger,
	}
}
