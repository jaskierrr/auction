//go:generate mockgen -source=./user_repository.go -destination=../../../../test/mock/user_repo_mock.go -package=mock

package repositories

import (
	"context"
	"log/slog"
	"main/internal/entities"
	"main/internal/infrastructure/database"
	pb "main/pkg/grpc"
)

type userRepo struct {
	db     database.DB
	logger *slog.Logger
}

type UserRepo interface {
	CreateUser(ctx context.Context, in *pb.CreateUserRequest) (entities.User, error)
	GetUser(ctx context.Context, in *pb.GetUserRequest) (entities.User, error)
	DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error)
}

func NewUserRepo(db database.DB, logger *slog.Logger) UserRepo {
	return &userRepo{
		db:     db,
		logger: logger,
	}
}
