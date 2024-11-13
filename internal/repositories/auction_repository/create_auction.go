package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
)

func (repo *auctionRepo) CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (entities.Auction, error) {
	sql, args, err := sq.Insert("auctions").
											Columns("lot_id", "status").
											Values(in.LotId, activeAuctionStatus).
											Suffix("returning id, lot_id, status").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.Auction{}, err
	}
	auction := entities.Auction{}
	err = repo.db.
		GetConn().
		QueryRow(ctx, sql, args...).
		Scan(&auction.Id, &auction.LotId, &auction.Status)

	if err != nil {
		return entities.Auction{}, err
	}

	repo.logger.Info(
		"Success create auction in storage",
		slog.Any("auctionID", auction.Id),
	)

	return auction, nil
}
