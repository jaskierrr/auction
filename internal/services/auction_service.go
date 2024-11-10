package service

import (
	"context"
	"main/internal/repositories/auction_repository"
	pb "main/pkg/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

type auctionService struct {
	repo repositories.AuctionRepo
}

type AuctionService interface {
	CreateLot(ctx context.Context, in *pb.CreateLotRequest) (pb.Lot, error)
	GetLot(ctx context.Context, in *pb.GetLotRequest) (pb.Lot, error)

	CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (pb.Auction, error)
	CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (pb.Auction, error)

	PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (pb.Bid, error)
	GetBid(ctx context.Context, in *pb.GetBidRequest) (pb.Bid, error)
}

func NewAuctionService(repo repositories.AuctionRepo) AuctionService {
	return &auctionService{
		repo: repo,
	}
}

func (s *auctionService) CreateLot(ctx context.Context, in *pb.CreateLotRequest) (pb.Lot, error) {
	lot, err := s.repo.CreateLot(ctx, in)
	if err != nil {
		return pb.Lot{}, err
	}

	return pb.Lot{
		Id:      lot.Id,
		Title: lot.Title,
		Description: lot.Description,
		StartingBid: lot.StartingBid,
		SellerId: lot.SellerId,
		Status: lot.Status,
	}, nil
}

func (s *auctionService) GetLot(ctx context.Context, in *pb.GetLotRequest) (pb.Lot, error) {
	lot, err := s.repo.GetLot(ctx, in)
	if err != nil {
		return pb.Lot{}, err
	}

	return pb.Lot{
		Id:      lot.Id,
		Title: lot.Title,
		Description: lot.Description,
		StartingBid: lot.StartingBid,
		SellerId: lot.SellerId,
		Status: lot.Status,
	}, nil
}

func (s *auctionService) CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (pb.Auction, error) {
	auction, err := s.repo.CreateAuction(ctx, in)
	if err != nil {
		return pb.Auction{}, err
	}

	return pb.Auction{
		Id:      auction.Id,
		LotId: auction.LotId,
		Status: auction.Status,
		WinnerId: auction.WinnerId,
	}, nil
}

func (s *auctionService) CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (pb.Auction, error) {
	auction, err := s.repo.CloseAuction(ctx, in)
	if err != nil {
		return pb.Auction{}, err
	}

	return pb.Auction{
		Id:      auction.Id,
		LotId: auction.LotId,
		Status: auction.Status,
		WinnerId: auction.WinnerId,
	}, nil
}

func (s *auctionService) PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (pb.Bid, error) {
	bid, err := s.repo.PlaceBid(ctx, in)
	if err != nil {
		return pb.Bid{}, err
	}

	return pb.Bid{
		Id: bid.Id,
		AuctionId: bid.AuctionId,
		BidderId: bid.BidderId,
		Amount: bid.Amount,
		CreatedAt: timestamppb.New(bid.CreatedAt),
	}, nil
}

func (s *auctionService) GetBid(ctx context.Context, in *pb.GetBidRequest) (pb.Bid, error) {
	bid, err := s.repo.GetBid(ctx, in)
	if err != nil {
		return pb.Bid{}, err
	}

	return pb.Bid{
		Id: bid.Id,
		AuctionId: bid.AuctionId,
		BidderId: bid.BidderId,
		Amount: bid.Amount,
		CreatedAt: timestamppb.New(bid.CreatedAt),
	}, nil
}
