package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) FindAllBids(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest) (entities.Auction, error) {
	sql, args, err := sq.Select("*").
		From("bids").
		Where(sq.Eq{"auction_id": in.AuctionId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return entities.Auction{}, err
	}

	bidsRow, err := tx.
		Query(ctx, sql, args...)

	auction := entities.Auction{
		Bids: make([]entities.Bid, 0),
	}
	var amountStr string

	for bidsRow.Next() {
		amountStr = ""
		bid := entities.Bid{}
		err := bidsRow.Scan(&bid.Id, &bid.AuctionId, &bid.BidderId, &amountStr, &bid.CreatedAt)
		if err != nil {
			repo.logger.Error(
				"Error",
				slog.Any("cant go through the entire bids array", err),
			)
			return entities.Auction{}, err
		}
		bid.Amount, _ = strconv.ParseFloat(amountStr, 64)
		auction.Bids = append(auction.Bids, bid)
	}
	bidsRow.Close()
	
	return auction, err
}
