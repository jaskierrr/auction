package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	"main/internal/infrastructure/database"
	pb "main/pkg/grpc"
)

const (
	activeLotStatus = "Active"
	closedLotStatus = "Closed"

	activeAuctionStatus = "Active"
	endedAuctionStatus = "Ended"
)

type auctionRepo struct {
	db     database.DB
	logger *slog.Logger
}

type AuctionRepo interface {
	CreateLot(ctx context.Context, in *pb.CreateLotRequest) (entities.Lot, error)
	GetLot(ctx context.Context, in *pb.GetLotRequest) (entities.Lot, error)

	CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (entities.Auction, error)
	CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (entities.Auction, error)

	PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (entities.Bid, error)
	GetBid(ctx context.Context, in *pb.GetBidRequest) (entities.Bid, error)
}

func NewAuctionRepo(db database.DB, logger *slog.Logger) AuctionRepo {
	return &auctionRepo{
		db:     db,
		logger: logger,
	}
}
