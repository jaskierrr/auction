package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"
	"strconv"

	sq "github.com/Masterminds/squirrel"
)


func (repo *auctionRepo) GetBid(ctx context.Context, in *pb.GetBidRequest) (entities.Bid, error) {
	sql, args, err := sq.Select("*").
											From("bids").
											Where(sq.Eq{"id": in.Id}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.Bid{}, err
	}

	var amountStr string
	bid := entities.Bid{}
	err = repo.db.
		GetConn().
		QueryRow(ctx, sql, args...).
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
