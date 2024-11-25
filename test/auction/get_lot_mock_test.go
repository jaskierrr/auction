package test

import (
	"context"
	"errors"
	"main/internal/entities"
	"main/internal/handlers"
	auction_service "main/internal/services/auction"
	pb "main/pkg/grpc"
	"main/pkg/logger"
	"main/test/mock"
	"reflect"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Get_Lot(t *testing.T) {
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

	service := auction_service.NewAuctionService(auctionRepoMock, logger)
	handlers := handlers.NewAuctionHandlers(service, logger, validator)

	lot := &pb.Lot{
		Id:          1,
		Title:       "Title",
		Description: "Desc",
		StartingBid: 100,
		SellerId:    1,
		Status:      "Active",
	}

	reqArgDef := &pb.GetLotRequest{
		LotId: 1,
	}

	reqArgErr := &pb.GetLotRequest{
		LotId: -1,
	}

	lotRepoResponse := entities.Lot{
		Id:          1,
		Title:       "Title",
		Description: "Desc",
		StartingBid: 100,
		SellerId:    1,
		Status:      "Active",
	}

	lotErr, errErr := &pb.LotResponse{Lot: &pb.Lot{}}, status.Errorf(codes.Unknown, "failed get lot: %v", errors.New("no rows in result set"))

	ctx := context.Background()

	tests := []struct {
		name       string
		args       *pb.GetLotRequest
		prepare    func(f *fields)
		wantResLot *pb.LotResponse
		wantResErr error
	}{
		{
			name: "valid",
			args: reqArgDef,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auctionRepo.EXPECT().GetLot(ctx, reqArgDef).Return(lotRepoResponse, nil),
				)
			},
			wantResLot: &pb.LotResponse{Lot: lot},
			wantResErr: nil,
		},
		{
			name: "wrong_ID",
			args: reqArgErr,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.auctionRepo.EXPECT().GetLot(ctx, reqArgErr).Return(entities.Lot{}, errors.New("no rows in result set")),
				)
			},
			wantResLot: lotErr,
			wantResErr: errErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(testFields)
			}

			response, err := handlers.GetLot(ctx, tt.args)

			if !reflect.DeepEqual(response, tt.wantResLot) {
				t.Errorf("\nGetLot() = %v\nwant = %v", response, tt.wantResLot)
			}
			if status.Code(err) != status.Code(tt.wantResErr) || (err != nil && err.Error() != tt.wantResErr.Error()) {
				t.Errorf("\nGetLot() = %v\nwant = %v", err, tt.wantResErr)
			}
		})
	}

}
