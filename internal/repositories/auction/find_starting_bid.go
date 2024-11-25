package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) FindStartingBid(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (float64, int64, error) {
	var startingBid float64
	var userID int64

	sql, args, err := sq.Select("l.starting_bid", "l.seller_id").
		From("auctions a").
		Join("lots l on a.lot_id = l.id").
		Where(sq.Eq{"a.id": in.AuctionId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return -1, -1, err
	}

	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&startingBid, &userID)

	return startingBid, userID, err
}
