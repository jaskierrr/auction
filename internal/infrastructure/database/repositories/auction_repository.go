package repositories

import (
	"log/slog"
	"main/internal/infrastructure/database"
)

type auctionRepo struct {
	db     database.DB
	logger *slog.Logger
}

type AuctionRepo interface {
}

func NewAuctionRepo(db database.DB, logger *slog.Logger) AuctionRepo {
	return &auctionRepo{
		db:     db,
		logger: logger,
	}
}
