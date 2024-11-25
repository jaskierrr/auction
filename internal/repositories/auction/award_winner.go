package repositories

import (
	"context"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) AwardingWinner(ctx context.Context, tx pgx.Tx, in *pb.CloseAuctionRequest, auction entities.Auction, winnerBid entities.Bid) (entities.Auction, error) {
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
