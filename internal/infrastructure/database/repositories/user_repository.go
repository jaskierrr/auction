package repositories

import (
	"context"
	"main/internal/domain/entities"
	"main/internal/infrastructure/database"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

type userRepo struct {
	db database.DB
	// logger *slog.Logger
}

type UserRepo interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (entities.User, error)
}

func NewUserRepo(db database.DB) UserRepo {
	return &userRepo{
		db: db,
	}
}

const postUserQuery = `insert into users (name) values (@name) returning *`

func (repo *userRepo) CreateUser(ctx context.Context, in *pb.CreateUserRequest) (entities.User, error) {
	args := pgx.NamedArgs{
		"name": in.Name,
	}
	user := entities.User{}
	err := repo.db.
		GetConn().
		QueryRow(ctx, postUserQuery, args).
		Scan(&user.Id, &user.Name, &user.Balance)

	if err != nil {
		return entities.User{}, err
	}

	// repo.logger.Info(
	// 	"Success POST user in storage",
	// 	slog.Any("ID", user.Id),
	// )

	return user, nil
}
