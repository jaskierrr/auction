package repositories

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/domain/entities"
	pb "main/pkg/grpc"
	"strconv"

	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (entities.Bid, error) {
	tx, err := repo.db.GetConn().
		BeginTx(
			ctx,
			pgx.TxOptions{
				IsoLevel:   pgx.Serializable,
				AccessMode: pgx.ReadWrite,
			})

	if err != nil {
		repo.logger.Error("failed to begin transaction deposit balance: " + err.Error())
		return entities.Bid{}, err
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
		"bidder_id":  in.BidderId,
		"amount":     in.Amount,
		"status":     activeAuctionStatus,
	}
	//find starting_bid
	var startingBid float64
	var userID int64
	err = tx.
		QueryRow(ctx, findStartingBidQuery, args).
		Scan(&startingBid, &userID)
	if err != nil {
		err = errors.Join(err, errors.New("starting bid not find"))
		return entities.Bid{}, err
	}
	if startingBid > in.Amount {
		err = errors.New("not enough money")
		return entities.Bid{}, err
	}
	if userID == in.BidderId {
		err = errors.New("user cant place bid on his lot")
		return entities.Bid{}, err
	}

	//check auction status
	var status string
	err = tx.
		QueryRow(ctx, checkAuctionStatusQuery, args).
		Scan(&status)
	if err != nil || status != "Active" {
		err = errors.Join(err, errors.New("cant place a bid, auction is closed"))
		return entities.Bid{}, err
	}

	//edit user balance
	_, err = tx.Exec(ctx, editUserBalance, args)
	if err != nil {
		err = errors.Join(err, errors.New("failed edit user balance"))
		return entities.Bid{}, err
	}

	//place bid
	bid := entities.Bid{}
	var amountStr string
	err = tx.
		QueryRow(ctx, placeBidQuery, args).
		Scan(&bid.Id, &bid.AuctionId, &bid.BidderId, &amountStr, &bid.CreatedAt)

	if err != nil {
		err = errors.Join(err, errors.New("failed place bid"))
		return entities.Bid{}, err
	}
	bid.Amount, _ = strconv.ParseFloat(amountStr, 64)

	//write trancsaction
	_, err = tx.Exec(ctx, writeTransactionPlaceBidQuery, args)
	if err != nil {
		err = errors.Join(err, errors.New("cant write transaction"))
		return entities.Bid{}, err
	}

	tx.Commit(ctx)

	repo.logger.Info(
		"Success place bid in storage, user lost some money hehehe",
		slog.Any("bidID", bid.Id),
	)

	return bid, nil
}
