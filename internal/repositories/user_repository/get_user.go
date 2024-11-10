package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

func (repo *userRepo) GetUser(ctx context.Context, in *pb.GetUserRequest) (entities.User, error) {
	args := pgx.NamedArgs{
		"userID": in.UserId,
	}
	user := entities.User{}
	err := repo.db.
		GetConn().
		QueryRow(ctx, getUserIDQuery, args).
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
