package service

import (
	"context"
	"errors"
	"log/slog"
	pb "main/pkg/grpc"
	"strconv"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *auctionService) PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (pb.Bid, error) {
	// start tx
	tx, err := s.repo.StartTx(ctx)
	if err != nil {
		s.logger.Error("failed to begin transaction to place bid: " + err.Error())
		return pb.Bid{}, err
	}

	defer func() {
		if err != nil {
			s.logger.Error("failed to write transaction to place bid: " + err.Error())
			_ = tx.Rollback(ctx)
		}
	}()

	// find starting_bid
	startingBid, userID, err := s.repo.FindStartingBid(ctx, tx, in)

	if err != nil {
		err = errors.Join(err, errors.New("starting bid not find"))
		return pb.Bid{}, err
	}
	if startingBid > in.Amount {
		err = errors.New("not enough money")
		return pb.Bid{}, err
	}
	if userID == in.BidderId {
		err = errors.New("user cant place bid on his lot")
		return pb.Bid{}, err
	}

	// check auction status
	status, err := s.repo.CheckAuctionStatus(ctx, tx, in)

	if err != nil || status != "Active" {
		err = errors.Join(err, errors.New("cant place a bid, auction is closed"))
		return pb.Bid{}, err
	}

	// edit user balance
	err = s.repo.EditUserBalance(ctx, tx, in)

	if err != nil {
		err = errors.Join(err, errors.New("failed edit user balance"))
		return pb.Bid{}, err
	}

	// insert bid
	bid, amountStr, err := s.repo.InsertBid(ctx, tx, in)

	if err != nil {
		err = errors.Join(err, errors.New("failed place bid"))
		return pb.Bid{}, err
	}
	bid.Amount, _ = strconv.ParseFloat(amountStr, 64)

	// write trancsaction
	err = s.repo.PlaceBidWriteTransaction(ctx, tx, in)

	if err != nil {
		err = errors.Join(err, errors.New("cant write transaction"))
		return pb.Bid{}, err
	}

	s.logger.Info(
		"Success place bid in storage, user lost some money hehehe",
		slog.Any("bidID", bid.Id),
	)

	tx.Commit(ctx)

	return pb.Bid{
		Id:        bid.Id,
		AuctionId: bid.AuctionId,
		BidderId:  bid.BidderId,
		Amount:    bid.Amount,
		CreatedAt: timestamppb.New(bid.CreatedAt),
	}, nil
}
