package repositories

import (
	"context"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) PlaceBidWriteTransaction(ctx context.Context, tx pgx.Tx, in *pb.PlaceBidRequest) error {
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
