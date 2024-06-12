package ws

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/persistence"
	"google.golang.org/protobuf/encoding/prototext"
	"google.golang.org/protobuf/proto"
)

type WsPubSubService struct {
	eventConfig config.EventConfigType
}

var WsPubSub WsPubSubService

func init() {
	WsPubSub = WsPubSubService{
		eventConfig: config.AppConfig.EventConfig,
	}
}

func (s WsPubSubService) pubsubHandler(cfg config.ChannelConfigType) func(ctx context.Context, hub *Hub) {
	cname := cfg.Name
	cid := cfg.Id

	return func(ctx context.Context, hub *Hub) {
		rdb := persistence.GetClient()
		pubsub := rdb.Subscribe(ctx, cname)

		loggerHandlers.LogInfo(cid, "subscribed to channel", "chanName", cname)

		defer func() {
			loggerHandlers.LogWarn(cid, "pubsub handler closed", "chanName", cname)
			pubsub.Close()
		}()

		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			reqID := uuid.NewString()
			if err != nil {
				loggerHandlers.LogError(reqID, "error receiving message", "chanName", cname, "error", err)
				break
			}

			loggerHandlers.LogInfo(reqID, "received message")
			var protoEvent model.WsEvent

			err = prototext.Unmarshal([]byte(msg.Payload), &protoEvent)
			if err != nil {
				loggerHandlers.LogError(reqID, "error un-marshalling proto string", "chanName", cname, "error", err)
				break
			}

			evReqID := protoEvent.GetReqId()
			if len(evReqID) > 0 {
				reqID = fmt.Sprintf("%s:%s", reqID, evReqID)
			} else {
				protoEvent.ReqId = reqID
			}

			loggerHandlers.LogInfo(reqID, "parsed message successfully", "chanName", cname)

			protoMsg, err := proto.Marshal(&protoEvent)
			if err != nil {
				loggerHandlers.LogError(reqID, "error marshalling protobuf", "chanName", cname, "error", err)
				break
			}

			loggerHandlers.LogInfo(reqID, "sending message to hub", "chanName", cname)
			hub.messages <- protoMsg
		}
	}
}

func (s WsPubSubService) PointerSubHandler(ctx context.Context, hub *Hub) {
	handler := s.pubsubHandler(s.eventConfig.ChannelUpdatePointer)
	handler(ctx, hub)
}

func (s WsPubSubService) ChatSubHandler(ctx context.Context, hub *Hub) {
	handler := s.pubsubHandler(s.eventConfig.ChannelChatMessage)
	handler(ctx, hub)
}

func (s WsPubSubService) DrawingSubHandler(ctx context.Context, hub *Hub) {
	handler := s.pubsubHandler(s.eventConfig.ChannelUpdateDrawing)
	handler(ctx, hub)
}

func (s WsPubSubService) BoardSubHandler(ctx context.Context, hub *Hub) {
	handler := s.pubsubHandler(s.eventConfig.ChannelUpdateBoard)
	handler(ctx, hub)
}

func (s WsPubSubService) UserSubHandler(ctx context.Context, hub *Hub) {
	handler := s.pubsubHandler(s.eventConfig.ChannelUpdateUser)
	handler(ctx, hub)
}

func (s WsPubSubService) DrawingDetectHandler(ctx context.Context, hub *Hub) {
	handler := s.pubsubHandler(s.eventConfig.ChannelDetectDrawing)
	handler(ctx, hub)
}
