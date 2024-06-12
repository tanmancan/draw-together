package ws

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/config"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/protobuf/proto"
)

var loggerClients = basiclogger.BasicLogger{Namespace: "internal.ws.client"}

type Client struct {
	id uuid.UUID

	hub *Hub

	con *websocket.Conn

	send chan []byte
}

func (c *Client) Read(ctx context.Context) {
	initReqID := helpers.GetReqIdFromContext(ctx)
	loggerClients.LogInfo(initReqID, "client read started", "client_id", c.id)

	defer func() {
		c.hub.unregister <- c.id
		c.con.Close()
		loggerClients.LogWarn(initReqID, "client closed", "clientID", c.id)
	}()

	c.con.SetReadLimit(config.AppConfig.NetworkServiceWebsocket.ReadLimit)
	readDeadline := config.AppConfig.NetworkServiceWebsocket.ReadDeadline
	c.con.SetReadDeadline(time.Now().Add(readDeadline))
	c.con.SetPongHandler(func(string) error {
		c.con.SetReadDeadline(time.Now().Add(readDeadline))
		return nil
	})

	for {
		messageType, msg, err := c.con.ReadMessage()
		reqID := uuid.NewString()
		loggerClients.LogInfo(reqID, "client read message received", "messageType", messageType)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				loggerClients.LogError(reqID, "unexpected close error")
			}

			loggerClients.LogError(reqID, "error reading message", "error", err, "message", msg, "messageType", messageType)
			break
		}

		var message model.WsEvent
		proto.Unmarshal(msg, &message)
		evReqID := message.GetReqId()
		if len(evReqID) > 0 {
			reqID = fmt.Sprintf("%s:%s", reqID, evReqID)
		} else {
			message.ReqId = reqID
		}

		loggerClients.LogInfo(reqID, "client read message sent to hub", "messageType", messageType)
		c.hub.messages <- msg
	}
}

func (c *Client) Write(ctx context.Context) {
	reqID := helpers.GetReqIdFromContext(ctx)
	loggerClients.LogInfo(reqID, "client write: started", "channel_id", c.id)
	ticker := time.NewTicker(config.AppConfig.NetworkServiceWebsocket.PingInterval)

	defer func() {
		c.con.Close()
		loggerClients.LogWarn(reqID, "client closed", "clientID", c.id)
	}()

	for {
		select {
		case msg, ok := <-c.send:
			writeDeadline := config.AppConfig.NetworkServiceWebsocket.WriteDeadline
			c.con.SetWriteDeadline(time.Now().Add(writeDeadline))
			if !ok {
				loggerClients.LogError(reqID, "error retrieving message from chan send")
				c.con.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			loggerClients.LogInfo(reqID, "client write: writing message", "clientID", c.id)
			loggerClients.LogDebug(reqID, "writing message", "clientID", c.id, "message", msg)
			c.con.WriteMessage(websocket.BinaryMessage, msg)
		case <-ticker.C:
			writeDeadline := config.AppConfig.NetworkServiceWebsocket.WriteDeadline
			c.con.SetWriteDeadline(time.Now().Add(writeDeadline))
			loggerClients.LogDebug(reqID, "pinging client", "clientID", c.id)
			err := c.con.WriteMessage(websocket.PingMessage, []byte("ping"))
			if err != nil {
				loggerClients.LogError(reqID, "error pinging client", "clientID", c.id, "error", err)
				return
			}
		}
	}
}
