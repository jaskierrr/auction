package repositories

import (
	"context"
	"main/internal/entities"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

func (repo *auctionRepo) ReturnMoney(ctx context.Context, tx pgx.Tx, bid entities.Bid) (entities.User, error) {
	user := entities.User{}

	sql, args, err := sq.Update("users").
		Set("balance", sq.Expr("balance + ?", bid.Amount)).
		Where(sq.Eq{"id": bid.BidderId}).
		Suffix("returning *").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	if err != nil {
		return entities.User{}, err
	}
	err = tx.
		QueryRow(ctx, sql, args...).
		Scan(&user.Id, &user.Name, &user.Balance)

	return user, err
}
