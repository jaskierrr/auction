package repositories

import (
	"context"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) CloseAuctionWriteTransaction(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest, bid entities.Bid) error {
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
