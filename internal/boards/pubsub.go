package boards

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/persistence"
	"github.com/tanmancan/openapi/azurecv"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type BoardPubSubService struct {
	eventConfig *config.EventConfigType
}

var loggerPubSub = basiclogger.BasicLogger{Namespace: "internal.boards.pubsub"}

var BoardPubSub BoardPubSubService

func init() {
	BoardPubSub = BoardPubSubService{
		eventConfig: &config.AppConfig.EventConfig,
	}
}

func (s *BoardPubSubService) QueueDrawingDetection(ctx context.Context, bid *model.UUID) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	user := helpers.GetUserFromContext(ctx)
	rdb := persistence.GetClient()
	channel := s.eventConfig.ChannelDetectQueue.Name

	metadata := &model.EventMetadata{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		BoardId:   bid,
		SenderId:  user.Id,
		CreatedAt: timestamppb.Now(),
	}

	evDrawingDetectQ := model.EventDrawingDetectQueue{
		Metadata: metadata,
		BoardId:  bid,
	}

	protoDrawingDetectEvent := model.WsEvent{
		ReqId: reqID,
		Event: &model.WsEvent_DrawingDetectQueue{
			DrawingDetectQueue: &evDrawingDetectQ,
		},
	}

	protoSer := protoDrawingDetectEvent.String()
	err := rdb.Publish(ctx, channel, protoSer).Err()
	if err != nil {
		loggerPubSub.LogError(reqID, "error publishing drawing detect", "error", err)
		return status.Error(codes.Internal, "error sending drawing detect")
	}

	return nil
}

func (s *BoardPubSubService) SubscribeDetectionQueue(ctx context.Context) {
	cfg := s.eventConfig.ChannelDetectQueue
	cname := cfg.Name
	cid := cfg.Id
	rdb := persistence.GetClient()

	pubsub := rdb.Subscribe(ctx, cname)
	loggerPubSub.LogInfo(cid, "subscribed to channel", "chanName", cname)

	defer func() {
		loggerPubSub.LogWarn(cid, "pubsub handler closed", "chanName", cname)
		pubsub.Close()
	}()

	for {
		reqID := uuid.NewString()
		ctx = context.WithValue(ctx, helpers.RequestID, reqID)
		msg, err := pubsub.ReceiveMessage(ctx)
		if err != nil {
			loggerPubSub.LogError(reqID, "error receiving message", "chanName", cname, "error", err)
			break
		}

		loggerPubSub.LogInfo(reqID, "received message", "msg", msg)
		protoEvent, err := HandleDetectQueueUnmarshalMsg(ctx, msg)
		if err != nil {
			loggerPubSub.LogError(reqID, "error un-marshalling proto string", "chanName", cname, "error", err)
			break
		}

		evReqID := protoEvent.GetReqId()
		if len(evReqID) > 0 {
			reqID = fmt.Sprintf("%s:%s", reqID, evReqID)
		} else {
			protoEvent.ReqId = reqID
		}

		bid := protoEvent.Event.(*model.WsEvent_DrawingDetectQueue).DrawingDetectQueue.BoardId
		loggerPubSub.LogInfo(reqID, "parsed message successfully", "chanName", cname, "boardID", bid)
		err = HandleDrawingDetect(ctx, bid)
		if err != nil {
			loggerPubSub.LogError(reqID, "error detecting drawing", "error", err)
			break
		}
	}
}

func (s *BoardPubSubService) EmitDetectionResult(ctx context.Context, bid *model.UUID, res *azurecv.ImageAnalysisResult) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	rdb := persistence.GetClient()
	channel := s.eventConfig.ChannelDetectDrawing.Name

	metadata := &model.EventMetadata{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		BoardId:   bid,
		SenderId:  bid,
		CreatedAt: timestamppb.Now(),
		Target:    &model.EventMetadata_AllUsers{AllUsers: true},
	}

	evDrawingDetect := model.EventDrawingDetect{
		Metadata:    metadata,
		Description: res.GetCaptionResult().Text,
	}

	protoDrawingDetectEvent := model.WsEvent{
		ReqId: reqID,
		Event: &model.WsEvent_DrawingDetect{
			DrawingDetect: &evDrawingDetect,
		},
	}

	protoSer := protoDrawingDetectEvent.String()
	err := rdb.Publish(ctx, channel, protoSer).Err()
	if err != nil {
		loggerPubSub.LogError(reqID, "error publishing drawing detect", "error", err)
		return status.Error(codes.Internal, "error sending drawing detect")
	}

	loggerPubSub.LogInfo(reqID, "broadcast detect result successful", "boardID", bid)
	return nil
}

func (s *BoardPubSubService) BoardUpdate(ctx context.Context, board *model.Board) error {
	reqID := helpers.GetReqIdFromContext(ctx)

	metadata := &model.EventMetadata{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		BoardId: board.Id,
		SenderId: &model.UUID{
			Value: reqID,
		},
		CreatedAt: timestamppb.Now(),
		Target:    &model.EventMetadata_AllUsers{AllUsers: true},
	}

	boardEv := model.EventBoardUpdate{
		Metadata: metadata,
		Board:    board,
	}

	protoBoardEvent := model.WsEvent{
		ReqId: reqID,
		Event: &model.WsEvent_BoardUpdate{
			BoardUpdate: &boardEv,
		},
	}

	protoSer := protoBoardEvent.String()
	loggerHandlers.LogDebug(reqID, "publishing board update message", "protoSerialized", protoSer)

	rdb := persistence.GetClient()
	channel := s.eventConfig.ChannelUpdateBoard.Name
	err := rdb.Publish(ctx, channel, protoSer).Err()
	if err != nil {
		errResponse := status.Error(codes.Internal, "internal error: unable to update board")
		loggerHandlers.LogError(reqID, "error publishing board update", "error", err, "errorResponse", errResponse)
		return errResponse
	}

	return nil
}

func (s *BoardPubSubService) DrawingUpdate(ctx context.Context, r *model.UpdateDrawingRequest) (bool, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	user := helpers.GetUserFromContext(ctx)
	board, err := HandleGetBoard(ctx, r.GetBoardId().Value)
	if err != nil {
		return false, err
	}

	userID, _ := helpers.ProtoToUUID(user.Id)
	boardID, _ := helpers.ProtoToUUID(board.Id)
	rdb := persistence.GetClient()
	imgData := r.ImageData

	metadata := &model.EventMetadata{
		Id: &model.UUID{
			Value: uuid.NewString(),
		},
		BoardId:   board.Id,
		SenderId:  user.Id,
		CreatedAt: timestamppb.Now(),
		Target:    &model.EventMetadata_AllButSender{AllButSender: true},
	}

	drawings := model.EventDrawingUpdate{
		Metadata:  metadata,
		ImageData: imgData,
	}

	protoDrawingEvent := model.WsEvent{
		ReqId: reqID,
		Event: &model.WsEvent_DrawingUpdate{
			DrawingUpdate: &drawings,
		},
	}

	protoSer := protoDrawingEvent.String()
	loggerHandlers.LogDebug(reqID, "publishing message", "protoSerialized", protoSer)

	channel := s.eventConfig.ChannelUpdateDrawing.Name
	err = rdb.Publish(ctx, channel, protoSer).Err()
	if err != nil {
		loggerHandlers.LogError(reqID, "error publishing drawing update", "error", err)
		return false, status.Error(codes.Internal, "error sending drawing update")
	}

	err = Repository.UpdateBoardDrawing(ctx, boardID, userID, imgData)
	if err != nil {
		loggerHandlers.LogError(reqID, "error saving board drawing", "userID", userID, "boardID", boardID)
		return false, err
	}

	return true, nil
}
