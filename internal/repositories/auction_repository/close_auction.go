package repositories

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (entities.Auction, error) {
	tx, err := repo.db.GetConn().
		BeginTx(
			ctx,
			pgx.TxOptions{
				IsoLevel:   pgx.Serializable,
				AccessMode: pgx.ReadWrite,
			})

	if err != nil {
		repo.logger.Error("failed to begin transaction deposit balance: " + err.Error())
		return entities.Auction{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
			repo.logger.Error(
				"Error",
				slog.Any("rollback transaction", err),
			)
		}
	}()

	args := pgx.NamedArgs{
		"auction_id": in.AuctionId,
		"status":     closedAuctionStatus,
	}

	//find all bids on this auction
	bidsRow, err := tx.
		Query(ctx, getBidsQuery, args)

	if err != nil {
		repo.logger.Error("failed find all bids on this auction: " + err.Error())
		return entities.Auction{}, err
	}

	auction := entities.Auction{
		Bids: make([]entities.Bid, 0),
	}
	var amountStr string

	for bidsRow.Next() {
		amountStr = ""
		bid := entities.Bid{}
		err := bidsRow.Scan(&bid.Id, &bid.AuctionId, &bid.BidderId, &amountStr, &bid.CreatedAt)
		if err != nil {
			repo.logger.Error(
				"Error",
				slog.Any("cant go through the entire bids array", err),
			)
			return entities.Auction{}, err
		}
		bid.Amount, _ = strconv.ParseFloat(amountStr, 64)
		auction.Bids = append(auction.Bids, bid)
	}
	bidsRow.Close()

	var winnerBid entities.Bid
	var winnerIndx int

	for i, v := range auction.Bids {
		if v.Amount > winnerBid.Amount {
			winnerBid = v
			winnerIndx = i
		}

		argsForTrans := pgx.NamedArgs{
			"auction_id": in.AuctionId,
			"bidder_id": v.BidderId,
			"amount": v.Amount,
		}
		_, err = tx.Exec(ctx, writeTransactionEndedAuctionQuery, argsForTrans)
		if err != nil {
			err = errors.Join(err, errors.New("cant write transaction"))
			return entities.Auction{}, err
		}
		repo.logger.Info("Trans write!", slog.Any("user", v.BidderId))
	}

	//end auction and awarding the winner
	args["winner_id"] = winnerBid.BidderId

	err = tx.
		QueryRow(ctx, endedAuctionQuery, args).
		Scan(&auction.Id, &auction.LotId, &auction.Status, &auction.WinnerId)

	if auction.Status == closedAuctionStatus {
		err = errors.New("cant end auction it is already over")
		return entities.Auction{}, err
	}

	if err != nil {
		repo.logger.Error("failed end auction and award the winner: " + err.Error())
		return entities.Auction{}, err
	}

	//find lot id
	var lotID int64
	err = tx.
		QueryRow(ctx, findLotId, args).
		Scan(&lotID)

	if err != nil {
		repo.logger.Error(
			"Error",
			slog.Any("cant find lot id", err),
		)
	}

	//close lot
	_, err = tx.Exec(ctx, closedLotQuery, lotID)

	if err != nil {
		repo.logger.Error(
			"Error",
			slog.Any("cant close lot", err),
		)
	}

	//bids without winner
	bidsWithoutWinner := append(auction.Bids[:winnerIndx], auction.Bids[winnerIndx+1:]...)

	//return money money money
	for _, v := range bidsWithoutWinner {
		returnMoneyArgs := pgx.NamedArgs{
			"bidder_id": v.BidderId,
			"amount":    v.Amount,
		}

		user := entities.User{}
		err = tx.
			QueryRow(ctx, returnMoney, returnMoneyArgs).
			Scan(&user.Id, &user.Name, &user.Balance)

		repo.logger.Info(
			"Return money to users",
			slog.Any("userID", user.Id),
			slog.Any("old balance", user.Balance-v.Amount),
			slog.Any("new balance", user.Balance),
		)
	}

	if err != nil {
		return entities.Auction{}, err
	}

	tx.Commit(ctx)

	repo.logger.Info(
		"Success close auction in storage",
		slog.Any("auctionID", auction.Id),
	)

	return auction, nil
}
