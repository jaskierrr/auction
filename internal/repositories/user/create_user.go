package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	sq "github.com/Masterminds/squirrel"

)

func (repo *userRepo) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (entities.User, error) {
	sql, args, err := sq.Insert("users").
											Columns("name").
											Values(in.Name).
											Suffix("returning *").
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
		"Success create user in storage",
		slog.Any("userID", user.Id),
	)

	return user, nil
}
