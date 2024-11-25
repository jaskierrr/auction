package handlers

import (
	"context"
	pb "main/pkg/grpc"

	"google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

func (s *AuctionHandlers) CreateLot(ctx context.Context, in *pb.CreateLotRequest) (*pb.LotResponse, error) {
	lot, err := s.service.CreateLot(ctx, in)

	if err != nil {
		s.logger.Error("failed create lot: " + err.Error())
		return &pb.LotResponse{}, status.Errorf(codes.Unknown, "failed create lot: %v", err)
	}

	return &pb.LotResponse{Lot: &lot}, nil
}

func (s *AuctionHandlers) GetLot(ctx context.Context, in *pb.GetLotRequest) (*pb.LotResponse, error) {
	lot, err := s.service.GetLot(ctx, in)

	if err != nil {
		s.logger.Error("failed get lot: " + err.Error())
		return &pb.LotResponse{
			Lot: &lot,
		}, status.Errorf(codes.Unknown, "failed get lot: %v", err)
	}

	return &pb.LotResponse{Lot: &lot}, nil
}

func (s *AuctionHandlers) CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (*pb.AuctionResponse, error) {
	auction, err := s.service.CreateAuction(ctx, in)

	if err != nil {
		s.logger.Error("failed create auction: " + err.Error())
		return &pb.AuctionResponse{}, status.Errorf(codes.Unknown, "failed create auction: %v", err)
	}

	return &pb.AuctionResponse{Auction: &auction}, nil
}

func (s *AuctionHandlers) CloseAuction(ctx context.Context, in *pb.CloseAuctionRequest) (*pb.AuctionResponse, error) {
	auction, err := s.service.CloseAuction(ctx, in)

	if err != nil {
		s.logger.Error("failed close auction: " + err.Error())
		return &pb.AuctionResponse{}, status.Errorf(codes.Unknown, "failed close auction: %v", err)
	}

	return &pb.AuctionResponse{Auction: &auction}, nil
}

func (s *AuctionHandlers) PlaceBid(ctx context.Context, in *pb.PlaceBidRequest) (*pb.BidResponse, error) {
	bid, err := s.service.PlaceBid(ctx, in)

	if err != nil {
		s.logger.Error("failed place bid: " + err.Error())
		return &pb.BidResponse{
			Bid: &bid,
		}, status.Errorf(codes.Unknown, "failed place bid: %v", err)
	}

	return &pb.BidResponse{Bid: &bid}, nil
}

func (s *AuctionHandlers) GetBid(ctx context.Context, in *pb.GetBidRequest) (*pb.BidResponse, error) {
	bid, err := s.service.GetBid(ctx, in)

	if err != nil {
		s.logger.Error("failed get bid: " + err.Error())
		return &pb.BidResponse{
			Bid: &bid,
		}, status.Errorf(codes.Unknown, "failed get bid: %v", err)
	}

	return &pb.BidResponse{Bid: &bid}, nil
}
