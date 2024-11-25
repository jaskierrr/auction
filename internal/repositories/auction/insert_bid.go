package repositories

import (
	"context"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) InsertBid(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (entities.Bid, string, error) {
	bid := entities.Bid{}
	amountStr := ""

	sql, args, err := sq.Insert("bids").
		Columns("auction_id", "bidder_id", "amount").
		Values(in.AuctionId, in.BidderId, in.Amount).
		Suffix("on conflict (auction_id, bidder_id) do update set amount = bids.amount + EXCLUDED.amount RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return bid, amountStr, err
	}
	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&bid.Id, &bid.AuctionId, &bid.BidderId, &amountStr, &bid.CreatedAt)

	return bid, amountStr, err
}
