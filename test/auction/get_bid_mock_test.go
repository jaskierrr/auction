package test

import (
	"context"
	"errors"
	"main/internal/entities"
	"main/internal/handlers"
	service "main/internal/services"
	pb "main/pkg/grpc"
	"main/pkg/logger"
	"main/test/mock"
	"reflect"
	"testing"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func Test_Get_Bid(t *testing.T) {
	t.Parallel()

	type fields struct {
		auctionRepo *mock.MockAuctionRepo
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewLogger()
	validator := validator.New(validator.WithRequiredStructEnabled())

	auctionRepoMock := mock.NewMockAuctionRepo(ctrl)

	testFields := &fields{
		auctionRepo: auctionRepoMock,
	}

	service := service.NewAuctionService(auctionRepoMock)
	handlers := handlers.NewAuctionHandlers(service, logger, validator)

	time := time.Now()
	timePB := timestamppb.New(time)

	bid := &pb.Bid{
		Id: 1,
		AuctionId: 1,
		BidderId: 1,
		Amount: 100,
		CreatedAt: timePB,
	}

	reqArgDef := &pb.GetBidRequest{
		Id: 1,
	}

	reqArgErr := &pb.GetBidRequest{
		Id: -1,
	}

	lotRepoResponse := entities.Bid{
		Id: 1,
		AuctionId: 1,
		BidderId: 1,
		Amount: 100,
		CreatedAt: time,
	}

	lotErr, errErr := &pb.BidResponse{Bid: &pb.Bid{}}, status.Errorf(codes.Unknown, "failed get bid: %v", errors.New("no rows in result set"))

	ctx := context.Background()

	tests := []struct {
		name        string
		args        *pb.GetBidRequest
		prepare     func(f *fields)
		wantResLot *pb.BidResponse
		wantResErr  error
	}{
		{
			name: "valid",
			args: reqArgDef,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auctionRepo.EXPECT().GetBid(ctx, reqArgDef).Return(lotRepoResponse, nil),
				)
			},
			wantResLot: &pb.BidResponse{Bid: bid},
			wantResErr:  nil,
		},
		{
			name: "wrong_ID",
			args: reqArgErr,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auctionRepo.EXPECT().GetBid(ctx, reqArgErr).Return(entities.Bid{}, errors.New("no rows in result set")),
				)
			},
			wantResLot: lotErr,
			wantResErr:  errErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(testFields)
			}

			response, err := handlers.GetBid(ctx, tt.args)


			if !reflect.DeepEqual(response, tt.wantResLot) {
				t.Errorf("\nGetBid() = %v\nwant = %v", response, tt.wantResLot)
			}
			if status.Code(err) != status.Code(tt.wantResErr) || (err != nil && err.Error() != tt.wantResErr.Error()) {
				t.Errorf("\nGetBid() = %v\nwant = %v", err, tt.wantResErr)
			}
		})
	}

}
