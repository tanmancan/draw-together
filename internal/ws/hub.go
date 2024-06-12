package ws

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/tanmancan/draw-together/internal/basiclogger"
	"github.com/tanmancan/draw-together/internal/boards"
	"github.com/tanmancan/draw-together/internal/helpers"
	"github.com/tanmancan/draw-together/internal/model"
	"google.golang.org/protobuf/proto"
)

var loggerHub = basiclogger.BasicLogger{Namespace: "internal.ws.hub"}

type Hub struct {
	clients map[uuid.UUID]*Client

	messages chan []byte

	register chan *Client

	unregister chan uuid.UUID
}

func (h *Hub) Run(ctx context.Context) {
	initReqID := helpers.GetReqIdFromContext(ctx)
	loggerHub.LogInfo(initReqID, "hub started")

	for {
		select {
		case client := <-h.register:
			h.clients[client.id] = client
			loggerHub.LogInfo(initReqID, "registered client", "clientID", client.id)
		case clientID := <-h.unregister:
			if c, ok := h.clients[clientID]; ok {
				delete(h.clients, clientID)
				close(c.send)
				loggerHub.LogInfo(initReqID, "un-registered client", "clientID", clientID)
			}
		case message := <-h.messages:
			reqID := uuid.NewString()
			loggerHub.LogInfo(reqID, "hub received message")
			loggerHub.LogDebug(reqID, "hub received message", "message", message)
			var protoMessage model.WsEvent
			proto.Unmarshal(message, &protoMessage)

			evReqID := protoMessage.GetReqId()
			if len(evReqID) > 0 {
				reqID = fmt.Sprintf("%s:%s", reqID, evReqID)
			}

			protoMessage.ReqId = reqID
			ctx = context.WithValue(ctx, helpers.RequestID, reqID)
			metadata := getEventMetadata(&protoMessage)
			if metadata == nil {
				loggerClients.LogWarn(reqID, "invalid event metadata found", "protoMessage", &protoMessage, "metadata", metadata)
			}

			allTarget := metadata.GetAllUsers()
			allTargetButSender := metadata.GetAllButSender()
			loggerHub.LogInfo(reqID, "sent message to all target", "allTarget", allTarget, "allTargetButSender", allTargetButSender)
			if allTarget || allTargetButSender {
				sendAllTarget(ctx, h, message, metadata, allTargetButSender)
			} else {
				sendSingleTarget(ctx, h, message, metadata)
			}
		}
	}
}

func getEventMetadata(protoMessage *model.WsEvent) *model.EventMetadata {
	var metadata *model.EventMetadata
	switch protoMessage.Event.(type) {
	case *model.WsEvent_ChatMessage:
		metadata = protoMessage.GetChatMessage().Metadata
	case *model.WsEvent_PointerUpdate:
		metadata = protoMessage.GetPointerUpdate().Metadata
	case *model.WsEvent_DrawingUpdate:
		metadata = protoMessage.GetDrawingUpdate().Metadata
	case *model.WsEvent_BoardUpdate:
		metadata = protoMessage.GetBoardUpdate().Metadata
	case *model.WsEvent_UserUpdate:
		metadata = protoMessage.GetUserUpdate().Metadata
	case *model.WsEvent_DrawingDetect:
		metadata = protoMessage.GetDrawingDetect().Metadata
	}

	return metadata
}

func sendAllTarget(ctx context.Context, h *Hub, message []byte, metadata *model.EventMetadata, excludeSender bool) {
	reqID := helpers.GetReqIdFromContext(ctx)
	senderID, err := helpers.ProtoToUUID(metadata.SenderId)
	if err != nil {
		loggerHub.LogError(reqID, "error parsing sender id", "error", err)
	}

	boardID, err := helpers.ProtoToUUID(metadata.BoardId)
	if err != nil {
		loggerHub.LogError(reqID, "error parsing board id", "error", err)
	}

	loggerHub.LogInfo(reqID, "send message to all clients: start", "senderID", senderID)
	loggerHub.LogDebug(reqID, "sending to all clients except sender", "senderID", senderID, "clients", h.clients)
	cFound := false
	for id, client := range h.clients {
		if ok := boards.Repository.BoardHasUser(ctx, boardID, id); ok {
			if excludeSender && id == senderID {
				cFound = true
				continue
			}
			loggerHub.LogInfo(reqID, "hub client found", "userID", id, "clientID", client.id)
			client.send <- message
			loggerHub.LogInfo(reqID, "sending message to client", "userID", id, "clientID", client.id)
			cFound = true
		}
	}

	if !cFound {
		loggerHub.LogWarn(reqID, "dropping message: no client found in hub for board message")
	}
}

func sendSingleTarget(ctx context.Context, h *Hub, message []byte, metadata *model.EventMetadata) {
	userTarget := metadata.GetUserId()
	reqID := helpers.GetReqIdFromContext(ctx)
	if userTarget != nil {
		boardID, err := helpers.ProtoToUUID(metadata.BoardId)
		if err != nil {
			loggerHub.LogError(reqID, "error parsing board id", "error", err)
		}

		targetID, err := helpers.ProtoToUUID(userTarget)
		if err != nil {
			loggerHub.LogError(reqID, "error parsing target id", "error", err)
		}

		if ok := boards.Repository.BoardHasUser(ctx, boardID, targetID); ok {
			if client, ok := h.clients[targetID]; ok {
				loggerHub.LogInfo(reqID, "hub client found", "clientID", client.id.String())
				client.send <- message
				loggerHub.LogInfo(reqID, "hub client send message: success", "clientID", client.id.String())
				loggerHub.LogDebug(reqID, "sent message", "message", message)
			} else {
				loggerHub.LogWarn(reqID, "dropping message. client not found", "userTarget", userTarget)
			}
		}
	}
}

func BuildHub() *Hub {
	return &Hub{
		clients:    make(map[uuid.UUID]*Client),
		messages:   make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan uuid.UUID),
	}
}
