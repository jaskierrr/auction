package repositories

import (
	"context"
	"errors"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"
	"strconv"

	sq "github.com/Masterminds/squirrel"
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

	//find all bids on this auction
	auction, err := findAllBids(ctx, repo, tx, in)

	if err != nil {
		err = errors.Join(errors.New("failed find all bids on this auction: "), err)
		return entities.Auction{}, err
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
	auction, err = awardingWinner(ctx, tx, in, auction, winnerBid)

	if err != nil {
		err = errors.Join(errors.New("cant award winner"), err)
		return entities.Auction{}, err
	}
	if auction.Status == endedAuctionStatus {
		err = errors.New("cant end auction it is already over")
		return entities.Auction{}, err
	}

	// end auction
	auction, err = endAuction(ctx, tx, in, auction)

	if err != nil {
		err = errors.Join(errors.New("failed end auction and award the winner: "), err)
		return entities.Auction{}, err
	}

	// find lot id
	lotID, err := findLotID(ctx, tx, in)

	if err != nil {
		err = errors.Join(errors.New("cant find lot id: "), err)
		return entities.Auction{}, err
	}

	// close lot
	err = closeLot(ctx, tx, in, lotID)

	if err != nil {
		repo.logger.Error(
			"Error",
			slog.Any("cant close lot", err),
		)
	}

	//! bids without winner
	if len(auction.Bids) > 1 {
		bidsWithoutWinner := append(auction.Bids[:winnerIndx], auction.Bids[winnerIndx+1:]...)

		repo.logger.Info(
			"arr bis without winner",
			slog.Any("without WINNER", bidsWithoutWinner),
		)
		//return money money money and write transaction
		for _, v := range bidsWithoutWinner {
			user, err := returnMoney(ctx, tx, v)

			repo.logger.Info(
				"Return money to users",
				slog.Any("userID", user.Id),
				slog.Any("Name", user.Name),
				slog.Any("old balance", user.Balance-v.Amount),
				slog.Any("new balance", user.Balance),
			)

			if err != nil {
				return entities.Auction{}, err
			}

			err = closeAuctionWriteTransaction(ctx, tx, in, v)

			if err != nil {
				err = errors.Join(err, errors.New("cant write transaction"))
				return entities.Auction{}, err
			}
			repo.logger.Info("Trans write!", slog.Any("user", v.BidderId))
	}
	} else {
		repo.logger.Info(
			"it was one bidder, auction closed without refund",
		)
	}

	repo.logger.Error(
		"Error",
		slog.Any("err", err),
	)

	tx.Commit(ctx)

	repo.logger.Info(
		"Success close auction in storage",
		slog.Any("auctionID", auction.Id),
	)

	return auction, nil
}

// getBidsQuery = `select * from bids where auction_id = @auction_id`


func findAllBids(ctx context.Context, repo *auctionRepo, tx pgx.Tx, in *pb.CloseAuctionRequest) (entities.Auction, error){
	sql, args, err := sq.Select("*").
											From("bids").
											Where(sq.Eq{"auction_id": in.AuctionId}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.Auction{}, err
	}

	bidsRow, err := tx.
	Query(ctx, sql, args...)

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

	repo.logger.Info(
		"auction",
		slog.Any("auction",auction),
	)
	return auction, err
}

// addWinnerAuctionQuery = `update auctions set winner_id = @winner_id where id = @auction_id returning *`

func awardingWinner(ctx context.Context,tx pgx.Tx, in *pb.CloseAuctionRequest, auction entities.Auction, winnerBid entities.Bid) (entities.Auction, error){
	sql, args, err := sq.Update("auctions").
											Set("winner_id", winnerBid.BidderId).
											Where(sq.Eq{"id": in.AuctionId}).
											Suffix("returning *").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.Auction{}, err
	}

	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&auction.Id, &auction.LotId, &auction.Status, &auction.WinnerId)

	return auction, err
}

// endAuctionQuery = `update auctions set status = @status where id = @auction_idreturning *`

func endAuction(ctx context.Context,tx pgx.Tx, in *pb.CloseAuctionRequest, auction entities.Auction) (entities.Auction, error){
	sql, args, err := sq.Update("auctions").
											Set("status", endedAuctionStatus).
											Where(sq.Eq{"id": in.AuctionId}).
											Suffix("returning *").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.Auction{}, err
	}

	err = tx.
	QueryRow(ctx, sql, args...).
	Scan(&auction.Id, &auction.LotId, &auction.Status, &auction.WinnerId)

	return auction, err
}

// findLotId = `select l.id from auctions a join lots l on a.lot_id = l.id where a.id = @auction_id and a.status = @status`

func findLotID(ctx context.Context,tx pgx.Tx, in *pb.CloseAuctionRequest) (int64, error){
	var lotID int64

	sql, args, err := sq.Select("l.id").
											From("auctions a").
											Join("lots l on a.lot_id = l.id").
											Where(sq.Eq{"a.id": in.AuctionId}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return 0, err
	}
	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&lotID)

	return lotID, err
}

//closedLotQuery = `update lots set status = 'Closed' where id = $1`

func closeLot(ctx context.Context,tx pgx.Tx, in *pb.CloseAuctionRequest, lotID int64) error{
	sql, _, err := sq.Select("l.id").
											From("auctions a").
											Join("lots l on a.lot_id = l.id").
											Where(sq.Eq{"a.id": in.AuctionId, "a.status": closedLotStatus}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, lotID)


	return err
}

// returnMoney = `update users set balance = balance + @amount where id = @bidder_id returning *`

func returnMoney(ctx context.Context,tx pgx.Tx, bid entities.Bid) (entities.User, error){
	user := entities.User{}

	sql, args, err := sq.Update("users").
											Set("balance", sq.Expr("balance + ?", bid.Amount)).
											Where(sq.Eq{"id": bid.BidderId}).
											Suffix("returning *").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.User{}, err
	}
	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&user.Id, &user.Name, &user.Balance)


	return user, err
}

// writeTransactionEndedAuctionQuery = `insert into transactions (sender_id, sender_type, recipient_id, recipient_type, amount, transaction_type) values (@auction_id, 'Auction', @bidder_id, 'User', @amount, 'Refund')`

func closeAuctionWriteTransaction(ctx context.Context,tx pgx.Tx, in *pb.CloseAuctionRequest, bid entities.Bid) error{
	sql, args, err := sq.Insert("transactions").
											Columns("sender_id", "sender_type", "recipient_id", "recipient_type", "amount", "transaction_type").
											Values(in.AuctionId, "Auction", bid.BidderId, "User", bid.Amount, "Refund").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)

	return err
}
