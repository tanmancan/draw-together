import { WsEvent } from "../../proto-ts/proto/model/events";
import {
  BoardUpdateInput,
  ChatMessageInput,
  DrawingDetectInput,
  DrawingUpdateInput,
  PointerUpdateInput,
  UserUpdateInput,
} from "./events";
import { protoFromBlob } from "./helpers";

export const onMessageHandler = async (e: MessageEvent) => {
  const { data } = e;
  const message = await protoFromBlob<WsEvent>(data, WsEvent);
  let ev: CustomEvent | null = null;

  switch (message?.event?.oneofKind) {
    case "chatMessage": {
      const { chatMessage } = message?.event ?? {};
      ev = new ChatMessageInput({
        detail: chatMessage,
      });
      break;
    }
    case "pointerUpdate": {
      const { pointerUpdate } = message?.event ?? {};
      ev = new PointerUpdateInput({
        detail: pointerUpdate,
      });
      break;
    }
    case "drawingUpdate": {
      const { drawingUpdate } = message?.event ?? {};
      ev = new DrawingUpdateInput({
        detail: drawingUpdate,
      });
      break;
    }
    case "boardUpdate": {
      const { boardUpdate } = message?.event ?? {};
      ev = new BoardUpdateInput({
        detail: boardUpdate,
      });
      break;
    }
    case "userUpdate": {
      const { userUpdate } = message?.event ?? {};
      ev = new UserUpdateInput({
        detail: userUpdate,
      });
      break;
    }
    case "drawingDetect": {
      const { drawingDetect } = message?.event ?? {};
      ev = new DrawingDetectInput({
        detail: drawingDetect,
      });
      break;
    }
    default:
      break;
  }

  if (ev) {
    window.dispatchEvent(ev);
  }
};
