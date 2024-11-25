package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *userRepo) UpdateBalance(ctx context.Context, tx pgx.Tx, in *pb.DepositBalanceRequest) (string, error) {
	sql, args, err := sq.Update("users").
		Set("balance", sq.Expr("balance + ?", in.Amount)).
		Where(sq.Eq{"id": in.UserId}).
		Suffix("returning balance::Text").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return "", err
	}

	var strBalance string
	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&strBalance)

	return strBalance, err
}
