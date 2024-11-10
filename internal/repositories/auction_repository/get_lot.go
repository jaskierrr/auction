package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) GetLot(ctx context.Context, in *pb.GetLotRequest) (entities.Lot, error) {
	args := pgx.NamedArgs{
		"lot_id": in.LotId,
	}
	lot := entities.Lot{}
	err := repo.db.
		GetConn().
		QueryRow(ctx, getLotQuery, args).
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
