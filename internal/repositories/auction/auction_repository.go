//go:generate mockgen -source=./auction_repository.go -destination=../../../test/mock/auction_repo_mock.go -package=mock

package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	"main/internal/infrastructure/database"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

const (
	activeLotStatus = "Active"
	closedLotStatus = "Closed"

	activeAuctionStatus = "Active"
	endedAuctionStatus  = "Ended"
)

type auctionRepo struct {
	db     database.DB
	logger *slog.Logger
}

type AuctionRepo interface {
	CreateLot(ctx context.Context, in *pb.CreateLotRequest) (entities.Lot, error)
	GetLot(ctx context.Context, in *pb.GetLotRequest) (entities.Lot, error)
	CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (entities.Auction, error)
	GetBid(ctx context.Context, in *pb.GetBidRequest) (entities.Bid, error)

	StartTx(ctx context.Context) (pgx.Tx, error)
	FindStartingBid(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (float64, int64, error)
	CheckAuctionStatus(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (string, error)
	EditUserBalance(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) error
	InsertBid(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (entities.Bid, string, error)
	PlaceBidWriteTransaction(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) error

	FindAllBids(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest) (entities.Auction, error)
	AwardingWinner(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest, auction entities.Auction, winnerBid entities.Bid) (entities.Auction, error)
	EndAuction(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest, auction entities.Auction) (entities.Auction, error)
	FindLotID(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest) (int64, error)
	CloseLot(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest, lotID int64) error
	ReturnMoney(ctx context.Context, tx pgx.Tx, bid entities.Bid) (entities.User, error)
	CloseAuctionWriteTransaction(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest, bid entities.Bid) error
}

func NewAuctionRepo(db database.DB, logger *slog.Logger) AuctionRepo {
	return &auctionRepo{
		db:     db,
		logger: logger,
	}
}
