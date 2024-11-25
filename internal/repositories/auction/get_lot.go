package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
)

func (repo *auctionRepo) GetLot(ctx context.Context, in *pb.GetLotRequest) (entities.Lot, error) {
	sql, args, err := sq.Select("*").
										From("lots").
										Where(sq.Eq{"id": in.LotId}).
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
		"Success get lot from storage",
		slog.Any("lotID", lot.Id),
	)

	return lot, nil
}
