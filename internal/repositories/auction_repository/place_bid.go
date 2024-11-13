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

  // find starting_bid
	startingBid, userID, err := findStartingBid(ctx, tx, in)

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

	// check auction status
	status, err := checkAuctionStatus(ctx, tx, in)

	if err != nil || status != "Active" {
		err = errors.Join(err, errors.New("cant place a bid, auction is closed"))
		return entities.Bid{}, err
	}

	// edit user balance
	err = editUserBalance(ctx, tx, in)

	if err != nil {
		err = errors.Join(err, errors.New("failed edit user balance"))
		return entities.Bid{}, err
	}

	// place bid
	bid, amountStr, err := placeBid(ctx, tx, in)

	if err != nil {
		err = errors.Join(err, errors.New("failed place bid"))
		return entities.Bid{}, err
	}
	bid.Amount, _ = strconv.ParseFloat(amountStr, 64)

	// write trancsaction
	err = PlaceBidwriteTransaction(ctx, tx, in)

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

func findStartingBid(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (float64, int64, error){
	var startingBid float64
	var userID int64

	sql, args, err := sq.Select("l.starting_bid", "l.seller_id").
											From("auctions a").
											Join("lots l on a.lot_id = l.id").
											Where(sq.Eq{"a.id": in.AuctionId}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return -1, -1, err
	}

	err = tx.
				QueryRow(ctx, sql, args...).
				Scan(&startingBid, &userID)

	return startingBid, userID, err
}

func checkAuctionStatus(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (string, error) {
	var status string

	sql, args, err := sq.Select("status").
											From("auctions").
											Where(sq.Eq{"id": in.AuctionId}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return "", err
	}

	err = tx.
	QueryRow(ctx, sql, args...).
	Scan(&status)

	return status, err
}

func editUserBalance(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) error {
	sql, args, err := sq.Update("users").
											Set("balance", sq.Expr("balance - ?", in.Amount)).
											Where(sq.Eq{"id": in.BidderId}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)

	return err
}

func placeBid(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) (entities.Bid, string, error)  {
	bid := entities.Bid{}
	amountStr := ""

	sql, args, err := sq.Insert("bids").
											Columns("auction_id", "bidder_id", "amount").
											Values(in.AuctionId, in.BidderId, in.Amount).
											Suffix("on conflict (auction_id, bidder_id) do update set amount = bids.amount + EXCLUDED.amount RETURNING *").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return bid, amountStr, err
	}
	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&bid.Id, &bid.AuctionId, &bid.BidderId, &amountStr, &bid.CreatedAt)

	return bid, amountStr, err
}

func PlaceBidwriteTransaction(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) error  {
	sql, args, err := sq.Insert("transactions").
											Columns("sender_id", "sender_type", "recipient_id", "recipient_type", "amount", "transaction_type").
											Values(in.BidderId, "User", in.AuctionId, "Auction", in.Amount, "Payment").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, sql, args...)

	return err
}
