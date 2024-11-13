package repositories

import (
	"context"
	"errors"
	"log/slog"
	pb "main/pkg/grpc"
	"strconv"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *userRepo) DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error) {
	tx, err := repo.db.GetConn().
		BeginTx(
			ctx,
			pgx.TxOptions{
				IsoLevel:   pgx.Serializable,
				AccessMode: pgx.ReadWrite,
			})

	if err != nil {
		repo.logger.Error("failed to begin transaction deposit balance: " + err.Error())
		return &pb.BalanceResponse{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	strBalance, err := updateBalance(ctx, tx, in)
	if err != nil {
		return &pb.BalanceResponse{}, err
	}

	balance, _ := strconv.ParseFloat(strBalance, 64)
	res := &pb.BalanceResponse{NewBalance: balance}

	err = writeTransaction(ctx, tx, in)
	if err != nil {
		err = errors.Join(errors.New("cant write transaction"), err)
		return &pb.BalanceResponse{}, err
	}

	tx.Commit(ctx)

	repo.logger.Info(
		"Success deposite balance",
		slog.Any("userID", in.UserId),
	)

	return res, nil
}

func updateBalance(ctx context.Context, tx pgx.Tx, in *pb.DepositBalanceRequest) (string, error) {
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

func writeTransaction(ctx context.Context, tx pgx.Tx, in *pb.DepositBalanceRequest) error {
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
