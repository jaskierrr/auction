package repositories

import (
	"context"
	"log/slog"
	"main/internal/domain/entities"
	"main/internal/infrastructure/database"
	pb "main/pkg/grpc"
	"strconv"

	"github.com/jackc/pgx/v5"
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

const (
	postUserQuery       = `insert into users (name) values (@name) returning *`
	getUserIDQuery      = `select * from users where id = @userID`
	depositBalanceQuery = `update users set balance = balance + @amount where id = @userID returning balance::Text`
)

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

	repo.logger.Info(
		"Success create user in storage",
		slog.Any("ID", user.Id),
	)

	return user, nil
}

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
		slog.Any("ID", user.Id),
	)

	return user, nil
}

func (repo *userRepo) DepositBalance(ctx context.Context, in *pb.DepositBalanceRequest) (*pb.BalanceResponse, error) {
	args := pgx.NamedArgs{
		"userID": in.UserId,
		"amount": in.Amount,
	}
	res := &pb.BalanceResponse{}

	var num string

	err := repo.db.
		GetConn().
		QueryRow(ctx, depositBalanceQuery, args).
		Scan(&num)

	if err != nil {
		repo.logger.Error("failed scan numeric: " + err.Error())
		return &pb.BalanceResponse{}, err
	}

	num2, _ := strconv.ParseFloat(num, 64)
	res.NewBalance = num2

	repo.logger.Info(
		"Success deposite balance",
		slog.Any("ID", in.UserId),
	)

	return res, nil
}
