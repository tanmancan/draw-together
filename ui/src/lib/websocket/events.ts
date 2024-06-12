import { EventBoardUpdate } from "../../proto-ts/proto/model/board";
import { EventChatMessage } from "../../proto-ts/proto/model/chat";
import {
  EventDrawingDetect,
  EventDrawingUpdate,
} from "../../proto-ts/proto/model/drawing";
import { EventPointerUpdate } from "../../proto-ts/proto/model/pointer";
import { EventUserUpdate } from "../../proto-ts/proto/model/user";

export enum WebsocketInputEventsEnum {
  CHAT_MESSAGE = "chatmessage",
  POINTER_UPDATE = "pointerupdate",
  DRAWING_UPDATE = "drawingupdate",
  BOARD_UPDATE = "boardupdate",
  USER_UPDATE = "userupdate",
  DRAWING_DETECT = "drawingdetect",
}

export class ChatMessageInput extends CustomEvent<EventChatMessage> {
  constructor(init: CustomEventInit<EventChatMessage>) {
    super(WebsocketInputEventsEnum.CHAT_MESSAGE, init);
  }
}

export class PointerUpdateInput extends CustomEvent<EventPointerUpdate> {
  constructor(init: CustomEventInit<EventPointerUpdate>) {
    super(WebsocketInputEventsEnum.POINTER_UPDATE, init);
  }
}

export class DrawingUpdateInput extends CustomEvent<EventDrawingUpdate> {
  constructor(init: CustomEventInit<EventDrawingUpdate>) {
    super(WebsocketInputEventsEnum.DRAWING_UPDATE, init);
  }
}

export class BoardUpdateInput extends CustomEvent<EventBoardUpdate> {
  constructor(init: CustomEventInit<EventBoardUpdate>) {
    super(WebsocketInputEventsEnum.BOARD_UPDATE, init);
  }
}

export class UserUpdateInput extends CustomEvent<EventUserUpdate> {
  constructor(init: CustomEventInit<EventUserUpdate>) {
    super(WebsocketInputEventsEnum.USER_UPDATE, init);
  }
}

export class DrawingDetectInput extends CustomEvent<EventDrawingDetect> {
  constructor(init: CustomEventInit<EventDrawingDetect>) {
    super(WebsocketInputEventsEnum.DRAWING_DETECT, init);
  }
}
