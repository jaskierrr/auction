package test

import (
	"context"
	"errors"
	app "main/internal/app/usercases"
	"main/internal/domain/entities"
	service "main/internal/domain/services"
	pb "main/pkg/grpc"
	"main/pkg/logger"
	"main/test/mock"
	"reflect"
	"testing"

	"github.com/golang/mock/gomock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Test_GetUser(t *testing.T) {
	type fields struct {
		userRepo *mock.MockUserRepo
	}

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	logger := logger.NewLogger()

	userRepoMock := mock.NewMockUserRepo(ctrl)

	testFields := &fields{
		userRepo: userRepoMock,
	}

	service := service.NewUserService(userRepoMock)
	usecase := app.NewUserUsecase(service, logger)

	user := &pb.User{
		Id:      1,
		Name:    "Ivan",
		Balance: 1000,
	}

	reqArgDef := pb.GetUserRequest{
		UserId: 1,
	}

	reqArgErr := pb.GetUserRequest{
		UserId: 111,
	}

	userRepoResponse := entities.User{
		Id:      1,
		Name:    "Ivan",
		Balance: 1000,
	}

	userErr, errErr := &pb.UserResponse{}, status.Errorf(codes.Unknown, "failed get user: %v", errors.New("no rows in result set"))

	ctx := context.Background()

	tests := []struct {
		name        string
		args        pb.GetUserRequest
		prepare     func(f *fields)
		wantResUser pb.UserResponse
		wantResErr  error
	}{
		{
			name: "valid",
			args: reqArgDef,
			prepare: func(f *fields) {
				// если указанные вызовы не станут выполняться в ожидаемом порядке, тест будет провален
				gomock.InOrder(
					f.userRepo.EXPECT().GetUser(ctx, &reqArgDef).Return(userRepoResponse, nil),
				)
			},
			wantResUser: pb.UserResponse{User: user},
			wantResErr:  nil,
		},
		{
			name: "wrong_ID",
			args: reqArgErr,
			prepare: func(f *fields) {
				gomock.InOrder(
					f.userRepo.EXPECT().GetUser(ctx, &reqArgErr).Return(entities.User{}, errErr),
				)
			},
			wantResUser: *userErr,
			wantResErr:  errErr,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.prepare != nil {
				tt.prepare(testFields)
			}

			response, err := usecase.GetUser(ctx, &tt.args)

			// func (x *UserResponse) GetUser() *User {
			// 	if x != nil {
			// 		return x.User
			// 	}
			// 	return nil
			// }

			if !reflect.DeepEqual(response, tt.wantResUser) {
				t.Errorf("GetUser() = %v, want = %v", response, tt.wantResUser)
			}
			if !reflect.DeepEqual(err, tt.wantResErr) {
				t.Errorf("GetUser() = %v, want = %v", err, tt.wantResErr)
			}
		})
	}

}
