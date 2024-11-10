package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) CreateLot(ctx context.Context, in *pb.CreateLotRequest) (entities.Lot, error) {
	args := pgx.NamedArgs{
		"title":        in.Title,
		"description":  in.Description,
		"starting_bid": in.StartingBid,
		"seller_id":    in.SellerId,
		"status":       activeLotStatus,
	}
	lot := entities.Lot{}
	err := repo.db.
		GetConn().
		QueryRow(ctx, createLotQuery, args).
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
