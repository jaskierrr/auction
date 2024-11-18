package test

import (
	"context"
	"errors"
	"main/internal/handlers"
	service "main/internal/services"
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

func Test_Deposit_Balance(t *testing.T) {
	t.Parallel()

	type fields struct {
		userRepo *mock.MockUserRepo
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewLogger()
	validator := validator.New(validator.WithRequiredStructEnabled())

	userRepoMock := mock.NewMockUserRepo(ctrl)

	testFields := &fields{
		userRepo: userRepoMock,
	}

	service := service.NewUserService(userRepoMock)
	handlers := handlers.NewUserHandlers(service, logger, validator)

	validUserReq := &pb.DepositBalanceRequest{
		UserId: 1,
		Amount: 100,
	}

	emtyUserReq := &pb.DepositBalanceRequest{
		UserId: 0,
		Amount: 100,
	}

	emtyAmountReq := &pb.DepositBalanceRequest{
		UserId: 1,
		Amount: 0,
	}

	balanceResponse := &pb.BalanceResponse{
		NewBalance: 100,
	}
	balanceErr := &pb.BalanceResponse{NewBalance: 0}
	err := errors.New("ERROR: new row for relation \"transactions\" violates check constraint \"transactions_amount_check\" (SQLSTATE 23514)\n")

	amountErr := status.Errorf(codes.Unknown, "failed deposite balance: %v", errors.Join(errors.New("cant write transaction"), err))
	userIDerr := status.Errorf(codes.Unknown, "failed deposite balance: %v", errors.New("no rows in result set"))

	ctx := context.Background()

	tests := []struct {
		name        string
		args        *pb.DepositBalanceRequest
		prepare     func(f *fields)
		wantResBalance *pb.BalanceResponse
		wantResErr  error
	}{
		{
			name: "valid",
			args: validUserReq,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.userRepo.EXPECT().DepositBalance(ctx, validUserReq).Return(balanceResponse, nil),
				)
			},
			wantResBalance: balanceResponse,
			wantResErr:  nil,
		},
		{
			name: "wrong_ID",
			args: emtyUserReq,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.userRepo.EXPECT().DepositBalance(ctx, emtyUserReq).Return(&pb.BalanceResponse{}, errors.New("no rows in result set")),
				)
			},
			wantResBalance: balanceErr,
			wantResErr:  userIDerr,
		},
		{
			name: "wrong_amount",
			args: emtyAmountReq,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.userRepo.EXPECT().DepositBalance(ctx, emtyAmountReq).Return(&pb.BalanceResponse{}, errors.Join(errors.New("cant write transaction"), err)),
				)
			},
			wantResBalance: balanceErr,
			wantResErr:  amountErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(testFields)
			}

			response, err := handlers.DepositBalance(ctx, tt.args)

			if !reflect.DeepEqual(response, tt.wantResBalance) {
				t.Errorf("\nCreateUser() = %v\nwant = %v", response, tt.wantResBalance)
			}
			if status.Code(err) != status.Code(tt.wantResErr) || (err != nil && err.Error() != tt.wantResErr.Error()) {
				t.Errorf("CreateUser() = %v\nwant = %v", err, tt.wantResErr)
			}
		})
	}

}
