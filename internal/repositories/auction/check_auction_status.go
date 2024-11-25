package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) CheckAuctionStatus(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (string, error) {
	var status string

	sql, args, err := sq.Select("status").
		From("auctions").
		Where(sq.Eq{"id": in.AuctionId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return "", err
	}

	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&status)

	return status, err
}
