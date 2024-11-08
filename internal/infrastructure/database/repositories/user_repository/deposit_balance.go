package repositories

import (
	"context"
	"errors"
	"log/slog"
	pb "main/pkg/grpc"
	"strconv"

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

	args := pgx.NamedArgs{
		"userID": in.UserId,
		"amount": in.Amount,
	}

	var strBalance string

	err = tx.
		QueryRow(ctx, depositBalanceQuery, args).
		Scan(&strBalance)

	if err != nil {
		return &pb.BalanceResponse{}, err
	}

	balance, _ := strconv.ParseFloat(strBalance, 64)
	res := &pb.BalanceResponse{NewBalance: balance}

	_, err = tx.
		Exec(ctx, writeTransactionQuery, args)

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
