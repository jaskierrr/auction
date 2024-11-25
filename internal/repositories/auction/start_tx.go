package repositories

import (
	"context"

	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) StartTx(ctx context.Context) (pgx.Tx, error) {
	tx, err := repo.db.GetConn().
		BeginTx(
			ctx,
			pgx.TxOptions{
				IsoLevel:   pgx.Serializable,
				AccessMode: pgx.ReadWrite,
			})

	if err != nil {
		return nil, err
	}

	return tx, nil
}
