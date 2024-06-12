package chat

import (
	"context"

	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ChatServiceServerImpl struct {
	service.UnimplementedChatServiceServer
}

var loggerService = basiclogger.BasicLogger{Namespace: "internal.chat.service"}

func (s *ChatServiceServerImpl) SendMessage(ctx context.Context, r *model.ChatMessageRequest) (*model.ChatMessageResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	loggerService.LogDebug(reqID, "received message", "message", r.GetMessage())

	success, err := ChatPubSub.PublishChatMessage(ctx, r)

	if err != nil {
		loggerHandlers.LogError(reqID, "error publishing message", "error", err)
		return nil, err
	}

	if !success {
		loggerHandlers.LogError(reqID, "error publishing message")
		return nil, status.Error(codes.Internal, "error publishing message")
	}

	return &model.ChatMessageResponse{
		Success: true,
	}, nil
}

func (s *ChatServiceServerImpl) GetBoardMessages(ctx context.Context, r *model.GetBoardMessagesRequest) (*model.GetBoardMessagesResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	messages, err := HandleGetAllMessage(ctx, r.BoardId.GetValue())
	if err != nil {
		return nil, err
	}

	loggerService.LogDebug(reqID, "messages loaded", "message_count", len(messages))

	res := model.GetBoardMessagesResponse{
		Messages: messages,
	}

	return &res, nil
}
