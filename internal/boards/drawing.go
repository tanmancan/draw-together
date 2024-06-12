package boards

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/detection"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/prototext"
)

var loggerDrawings = basiclogger.BasicLogger{Namespace: "internal.boards.drawings"}

func HandleGetBoardDrawings(ctx context.Context, bid *model.UUID) (BoardDrawings, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	boardID, _ := helpers.ProtoToUUID(bid)
	boardDrawing, err := Repository.GetBoardDrawings(ctx, boardID)
	if err != nil {
		loggerDrawings.LogError(reqID, "error loading boards", "error", err)
		return nil, status.Error(codes.Internal, "error loading drawings")
	}

	return boardDrawing, nil
}

func HandleDrawingsToDataLayers(ctx context.Context, drawings BoardDrawings) []*model.ImageData {
	reqId := helpers.GetReqIdFromContext(ctx)

	var dl []*model.ImageData
	for _, imgData := range drawings {
		dl = append(dl, imgData)
	}

	if len(dl) == 0 {
		loggerDrawings.LogWarn(reqId, "no image layers found")
	}

	return dl
}

func HandleDrawingDetect(ctx context.Context, bid *model.UUID) error {
	reqID := helpers.GetReqIdFromContext(ctx)
	boardDrawings, err := HandleGetBoardDrawings(ctx, bid)
	if err != nil {
		return err
	}

	dl := HandleDrawingsToDataLayers(ctx, boardDrawings)
	if len(dl) == 0 {
		return status.Error(codes.NotFound, "no drawings found")
	}

	gis := detection.NewGenerateImageService(ctx)
	gis.SetBoardImageDataLayers(dl)
	reader, err := gis.BoardDrawingsToImageBuffer()
	if err != nil {
		return status.Error(codes.NotFound, "no drawings found")
	}

	ds := detection.NewDetectionService(ctx)
	ds.SetBoardImageReader(reader)
	res, err := ds.DetectImage()
	if err != nil {
		loggerDrawings.LogError(reqID, "error loading images")
		return status.Error(codes.Internal, "error loading images")
	}

	err = BoardPubSub.EmitDetectionResult(ctx, bid, res)
	if err != nil {
		loggerDrawings.LogError(reqID, "error detecting image", "error", err)
		return status.Error(codes.Internal, "error detecting image")
	}

	return nil
}

func HandleDetectQueueUnmarshalMsg(ctx context.Context, msg *redis.Message) (*model.WsEvent, error) {
	var protoEvent model.WsEvent
	var err error = nil
	err = prototext.Unmarshal([]byte(msg.Payload), &protoEvent)
	if err != nil {
		return nil, err
	}

	switch protoEvent.Event.(type) {
	case *model.WsEvent_DrawingDetectQueue:
	default:
		err = errors.New("invalid message type")
		protoEvent = model.WsEvent{}
	}

	return &protoEvent, err
}
