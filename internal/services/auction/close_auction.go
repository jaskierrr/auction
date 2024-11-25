package service

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"
)

func (s *auctionService) CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (pb.Auction, error) {
	tx, err := s.repo.StartTx(ctx)

	if err != nil {
		s.logger.Error("failed to begin transaction deposit balance: " + err.Error())
		return pb.Auction{}, err
	}

	//find all bids on this auction
	auction, err := s.repo.FindAllBids(ctx, tx, in)

	if err != nil {
		err = errors.Join(errors.New("failed find all bids on this auction: "), err)
		return pb.Auction{}, err
	}

	var winnerBid entities.Bid
	var winnerIndx int

	for i, v := range auction.Bids {
		if v.Amount > winnerBid.Amount {
			winnerBid = v
			winnerIndx = i
		}
	}

	//awarding the winner
	auction, err = s.repo.AwardingWinner(ctx, tx, in, auction, winnerBid)

	if err != nil {
		err = errors.Join(errors.New("cant award winner"), err)
		return pb.Auction{}, err
	}
	if auction.Status == endedAuctionStatus {
		err = errors.New("cant end auction it is already over")
		return pb.Auction{}, err
	}

	// end auction
	auction, err = s.repo.EndAuction(ctx, tx, in, auction)

	if err != nil {
		err = errors.Join(errors.New("failed end auction and award the winner: "), err)
		return pb.Auction{}, err
	}

	// find lot id
	lotID, err := s.repo.FindLotID(ctx, tx, in)

	if err != nil {
		err = errors.Join(errors.New("cant find lot id: "), err)
		return pb.Auction{}, err
	}

	// close lot
	err = s.repo.CloseLot(ctx, tx, in, lotID)

	if err != nil {
		s.logger.Error(
			"Error",
			slog.Any("cant close lot", err),
		)
	}

	//! bids without winner
	if len(auction.Bids) > 1 {
		bidsWithoutWinner := append(auction.Bids[:winnerIndx], auction.Bids[winnerIndx+1:]...)

		//return money money money and write transaction
		for _, v := range bidsWithoutWinner {
			user, err := s.repo.ReturnMoney(ctx, tx, v)

			s.logger.Info(
				"Return money to user",
				slog.Any("userID", user.Id),
				slog.Any("Name", user.Name),
				slog.Any("old balance", user.Balance-v.Amount),
				slog.Any("new balance", user.Balance),
			)

			if err != nil {
				return pb.Auction{}, err
			}

			err = s.repo.CloseAuctionWriteTransaction(ctx, tx, in, v)

			if err != nil {
				err = errors.Join(err, errors.New("cant write transaction"))
				return pb.Auction{}, err
			}
		}
	} else {
		s.logger.Info(
			"it was one bidder, auction closed without refund",
		)
	}

	tx.Commit(ctx)

	s.logger.Info(
		"Success close auction in storage",
		slog.Any("auctionID", auction.Id),
		slog.Any("winnerID", auction.WinnerId),
	)

	return pb.Auction{
		Id:       auction.Id,
		LotId:    auction.LotId,
		Status:   auction.Status,
		WinnerId: auction.WinnerId,
	}, nil
}
