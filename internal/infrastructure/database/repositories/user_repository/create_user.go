package repositories

import (
	"context"
	"log/slog"
	"main/internal/domain/entities"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

func (repo *userRepo) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (entities.User, error) {
	args := pgx.NamedArgs{
		"name": in.Name,
	}
	user := entities.User{}
	err := repo.db.
		GetConn().
		QueryRow(ctx, createUserQuery, args).
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
