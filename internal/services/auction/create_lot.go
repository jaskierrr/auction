package service

import (
	"context"
	pb "main/pkg/grpc"
)

func (s *auctionService) CreateLot(ctx context.Context, in *pb.CreateLotRequest) (pb.Lot, error) {
	lot, err := s.repo.CreateLot(ctx, in)
	if err != nil {
		return pb.Lot{}, err
	}

	return pb.Lot{
		Id:          lot.Id,
		Title:       lot.Title,
		Description: lot.Description,
		StartingBid: lot.StartingBid,
		SellerId:    lot.SellerId,
		Status:      lot.Status,
	}, nil
}
