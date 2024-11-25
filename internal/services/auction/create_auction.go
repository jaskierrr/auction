package service

import (
	"context"
	pb "main/pkg/grpc"
)

func (s *auctionService) CreateAuction(ctx context.Context, in *pb.CreateAuctionRequest) (pb.Auction, error) {
	auction, err := s.repo.CreateAuction(ctx, in)
	if err != nil {
		return pb.Auction{}, err
	}

	return pb.Auction{
		Id:       auction.Id,
		LotId:    auction.LotId,
		Status:   auction.Status,
		WinnerId: auction.WinnerId,
	}, nil
}
