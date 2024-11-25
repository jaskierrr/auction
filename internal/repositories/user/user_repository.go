//go:generate mockgen -source=./user_repository.go -destination=../../../test/mock/user_repo_mock.go -package=mock

package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	"main/internal/infrastructure/database"
	pb "main/pkg/grpc"

	"github.com/jackc/pgx/v5"
)

type userRepo struct {
	db     database.DB
	logger *slog.Logger
}

type UserRepo interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (entities.User, error)
	GetUser(ctx context.Context, in *pb.GetUserRequest) (entities.User, error)
	PlaceBidWriteTransaction(ctx context.Context, tx pgx.Tx, in *pb.DepositBalanceRequest) error
	UpdateBalance(ctx context.Context, tx pgx.Tx, in *pb.DepositBalanceRequest) (string, error)

	StartTx(ctx context.Context) (pgx.Tx, error)
}

func NewUserRepo(db database.DB, logger *slog.Logger) UserRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}
