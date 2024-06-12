package main

import (
	"context"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/boards"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
)

func main() {
	config.SetupConfigFlags()
	reqID := uuid.NewString()
	ctx := context.WithValue(context.Background(), helpers.RequestID, reqID)

	boards.BoardPubSub.SubscribeDetectionQueue(ctx)
}
