package app

import (
	"context"
	pb "main/pkg/grpc"

	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *AuctionUsecase) CreateLot(ctx context.Context, in *pb.CreateLotRequest) (*pb.LotResponse, error) {
	lot, err := s.service.CreateLot(ctx, in)

	if err != nil {
		s.logger.Error("failed create lot: " + err.Error())
		return &pb.LotResponse{}, status.Errorf(codes.Unknown, "failed create lot: %v\n", err)
	}

	return &pb.LotResponse{Lot: &lot}, nil
}

func (s *AuctionUsecase) GetLot(ctx context.Context, in *pb.GetLotRequest) (*pb.LotResponse, error) {
	lot, err := s.service.GetLot(ctx, in)

	if err != nil {
		s.logger.Error("failed get lot: " + err.Error())
		return &pb.LotResponse{}, status.Errorf(codes.Unknown, "failed get lot: %v\n", err)
	}

	return &pb.LotResponse{Lot: &lot}, nil
}

func (s *AuctionUsecase) CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (*pb.AuctionResponse, error) {
	auction, err := s.service.CreateAuction(ctx, in)

	if err != nil {
		s.logger.Error("failed create auction: " + err.Error())
		return &pb.AuctionResponse{}, status.Errorf(codes.Unknown, "failed create auction: %v\n", err)
	}

	return &pb.AuctionResponse{Auction: &auction}, nil
}

func (s *AuctionUsecase) CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (*pb.AuctionResponse, error) {
	auction, err := s.service.CloseAuction(ctx, in)

	if err != nil {
		s.logger.Error("failed close auction: " + err.Error())
		return &pb.AuctionResponse{}, status.Errorf(codes.Unknown, "failed close auction: %v\n", err)
	}

	return &pb.AuctionResponse{Auction: &auction}, nil
}

func (s *AuctionUsecase) PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (*pb.BidResponse, error) {
	bid, err := s.service.PlaceBid(ctx, in)

	if err != nil {
		s.logger.Error("failed place bid: " + err.Error())
		return &pb.BidResponse{}, status.Errorf(codes.Unknown, "failed place bid: %v\n", err)
	}

	return &pb.BidResponse{Bid: &bid}, nil
}

func (s *AuctionUsecase) GetBid(ctx context.Context, in *pb.GetBidRequest) (*pb.BidResponse, error) {
	bid, err := s.service.GetBid(ctx, in)

	if err != nil {
		s.logger.Error("failed get bid: " + err.Error())
		return &pb.BidResponse{}, status.Errorf(codes.Unknown, "failed place bid: %v\n", err)
	}

	return &pb.BidResponse{Bid: &bid}, nil
}
