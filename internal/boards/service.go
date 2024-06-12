package boards

import (
	"context"

	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BoardServiceServerImpl struct {
	service.UnimplementedBoardServiceServer
}

var loggerService = basiclogger.BasicLogger{Namespace: "internal.boards.service"}

func (s *BoardServiceServerImpl) CreateBoard(ctx context.Context, r *model.CreateBoardRequest) (*model.GetBoardResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	board, err := HandleCreateBoard(ctx, r.GetName())

	if err != nil {
		loggerService.LogError(reqID, "error creating board", "error", err)
		return nil, err
	}

	return &model.GetBoardResponse{
		Board: board,
	}, nil
}

func (s *BoardServiceServerImpl) GetBoard(ctx context.Context, r *model.GetBoardRequest) (*model.GetBoardResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)

	board, err := HandleGetBoard(ctx, r.GetId().Value)

	if err != nil {
		loggerService.LogError(reqID, "error fetching board", "error", err)
		return nil, err
	}

	return &model.GetBoardResponse{
		Board: board,
	}, nil
}

func (s *BoardServiceServerImpl) DeleteBoard(ctx context.Context, r *model.GetBoardRequest) (*model.GetBoardResponse, error) {
	return nil, nil
}

func (s *BoardServiceServerImpl) UpdateDrawing(ctx context.Context, r *model.UpdateDrawingRequest) (*model.UpdateDrawingResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	success, err := BoardPubSub.DrawingUpdate(ctx, r)
	if err != nil {
		return nil, err
	}
	res := &model.UpdateDrawingResponse{
		Success: success,
	}
	loggerService.LogDebug(reqID, "successfully updated drawing event")
	return res, nil
}

func (s *BoardServiceServerImpl) GetBoardDrawings(ctx context.Context, r *model.GetBoardDrawingsRequest) (*model.GetBoardDrawingsResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	bid := r.GetBoardId()
	boardDrawings, err := HandleGetBoardDrawings(ctx, bid)
	if err != nil {
		return nil, err
	}

	res := &model.GetBoardDrawingsResponse{
		Drawings: boardDrawings,
	}

	loggerService.LogInfo(reqID, "fetched board drawings", "boardID", r.BoardId.GetValue(), "count_drawing", len(boardDrawings))
	return res, nil
}

func (s *BoardServiceServerImpl) DrawingDetect(ctx context.Context, r *model.GetBoardDrawingsRequest) (*model.DrawingDetectResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	bid := r.GetBoardId()
	err := BoardPubSub.QueueDrawingDetection(ctx, bid)
	if err != nil {
		loggerService.LogError(reqID, "error queuing drawing for detection", "error", err)
		return nil, status.Error(codes.Internal, "error detecting drawing")
	}

	res := &model.DrawingDetectResponse{
		Success: true,
	}

	loggerService.LogInfo(reqID, "drawing detect successful", "boardID", bid)
	return res, nil
}
