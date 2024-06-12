package chat

import (
	"context"

	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/boards"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var loggerHandlers = basiclogger.BasicLogger{Namespace: "internal.chat.handlers"}

func HandleGetAllMessage(ctx context.Context, boardID string) ([]*model.EventChatMessage, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	board, err := boards.HandleGetBoard(ctx, boardID)
	if err != nil {
		return nil, err
	}

	bid, err := helpers.ProtoToUUID(board.Id)
	if err != nil {
		loggerService.LogError(reqID, "unable to parse board ID", "error", err)
		return nil, status.Error(codes.NotFound, "board not found")
	}

	messages, err := Repository.GetAllChat(ctx, bid)
	if err != nil {
		return nil, err
	}

	return messages, nil
}
