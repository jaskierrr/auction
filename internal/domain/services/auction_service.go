package service

import (
	"main/internal/infrastructure/database/repositories"
)

type auctionService struct {
	repo repositories.AuctionRepo
}

type AuctionService interface {
}

func NewAuctionService(repo repositories.AuctionRepo) AuctionService {
	return &auctionService{
		repo: repo,
	}
}
