package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"
)



func (repo *userRepo) GetUser(ctx context.Context, in *pb.GetUserRequest) (entities.User, error) {
	sql, args, err := sq.Select("*").
											From("users").
											Where(sq.Eq{"id": in.UserId}).
											PlaceholderFormat(sq.Dollar).
											ToSql()

	if err != nil {
		return entities.User{}, err
	}

	user := entities.User{}
	err = repo.db.
		GetConn().
		QueryRow(ctx, sql, args...).
		Scan(&user.Id, &user.Name, &user.Balance)

	if err != nil {
		return entities.User{}, err
	}

	repo.logger.Info(
		"Success get user from storage",
		slog.Any("userID", user.Id),
	)

	return user, nil
}
