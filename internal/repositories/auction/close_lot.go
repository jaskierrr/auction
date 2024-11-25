package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) CloseLot(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest, lotID int64) error {
	sql, args, err := sq.Update("lots").
											Set("status", closedLotStatus).
											Where(sq.Eq{"id": lotID}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)

	return err
}
