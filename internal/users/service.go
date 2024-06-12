package users

import (
	"context"

	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/service"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserServiceServerImpl struct {
	service.UnimplementedUserServiceServer
}

var loggerService basiclogger.BasicLogger = basiclogger.BasicLogger{Namespace: "internal.grpc.service"}

func (s *UserServiceServerImpl) CreateUser(ctx context.Context, r *model.CreateUserRequest) (*model.CreateUserResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	u, token, err := HandleCreateUser(ctx, r)

	if err != nil {
		loggerService.LogError(reqID, "error creating user", "error", err)
		return nil, err
	}

	return &model.CreateUserResponse{
		User:  u,
		Token: token,
	}, nil
}

func (s *UserServiceServerImpl) GetUser(ctx context.Context, e *emptypb.Empty) (*model.GetUserResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	u, err := HandleGetUser(ctx)

	if err != nil {
		loggerService.LogError(reqID, "error fetching user", "error", err)
		return nil, err
	}

	return &model.GetUserResponse{
		User: u,
	}, nil
}

func (s *UserServiceServerImpl) DeleteUser(ctx context.Context, e *emptypb.Empty) (*model.GetUserResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	u, err := HandleGetUser(ctx)

	if err != nil {
		loggerService.LogError(reqID, "error fetching user", "error", err)
		return nil, err
	}

	return &model.GetUserResponse{
		User: u,
	}, nil
}
