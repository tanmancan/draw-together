package pointer

import (
	"context"

	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"github.com/tanmancan/draw-together/internal/service"
)

type PointerServiceServerImpl struct {
	service.UnimplementedPointerServiceServer
}

var loggerService = basiclogger.BasicLogger{Namespace: "internal.pointer.service"}

func (s *PointerServiceServerImpl) UpdatePointer(ctx context.Context, r *model.UpdatePointerRequest) (*model.UpdatePointerResponse, error) {
	reqID := helpers.GetReqIdFromContext(ctx)
	loggerService.LogInfo(reqID, "update pointer: start")
	success, err := PointerPubSub.UpdatePointer(ctx, r)
	if err != nil {
		return nil, err
	}
	res := &model.UpdatePointerResponse{
		Success: success,
	}
	loggerService.LogInfo(reqID, "update pointer: done")
	return res, nil
}
