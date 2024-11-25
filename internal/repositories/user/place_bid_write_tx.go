package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *userRepo) PlaceBidWriteTransaction(ctx context.Context, tx pgx.Tx, in *pb.DepositBalanceRequest) error {
	sql, args, err := sq.Insert("transactions").
												Columns("recipient_id", "recipient_type", "amount", "transaction_type").
												Values(in.UserId, "User", in.Amount, "Deposit").
												PlaceholderFormat(sq.Dollar).
												ToSql()
	if err != nil {
		return err
	}

	_, err = tx.
		Exec(ctx, sql, args...)

	return err
}
