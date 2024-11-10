package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (entities.Auction, error) {
	args := pgx.NamedArgs{
		"lot_id": in.LotId,
		"status": activeAuctionStatus,
	}
	auction := entities.Auction{}
	err := repo.db.
		GetConn().
		QueryRow(ctx, createAuctionQuery, args).
		Scan(&auction.Id, &auction.LotId, &auction.Status)

	if err != nil {
		return entities.Auction{}, err
	}

	repo.logger.Info(
		"Success create auction in storage",
		slog.Any("auctionID", auction.Id),
	)

	return auction, nil
}
