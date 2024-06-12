package ws

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/tanmancan/draw-together/internal/auth"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
)

var loggerHandlers = basiclogger.BasicLogger{Namespace: "internal.ws.handlers"}

var upgrader = websocket.Upgrader{
	CheckOrigin:     helpers.CheckOriginHandler,
	ReadBufferSize:  config.AppConfig.NetworkServiceWebsocket.ReadBuffer,
	WriteBufferSize: config.AppConfig.NetworkServiceWebsocket.WriteBuffer,
}

func WsInitHandler(hub *Hub, w http.ResponseWriter, r *http.Request) {
	r = r.WithContext(helpers.CreateContextWithRequestID(r.Context()))
	ctx := r.Context()
	reqID := helpers.GetReqIdFromContext(r.Context())
	q := r.URL.Query()
	token := q.Get("token")
	userID, err := auth.JwtAuthHandler(ctx, token)
	if err != nil {
		loggerHandlers.LogError(reqID, "Error upgrading connection", "error", err)
		return
	}

	if hc, ok := hub.clients[userID]; ok {
		loggerHandlers.LogWarn(reqID, "client is already registered to hub", "userID", userID)
		delete(hub.clients, userID)
		close(hc.send)
		loggerHandlers.LogWarn(reqID, "closed and removed old client", "userID", userID)
	}

	loggerHandlers.LogInfo(reqID, "websocket init", "request", r)

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		loggerHandlers.LogError(reqID, "Error upgrading connection", "error", err)
		return
	}

	loggerHandlers.LogInfo(reqID, "client successfully upgraded")

	client := &Client{
		hub:  hub,
		id:   userID,
		con:  c,
		send: make(chan []byte),
	}

	client.hub.register <- client
	loggerHandlers.LogInfo(reqID, "client registered", "userID", userID)

	go client.Read(ctx)
	go client.Write(ctx)
}
