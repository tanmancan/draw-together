package boards

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var loggerHandlers = basiclogger.BasicLogger{Namespace: "internal.boards.handlers"}

func HandleCreateBoard(ctx context.Context, boardName string) (*model.Board, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	user := helpers.GetUserFromContext(ctx)
	if user == nil {
		err := status.Error(codes.PermissionDenied, "permission denied")
		loggerHandlers.LogError(reqID, "user not found", "error", err)
		return nil, err
	}

	if boardName == "" {
		logMsg := fmt.Sprintf("empty name parameter provided: %s", boardName)
		err := status.Error(codes.InvalidArgument, "invalid argument: name is required")
		loggerHandlers.LogError(reqID, logMsg, "error", err)
		return nil, err
	}

	board, err := Repository.CreateBoard(ctx, boardName, user)

	if err != nil {
		loggerHandlers.LogError(reqID, "error creating board", "error", err)
		return nil, status.Error(codes.Internal, "internal error: unable to create board")
	}

	return board, nil
}

func HandleGetBoard(ctx context.Context, id string) (*model.Board, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	userLoad := helpers.GetUserFromContext(ctx)
	if userLoad == nil {
		err := status.Error(codes.PermissionDenied, "permission denied")
		loggerHandlers.LogError(reqID, "user not found", "error", err)
		return nil, err
	}

	userUuid, err := helpers.ProtoToUUID(userLoad.Id)
	if err != nil {
		errResponse := status.Error(codes.PermissionDenied, "permission denied")
		loggerHandlers.LogError(reqID, "error parsing user id", "error", err, "errorResponse", errResponse)
		return nil, errResponse
	}

	if id == "" {
		loggerHandlers.LogError(reqID, "board id not provided", "boardID", id, "code", codes.InvalidArgument)
		return nil, status.Error(codes.InvalidArgument, "invalid argument: board id required")
	}

	boardID, err := uuid.Parse(id)
	if err != nil {
		loggerHandlers.LogError(reqID, "error parsing board id:", "boardID", id, "code", codes.InvalidArgument)
		return nil, status.Error(codes.InvalidArgument, "invalid argument: valid board id required")
	}

	addedNewUser := false
	if ok := Repository.BoardHasUser(ctx, boardID, userUuid); !ok {
		addedNewUser, err = HandleAddBoardUser(ctx, boardID, userUuid)
		if err != nil {
			return nil, err
		}
	}

	board, err := Repository.GetByID(ctx, boardID)
	if err != nil {
		loggerHandlers.LogError(reqID, "error fetching board", "error", err)
		return nil, status.Error(codes.NotFound, "not found")
	}

	if addedNewUser {
		err := BoardPubSub.BoardUpdate(ctx, board)
		if err != nil {
			loggerHandlers.LogError(reqID, "error sending board update", "error", err)
			return nil, status.Error(codes.Internal, "internal error: unable to load board")
		}
	}

	return board, nil
}

func HandleAddBoardUser(ctx context.Context, boardID uuid.UUID, userID uuid.UUID) (bool, error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	if ok := Repository.BoardHasUser(ctx, boardID, userID); ok {
		loggerHandlers.LogWarn(reqID, "noop: user already exists in board", "userID", userID, "boardID", boardID)
		return true, nil
	}

	loggerRepository.LogInfo(reqID, "user does not exist in board. continue add user", "userID", userID, "boardID", boardID)
	err := Repository.AddUser(ctx, boardID, userID)
	if err != nil {
		errResponse := status.Error(codes.Internal, "internal error: unable to add user to board")
		loggerRepository.LogError(reqID, "error adding user to board", "error", err)
		return false, errResponse
	}
	loggerRepository.LogInfo(reqID, "added user to board", "userID", userID)

	return true, nil
}
