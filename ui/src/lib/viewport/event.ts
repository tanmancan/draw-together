import { PointerStack } from "./types";

export enum ViewPortEventsEnum {
  POINTER_UPDATE_EVENT = "pointerupdate",
  DRAWING_UPDATE_EVENT = "drawingupdate",
}

export interface IPointerUpdateEventPayload {
  boardID: string;
  userID: string;
  pointerStack: PointerStack;
}

export class PointerUpdateEvent extends CustomEvent<IPointerUpdateEventPayload> {
  constructor(init: CustomEventInit<IPointerUpdateEventPayload>) {
    super(ViewPortEventsEnum.POINTER_UPDATE_EVENT, init);
  }
}

export interface IDrawingUpdateEventPayload {
  boardID: string;
  userID: string;
  imageData: Uint8Array;
}

export class DrawingUpdateEvent extends CustomEvent<IDrawingUpdateEventPayload> {
  constructor(init: CustomEventInit<IDrawingUpdateEventPayload>) {
    super(ViewPortEventsEnum.DRAWING_UPDATE_EVENT, init);
  }
}
