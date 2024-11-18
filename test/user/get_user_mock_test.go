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

	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_Get_User(t *testing.T) {
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

	user := &pb.User{
		Id:      1,
		Name:    "Ivan",
		Balance: 1000,
	}

	reqArgDef := &pb.GetUserRequest{
		UserId: 1,
	}

	reqArgErr := &pb.GetUserRequest{
		UserId: -1,
	}

	userRepoResponse := entities.User{
		Id:      1,
		Name:    "Ivan",
		Balance: 1000,
	}

	userErr, errErr := &pb.UserResponse{User: &pb.User{}}, status.Errorf(codes.Unknown, "failed get user: %v", errors.New("no rows in result set"))

	ctx := context.Background()

	tests := []struct {
		name        string
		args        *pb.GetUserRequest
		prepare     func(f *fields)
		wantResUser *pb.UserResponse
		wantResErr  error
	}{
		{
			name: "valid",
			args: reqArgDef,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.userRepo.EXPECT().GetUser(ctx, reqArgDef).Return(userRepoResponse, nil),
				)
			},
			wantResUser: &pb.UserResponse{User: user},
			wantResErr:  nil,
		},
		{
			name: "wrong_ID",
			args: reqArgErr,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.userRepo.EXPECT().GetUser(ctx, reqArgErr).Return(entities.User{}, errors.New("no rows in result set")),
				)
			},
			wantResUser: userErr,
			wantResErr:  errErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(testFields)
			}

			response, err := handlers.GetUser(ctx, tt.args)


			if !reflect.DeepEqual(response, tt.wantResUser) {
				t.Errorf("\nGetUser() = %v\nwant = %v", response, tt.wantResUser)
			}
			if status.Code(err) != status.Code(tt.wantResErr) || (err != nil && err.Error() != tt.wantResErr.Error()) {
				t.Errorf("\nGetUser() = %v\nwant = %v", err, tt.wantResErr)
			}
		})
	}

}
