// @generated by protobuf-ts 2.9.4 with parameter ts_nocheck
// @generated from protobuf file "proto/model/events.proto" (package "model", syntax proto3)
// tslint:disable
// @ts-nocheck
import type { BinaryWriteOptions } from "@protobuf-ts/runtime";
import type { IBinaryWriter } from "@protobuf-ts/runtime";
import { WireType } from "@protobuf-ts/runtime";
import type { BinaryReadOptions } from "@protobuf-ts/runtime";
import type { IBinaryReader } from "@protobuf-ts/runtime";
import { UnknownFieldHandler } from "@protobuf-ts/runtime";
import type { PartialMessage } from "@protobuf-ts/runtime";
import { reflectionMergePartial } from "@protobuf-ts/runtime";
import { MessageType } from "@protobuf-ts/runtime";
import { EventDrawingDetectQueue } from "./drawing";
import { EventDrawingDetect } from "./drawing";
import { EventUserUpdate } from "./user";
import { EventDrawingUpdate } from "./drawing";
import { EventBoardUpdate } from "./board";
import { EventPointerUpdate } from "./pointer";
import { EventChatMessage } from "./chat";
/**
 * @generated from protobuf message model.WsEvent
 */
export interface WsEvent {
    /**
     * @generated from protobuf field: string req_id = 1;
     */
    reqId: string;
    /**
     * @generated from protobuf oneof: event
     */
    event: {
        oneofKind: "chatMessage";
        /**
         * @generated from protobuf field: model.EventChatMessage chat_message = 2;
         */
        chatMessage: EventChatMessage;
    } | {
        oneofKind: "pointerUpdate";
        /**
         * @generated from protobuf field: model.EventPointerUpdate pointer_update = 3;
         */
        pointerUpdate: EventPointerUpdate;
    } | {
        oneofKind: "boardUpdate";
        /**
         * @generated from protobuf field: model.EventBoardUpdate board_update = 4;
         */
        boardUpdate: EventBoardUpdate;
    } | {
        oneofKind: "drawingUpdate";
        /**
         * @generated from protobuf field: model.EventDrawingUpdate drawing_update = 5;
         */
        drawingUpdate: EventDrawingUpdate;
    } | {
        oneofKind: "userUpdate";
        /**
         * @generated from protobuf field: model.EventUserUpdate user_update = 6;
         */
        userUpdate: EventUserUpdate;
    } | {
        oneofKind: "drawingDetect";
        /**
         * @generated from protobuf field: model.EventDrawingDetect drawing_detect = 7;
         */
        drawingDetect: EventDrawingDetect;
    } | {
        oneofKind: "drawingDetectQueue";
        /**
         * @generated from protobuf field: model.EventDrawingDetectQueue drawing_detect_queue = 8;
         */
        drawingDetectQueue: EventDrawingDetectQueue;
    } | {
        oneofKind: undefined;
    };
}
// @generated message type with reflection information, may provide speed optimized methods
class WsEvent$Type extends MessageType<WsEvent> {
    constructor() {
        super("model.WsEvent", [
            { no: 1, name: "req_id", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "chat_message", kind: "message", oneof: "event", T: () => EventChatMessage },
            { no: 3, name: "pointer_update", kind: "message", oneof: "event", T: () => EventPointerUpdate },
            { no: 4, name: "board_update", kind: "message", oneof: "event", T: () => EventBoardUpdate },
            { no: 5, name: "drawing_update", kind: "message", oneof: "event", T: () => EventDrawingUpdate },
            { no: 6, name: "user_update", kind: "message", oneof: "event", T: () => EventUserUpdate },
            { no: 7, name: "drawing_detect", kind: "message", oneof: "event", T: () => EventDrawingDetect },
            { no: 8, name: "drawing_detect_queue", kind: "message", oneof: "event", T: () => EventDrawingDetectQueue }
        ]);
    }
    create(value?: PartialMessage<WsEvent>): WsEvent {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.reqId = "";
        message.event = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<WsEvent>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: WsEvent): WsEvent {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string req_id */ 1:
                    message.reqId = reader.string();
                    break;
                case /* model.EventChatMessage chat_message */ 2:
                    message.event = {
                        oneofKind: "chatMessage",
                        chatMessage: EventChatMessage.internalBinaryRead(reader, reader.uint32(), options, (message.event as any).chatMessage)
                    };
                    break;
                case /* model.EventPointerUpdate pointer_update */ 3:
                    message.event = {
                        oneofKind: "pointerUpdate",
                        pointerUpdate: EventPointerUpdate.internalBinaryRead(reader, reader.uint32(), options, (message.event as any).pointerUpdate)
                    };
                    break;
                case /* model.EventBoardUpdate board_update */ 4:
                    message.event = {
                        oneofKind: "boardUpdate",
                        boardUpdate: EventBoardUpdate.internalBinaryRead(reader, reader.uint32(), options, (message.event as any).boardUpdate)
                    };
                    break;
                case /* model.EventDrawingUpdate drawing_update */ 5:
                    message.event = {
                        oneofKind: "drawingUpdate",
                        drawingUpdate: EventDrawingUpdate.internalBinaryRead(reader, reader.uint32(), options, (message.event as any).drawingUpdate)
                    };
                    break;
                case /* model.EventUserUpdate user_update */ 6:
                    message.event = {
                        oneofKind: "userUpdate",
                        userUpdate: EventUserUpdate.internalBinaryRead(reader, reader.uint32(), options, (message.event as any).userUpdate)
                    };
                    break;
                case /* model.EventDrawingDetect drawing_detect */ 7:
                    message.event = {
                        oneofKind: "drawingDetect",
                        drawingDetect: EventDrawingDetect.internalBinaryRead(reader, reader.uint32(), options, (message.event as any).drawingDetect)
                    };
                    break;
                case /* model.EventDrawingDetectQueue drawing_detect_queue */ 8:
                    message.event = {
                        oneofKind: "drawingDetectQueue",
                        drawingDetectQueue: EventDrawingDetectQueue.internalBinaryRead(reader, reader.uint32(), options, (message.event as any).drawingDetectQueue)
                    };
                    break;
                default:
                    let u = options.readUnknownField;
                    if (u === "throw")
                        throw new globalThis.Error(`Unknown field ${fieldNo} (wire type ${wireType}) for ${this.typeName}`);
                    let d = reader.skip(wireType);
                    if (u !== false)
                        (u === true ? UnknownFieldHandler.onRead : u)(this.typeName, message, fieldNo, wireType, d);
            }
        }
        return message;
    }
    internalBinaryWrite(message: WsEvent, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string req_id = 1; */
        if (message.reqId !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.reqId);
        /* model.EventChatMessage chat_message = 2; */
        if (message.event.oneofKind === "chatMessage")
            EventChatMessage.internalBinaryWrite(message.event.chatMessage, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* model.EventPointerUpdate pointer_update = 3; */
        if (message.event.oneofKind === "pointerUpdate")
            EventPointerUpdate.internalBinaryWrite(message.event.pointerUpdate, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* model.EventBoardUpdate board_update = 4; */
        if (message.event.oneofKind === "boardUpdate")
            EventBoardUpdate.internalBinaryWrite(message.event.boardUpdate, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* model.EventDrawingUpdate drawing_update = 5; */
        if (message.event.oneofKind === "drawingUpdate")
            EventDrawingUpdate.internalBinaryWrite(message.event.drawingUpdate, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        /* model.EventUserUpdate user_update = 6; */
        if (message.event.oneofKind === "userUpdate")
            EventUserUpdate.internalBinaryWrite(message.event.userUpdate, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* model.EventDrawingDetect drawing_detect = 7; */
        if (message.event.oneofKind === "drawingDetect")
            EventDrawingDetect.internalBinaryWrite(message.event.drawingDetect, writer.tag(7, WireType.LengthDelimited).fork(), options).join();
        /* model.EventDrawingDetectQueue drawing_detect_queue = 8; */
        if (message.event.oneofKind === "drawingDetectQueue")
            EventDrawingDetectQueue.internalBinaryWrite(message.event.drawingDetectQueue, writer.tag(8, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.WsEvent
 */
export const WsEvent = new WsEvent$Type();
