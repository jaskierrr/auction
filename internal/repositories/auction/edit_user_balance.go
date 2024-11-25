package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) EditUserBalance(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) error {
	sql, args, err := sq.Update("users").
		Set("balance", sq.Expr("balance - ?", in.Amount)).
		Where(sq.Eq{"id": in.BidderId}).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)

	return err
}
