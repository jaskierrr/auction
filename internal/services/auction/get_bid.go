package service

import (
	"context"
	pb "main/pkg/grpc"

	"google.golang.org/protobuf/types/known/timestamppb"
)

func (s *auctionService) GetBid(ctx context.Context, in *pb.GetBidRequest) (pb.Bid, error) {
	bid, err := s.repo.GetBid(ctx, in)
	if err != nil {
		return pb.Bid{}, err
	}

	return pb.Bid{
		Id:        bid.Id,
		AuctionId: bid.AuctionId,
		BidderId:  bid.BidderId,
		Amount:    bid.Amount,
		CreatedAt: timestamppb.New(bid.CreatedAt),
	}, nil
}
