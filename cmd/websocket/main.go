package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/ws"
)

var logger = basiclogger.BasicLogger{Namespace: "websocket.main"}

func main() {
	config.SetupConfigFlags()
	reqID := uuid.NewString()
	ctx := context.WithValue(context.Background(), helpers.RequestID, reqID)

	hub := ws.BuildHub()

	go hub.Run(ctx)
	go ws.WsPubSub.PointerSubHandler(ctx, hub)
	go ws.WsPubSub.ChatSubHandler(ctx, hub)
	go ws.WsPubSub.DrawingSubHandler(ctx, hub)
	go ws.WsPubSub.BoardSubHandler(ctx, hub)
	go ws.WsPubSub.DrawingDetectHandler(ctx, hub)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ws.WsInitHandler(hub, w, r)
	})

	h := config.AppConfig.NetworkServiceWebsocket.Host
	p := config.AppConfig.NetworkServiceWebsocket.Port
	s := fmt.Sprintf("%s:%d", h, p)

	msg := fmt.Sprintf("websocket service listening at %s:%d", h, p)
	logger.LogInfo(reqID, msg)

	keyFile := "/misc/cert/cert.key"
	certFile := "/misc/cert/cert.crt"
	log.Fatal(http.ListenAndServeTLS(s, certFile, keyFile, nil), reqID)
}
