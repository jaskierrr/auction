package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
)

func (repo *auctionRepo) CreateLot(ctx context.Context, in *pb.CreateLotRequest) (entities.Lot, error) {
	sql, args, err := sq.Insert("lots").
											Columns("title", "description", "starting_bid", "seller_id", "status").
											Values(in.Title, in.Description, in.StartingBid, in.SellerId, activeLotStatus).
											Suffix("returning *").
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.Lot{}, err
	}

	lot := entities.Lot{}
	err = repo.db.
		GetConn().
		QueryRow(ctx, sql, args...).
		Scan(&lot.Id, &lot.Title, &lot.Description, &lot.StartingBid, &lot.SellerId, &lot.Status)

	if err != nil {
		return entities.Lot{}, err
	}

	repo.logger.Info(
		"Success create lot in storage",
		slog.Any("lotID", lot.Id),
	)

	return lot, nil
}
