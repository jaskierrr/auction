package user_service

import (
	"context"
	"errors"
	pb "main/pkg/grpc"
	"strconv"
)

func (s *userService) DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error) {
	if in.Amount <= 0 {
		err := errors.New("amount must be greater than 0")
		return &pb.BalanceResponse{}, err
	}

	// start tx
	tx, err := s.repo.StartTx(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction deposit balance: " + err.Error())
		return &pb.BalanceResponse{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		}
	}()

	// update balance
	strBalance, err := s.repo.UpdateBalance(ctx, tx, in)
	if err != nil {
		return nil, err
	}

	balanceValue, _ := strconv.ParseFloat(strBalance, 64)
	balance := &pb.BalanceResponse{NewBalance: balanceValue}

	// write transaction
	err = s.repo.PlaceBidWriteTransaction(ctx, tx, in)
	if err != nil {
		err = errors.Join(errors.New("cant write transaction"), err)
		return nil, err
	}

	tx.Commit(ctx)

	return balance, nil
}
