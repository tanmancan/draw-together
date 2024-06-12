package chat

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/persistence"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
)

type ChatRepository interface {
	AddChat(ctx context.Context, boardID uuid.UUID, chat *model.EventChatMessage) error
	DeleteChat(ctx context.Context, boardID uuid.UUID, chatID string) error
	GetAllChat(ctx context.Context, boardID uuid.UUID) ([]*model.EventChatMessage, error)
}

type ChatRepositoryImpl struct {
	client            *redis.Client
	primaryHashPrefix string
}

var Repository ChatRepository

var loggerRepository = basiclogger.BasicLogger{Namespace: "internal.chat.repository"}

func init() {
	MakeChatRepository(nil)
}

func MakeChatRepository(r ChatRepository) {
	if r != nil {
		Repository = r
	} else {
		Repository = &ChatRepositoryImpl{
			persistence.GetClient(),
			"chat",
		}
	}
}

func (r *ChatRepositoryImpl) MakeChatHash(boardID uuid.UUID) string {
	return fmt.Sprintf("%s:%s", r.primaryHashPrefix, boardID.String())
}

func (r *ChatRepositoryImpl) AddChat(ctx context.Context, boardID uuid.UUID, chat *model.EventChatMessage) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	chatHash := r.MakeChatHash(boardID)

	err := r.client.HSet(ctx, chatHash, chat.Metadata.GetId().GetValue(), chat.String()).Err()
	if err != nil {
		loggerRepository.LogError(reqID, "error saving chat message", "error", err)
		return status.Error(codes.Internal, "error saving chat message")
	}

	return nil
}

func (r *ChatRepositoryImpl) DeleteChat(ctx context.Context, boardID uuid.UUID, chatID string) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	chatHash := r.MakeChatHash(boardID)

	err := r.client.HDel(ctx, chatHash, chatID).Err()
	if err != nil {
		loggerRepository.LogError(reqID, "error deleting chat message", "error", err)
		return status.Error(codes.Internal, "error deleting chat message")
	}

	return nil
}

func (r *ChatRepositoryImpl) GetAllChat(ctx context.Context, boardID uuid.UUID) ([]*model.EventChatMessage, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	chatHash := r.MakeChatHash(boardID)

	chatSer, err := r.client.HGetAll(ctx, chatHash).Result()
	if err != nil {
		loggerRepository.LogError(reqID, "error loading messages", "boardID", boardID, "error", err)
		return nil, status.Error(codes.Internal, "error loading messages")
	}

	chatCollection := []*model.EventChatMessage{}

	for _, chat := range chatSer {
		var c = &model.EventChatMessage{}
		err := prototext.Unmarshal([]byte(chat), c)
		if err != nil {
			loggerRepository.LogInfo(reqID, "error un marshalling proto", "error", err)
			continue
		}
		chatCollection = append(chatCollection, c)
	}

	return chatCollection, nil
}
