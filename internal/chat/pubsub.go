package chat

import (
	"context"
	"strings"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/boards"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/persistence"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ChatPubSubService struct {
	eventConfig *config.EventConfigType
}

var loggerPubSub = basiclogger.BasicLogger{Namespace: "internal.chat.pubsub"}

var ChatPubSub ChatPubSubService

func init() {
	ChatPubSub = ChatPubSubService{
		eventConfig: &config.AppConfig.EventConfig,
	}
}

func (s *ChatPubSubService) PublishChatMessage(ctx context.Context, r *model.ChatMessageRequest) (bool, error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	if len(strings.Trim(r.Message[0], " ")) == 0 || r.Message == nil {
		return false, nil
	}

	user := helpers.GetUserFromContext(ctx)
	board, err := boards.HandleGetBoard(ctx, r.GetBoardId().Value)
	if err != nil {
		return false, err
	}

	boardID, _ := helpers.ProtoToUUID(board.Id)
	rdb := persistence.GetClient()

	for _, message := range r.GetMessage() {
		metadata := &model.EventMetadata{
			Id: &model.UUID{
				Value: uuid.NewString(),
			},
			BoardId:   board.Id,
			SenderId:  user.Id,
			CreatedAt: timestamppb.Now(),
			Target:    &model.EventMetadata_AllUsers{AllUsers: true},
		}
		chat := model.EventChatMessage{
			Metadata: metadata,
			Body:     strings.Trim(message, " "),
		}
		protoChatEvent := model.WsEvent{
			Event: &model.WsEvent_ChatMessage{
				ChatMessage: &chat,
			},
		}

		protoSer := protoChatEvent.String()
		loggerPubSub.LogDebug(reqID, "publishing message", "protoSerialized", protoSer)

		channel := config.AppConfig.EventConfig.ChannelChatMessage.Name
		err := rdb.Publish(ctx, channel, protoSer).Err()
		if err != nil {
			loggerPubSub.LogError(reqID, "error publishing message", "error", err)
			return false, status.Error(codes.Internal, "error sending message")
		}
		Repository.AddChat(ctx, boardID, &chat)
	}

	return true, nil
}
