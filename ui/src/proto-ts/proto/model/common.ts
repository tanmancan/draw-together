// @generated by protobuf-ts 2.9.4 with parameter ts_nocheck
// @generated from protobuf file "proto/model/common.proto" (package "model", syntax proto3)
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
import { Timestamp } from "../../google/protobuf/timestamp";
/**
 * @generated from protobuf message model.UUID
 */
export interface UUID {
    /**
     * A string representation of UUID
     * value = "08f0097e-9d04-49bf-a92d-1ed8e1d95329"
     *
     * @generated from protobuf field: string value = 1;
     */
    value: string;
}
/**
 * @generated from protobuf message model.EventMetadata
 */
export interface EventMetadata {
    /**
     * Event ID
     *
     * @generated from protobuf field: model.UUID id = 1;
     */
    id?: UUID;
    /**
     * Board ID
     *
     * @generated from protobuf field: model.UUID board_id = 2;
     */
    boardId?: UUID;
    /**
     * Sender User ID
     *
     * @generated from protobuf field: model.UUID sender_id = 3;
     */
    senderId?: UUID;
    /**
     * When message was created in system
     *
     * @generated from protobuf field: google.protobuf.Timestamp created_at = 4;
     */
    createdAt?: Timestamp;
    /**
     * @generated from protobuf oneof: target
     */
    target: {
        oneofKind: "allUsers";
        /**
         * Sent to all users in the board
         *
         * @generated from protobuf field: bool all_users = 5;
         */
        allUsers: boolean;
    } | {
        oneofKind: "userId";
        /**
         * Sent only to specified user
         *
         * @generated from protobuf field: model.UUID user_id = 6;
         */
        userId: UUID;
    } | {
        oneofKind: "allButSender";
        /**
         * Sent to all but the sender
         *
         * @generated from protobuf field: bool all_but_sender = 7;
         */
        allButSender: boolean;
    } | {
        oneofKind: undefined;
    };
}
// @generated message type with reflection information, may provide speed optimized methods
class UUID$Type extends MessageType<UUID> {
    constructor() {
        super("model.UUID", [
            { no: 1, name: "value", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<UUID>): UUID {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.value = "";
        if (value !== undefined)
            reflectionMergePartial<UUID>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: UUID): UUID {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string value */ 1:
                    message.value = reader.string();
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
    internalBinaryWrite(message: UUID, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string value = 1; */
        if (message.value !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.value);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.UUID
 */
export const UUID = new UUID$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EventMetadata$Type extends MessageType<EventMetadata> {
    constructor() {
        super("model.EventMetadata", [
            { no: 1, name: "id", kind: "message", T: () => UUID },
            { no: 2, name: "board_id", kind: "message", T: () => UUID },
            { no: 3, name: "sender_id", kind: "message", T: () => UUID },
            { no: 4, name: "created_at", kind: "message", T: () => Timestamp },
            { no: 5, name: "all_users", kind: "scalar", oneof: "target", T: 8 /*ScalarType.BOOL*/ },
            { no: 6, name: "user_id", kind: "message", oneof: "target", T: () => UUID },
            { no: 7, name: "all_but_sender", kind: "scalar", oneof: "target", T: 8 /*ScalarType.BOOL*/ }
        ]);
    }
    create(value?: PartialMessage<EventMetadata>): EventMetadata {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.target = { oneofKind: undefined };
        if (value !== undefined)
            reflectionMergePartial<EventMetadata>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: EventMetadata): EventMetadata {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* model.UUID id */ 1:
                    message.id = UUID.internalBinaryRead(reader, reader.uint32(), options, message.id);
                    break;
                case /* model.UUID board_id */ 2:
                    message.boardId = UUID.internalBinaryRead(reader, reader.uint32(), options, message.boardId);
                    break;
                case /* model.UUID sender_id */ 3:
                    message.senderId = UUID.internalBinaryRead(reader, reader.uint32(), options, message.senderId);
                    break;
                case /* google.protobuf.Timestamp created_at */ 4:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
                    break;
                case /* bool all_users */ 5:
                    message.target = {
                        oneofKind: "allUsers",
                        allUsers: reader.bool()
                    };
                    break;
                case /* model.UUID user_id */ 6:
                    message.target = {
                        oneofKind: "userId",
                        userId: UUID.internalBinaryRead(reader, reader.uint32(), options, (message.target as any).userId)
                    };
                    break;
                case /* bool all_but_sender */ 7:
                    message.target = {
                        oneofKind: "allButSender",
                        allButSender: reader.bool()
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
    internalBinaryWrite(message: EventMetadata, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.UUID id = 1; */
        if (message.id)
            UUID.internalBinaryWrite(message.id, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* model.UUID board_id = 2; */
        if (message.boardId)
            UUID.internalBinaryWrite(message.boardId, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        /* model.UUID sender_id = 3; */
        if (message.senderId)
            UUID.internalBinaryWrite(message.senderId, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* google.protobuf.Timestamp created_at = 4; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* bool all_users = 5; */
        if (message.target.oneofKind === "allUsers")
            writer.tag(5, WireType.Varint).bool(message.target.allUsers);
        /* model.UUID user_id = 6; */
        if (message.target.oneofKind === "userId")
            UUID.internalBinaryWrite(message.target.userId, writer.tag(6, WireType.LengthDelimited).fork(), options).join();
        /* bool all_but_sender = 7; */
        if (message.target.oneofKind === "allButSender")
            writer.tag(7, WireType.Varint).bool(message.target.allButSender);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.EventMetadata
 */
export const EventMetadata = new EventMetadata$Type();
