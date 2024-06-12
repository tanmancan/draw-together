package pointer

import (
	"context"

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

var loggerPubSub = basiclogger.BasicLogger{Namespace: "internal.pointer.pubsub"}

type PointerPubSubService struct {
	eventConfig *config.EventConfigType
}

var PointerPubSub PointerPubSubService

func init() {
	PointerPubSub = PointerPubSubService{
		eventConfig: &config.AppConfig.EventConfig,
	}
}

func (s *PointerPubSubService) UpdatePointer(ctx context.Context, r *model.UpdatePointerRequest) (bool, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	user := helpers.GetUserFromContext(ctx)
	board, err := boards.HandleGetBoard(ctx, r.GetBoardId().GetValue())
	if err != nil {
		return false, err
	}

	rdb := persistence.GetClient()
	events := r.PointerPositions

	metadata := &model.EventMetadata{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		BoardId:   board.Id,
		SenderId:  user.Id,
		CreatedAt: timestamppb.Now(),
		Target:    &model.EventMetadata_AllButSender{AllButSender: true},
	}
	pointers := model.EventPointerUpdate{
		Metadata:         metadata,
		PointerPositions: events,
	}
	protoPointerEvent := model.WsEvent{
		ReqId: reqID,
		Event: &model.WsEvent_PointerUpdate{
			PointerUpdate: &pointers,
		},
	}

	protoSer := protoPointerEvent.String()
	loggerPubSub.LogDebug(reqID, "publishing message", "protoSerialized", protoSer)

	channel := config.AppConfig.EventConfig.ChannelUpdatePointer.Name
	loggerPubSub.LogInfo(reqID, "publish pointer message start", "chanName", channel, "userID", user.Id.Value, "boardID", board.Id.Value)
	err = rdb.Publish(ctx, channel, protoSer).Err()
	if err != nil {
		loggerPubSub.LogError(reqID, "error publishing pointer update", "error", err)
		return false, status.Error(codes.Internal, "error sending pointer update")
	}

	loggerPubSub.LogInfo(reqID, "publish pointer message success", "chanName", channel, "userID", user.Id.Value, "boardID", board.Id.Value)

	return true, nil
}
