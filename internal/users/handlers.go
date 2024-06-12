package users

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/auth"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var loggerHandlers = basiclogger.BasicLogger{Namespace: "internal.users.handlers"}

func validateHasCtxUser(ctx context.Context) (*model.User, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	u := helpers.GetUserFromContext(ctx)
	if u == nil {
		err := status.Error(codes.PermissionDenied, "permission denied")
		loggerHandlers.LogError(reqID, "unable to load user from context", "error", err)
		return nil, err
	}
	loggerHandlers.LogInfo(reqID, "user loaded from context", "u", u)
	return u, nil
}

func validateHasStoreUser(ctx context.Context, userID uuid.UUID) (*model.User, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	u, err := Repository.GetByID(ctx, userID)
	if err != nil {
		errResponse := status.Error(codes.PermissionDenied, "permission denied")
		loggerHandlers.LogError(reqID, "unable to load user from store", "userID", userID, "error", err, "errorResponse", errResponse)
		return nil, errResponse
	}

	if u == nil {
		err := status.Error(codes.PermissionDenied, "permission denied")
		loggerHandlers.LogError(reqID, "unable to load user from store", "userID", userID, "error", err)
		return nil, err
	}

	return u, nil
}

func validateHasName(ctx context.Context, name string) (string, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	if name == "" {
		err := status.Error(codes.InvalidArgument, "name is required")
		loggerHandlers.LogError(reqID, "empty name parameter provided", "name", name, "error", err)
		return "", err
	}
	return name, nil
}

func validateUserUUID(ctx context.Context, u *model.User) (uuid.UUID, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	userID, err := helpers.ProtoToUUID(u.Id)
	if err != nil {
		errResponse := status.Error(codes.PermissionDenied, "permission denied")
		loggerHandlers.LogError(reqID, "unable to parse user ID", "userID", userID, "error", err, "errorResponse", errResponse)
		return uuid.Nil, errResponse
	}

	return userID, nil
}

func HandleCreateUser(ctx context.Context, r *model.CreateUserRequest) (*model.User, string, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	if _, err := validateHasCtxUser(ctx); err == nil {
		return nil, "", status.Error(codes.Internal, "error creating user")
	}

	name, err := validateHasName(ctx, r.GetName())
	if err != nil {
		return nil, "", err
	}

	user, err := Repository.CreateUser(ctx, name)
	if err != nil {
		loggerHandlers.LogError(reqID, "error creating user", "error", err)
		return nil, "", status.Error(codes.Internal, "error creating user")
	}

	token, err := auth.GenerateUserToken(ctx, user.GetId().GetValue(), user.GetName())
	if err != nil {
		loggerHandlers.LogError(reqID, "error creating user", "error", err)
		return nil, "", status.Error(codes.Internal, "error creating user")
	}

	return user, token, nil
}

func HandleGetUser(ctx context.Context) (*model.User, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	loggerHandlers.LogInfo(reqID, "handle get user: start")
	user, err := validateHasCtxUser(ctx)
	if err != nil {
		return nil, err
	}

	userID, err := validateUserUUID(ctx, user)
	if err != nil {
		return nil, err
	}

	user, err = validateHasStoreUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	loggerHandlers.LogInfo(reqID, "handle get user: done")
	return user, nil
}

func HandleDeleteUser(ctx context.Context, userID uuid.UUID) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	loggerHandlers.LogInfo(reqID, "user delete: start", "userID", userID)
	if _, err := validateHasStoreUser(ctx, userID); err != nil {
		return err
	}

	err := Repository.DeleteUser(ctx, userID)
	if err != nil {
		loggerHandlers.LogError(reqID, "user delete failed", "userID", userID, "error", err)
		return err
	}
	loggerHandlers.LogInfo(reqID, "user delete: done", "userID", userID)

	return nil
}
