package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) GetBid(ctx context.Context, in *pb.GetBidRequest) (entities.Bid, error) {
	args := pgx.NamedArgs{
		"auction_id": in.AuctionId,
	}

	var amountStr string
	bid := entities.Bid{}
	err := repo.db.
		GetConn().
		QueryRow(ctx, getBidQuery, args).
		Scan(&bid.Id, &bid.AuctionId, &bid.BidderId, &amountStr, &bid.CreatedAt)

	if err != nil {
		return entities.Bid{}, err
	}

	bid.Amount, _ = strconv.ParseFloat(amountStr, 64)

	repo.logger.Info(
		"Success get bid from storage",
		slog.Any("bidID", bid.Id),
	)

	return bid, nil
}
