package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) FindLotID(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest) (int64, error) {
	var lotID int64

	sql, args, err := sq.Select("l.id").
		From("auctions a").
		Join("lots l on a.lot_id = l.id").
		Where(sq.Eq{"a.id": in.AuctionId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return 0, err
	}
	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&lotID)

	return lotID, err
}
