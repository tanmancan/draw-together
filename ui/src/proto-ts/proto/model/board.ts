// @generated by protobuf-ts 2.9.4 with parameter ts_nocheck
// @generated from protobuf file "proto/model/board.proto" (package "model", syntax proto3)
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
import { EventMetadata } from "./common";
import { Timestamp } from "../../google/protobuf/timestamp";
import { User } from "./user";
import { UUID } from "./common";
/**
 * @generated from protobuf message model.Board
 */
export interface Board {
    /**
     * Board ID - UUID
     *
     * @generated from protobuf field: model.UUID id = 1;
     */
    id?: UUID;
    /**
     * Board name
     *
     * @generated from protobuf field: string name = 2;
     */
    name: string;
    /**
     * Board owner
     *
     * @generated from protobuf field: model.User owner = 3;
     */
    owner?: User;
    /**
     * Board users
     *
     * @generated from protobuf field: repeated model.User board_users = 4;
     */
    boardUsers: User[];
    /**
     * Time board was created in system
     *
     * @generated from protobuf field: google.protobuf.Timestamp created_at = 5;
     */
    createdAt?: Timestamp;
}
/**
 * @generated from protobuf message model.CreateBoardRequest
 */
export interface CreateBoardRequest {
    /**
     * @generated from protobuf field: string name = 1;
     */
    name: string;
}
/**
 * @generated from protobuf message model.GetBoardRequest
 */
export interface GetBoardRequest {
    /**
     * @generated from protobuf field: model.UUID id = 1;
     */
    id?: UUID;
}
/**
 * @generated from protobuf message model.GetBoardResponse
 */
export interface GetBoardResponse {
    /**
     * @generated from protobuf field: model.Board board = 1;
     */
    board?: Board;
}
/**
 * @generated from protobuf message model.EventBoardUpdate
 */
export interface EventBoardUpdate {
    /**
     * @generated from protobuf field: model.EventMetadata metadata = 1;
     */
    metadata?: EventMetadata;
    /**
     * @generated from protobuf field: model.Board board = 2;
     */
    board?: Board;
}
// @generated message type with reflection information, may provide speed optimized methods
class Board$Type extends MessageType<Board> {
    constructor() {
        super("model.Board", [
            { no: 1, name: "id", kind: "message", T: () => UUID },
            { no: 2, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "owner", kind: "message", T: () => User },
            { no: 4, name: "board_users", kind: "message", repeat: 1 /*RepeatType.PACKED*/, T: () => User },
            { no: 5, name: "created_at", kind: "message", T: () => Timestamp }
        ]);
    }
    create(value?: PartialMessage<Board>): Board {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.name = "";
        message.boardUsers = [];
        if (value !== undefined)
            reflectionMergePartial<Board>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: Board): Board {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* model.UUID id */ 1:
                    message.id = UUID.internalBinaryRead(reader, reader.uint32(), options, message.id);
                    break;
                case /* string name */ 2:
                    message.name = reader.string();
                    break;
                case /* model.User owner */ 3:
                    message.owner = User.internalBinaryRead(reader, reader.uint32(), options, message.owner);
                    break;
                case /* repeated model.User board_users */ 4:
                    message.boardUsers.push(User.internalBinaryRead(reader, reader.uint32(), options));
                    break;
                case /* google.protobuf.Timestamp created_at */ 5:
                    message.createdAt = Timestamp.internalBinaryRead(reader, reader.uint32(), options, message.createdAt);
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
    internalBinaryWrite(message: Board, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.UUID id = 1; */
        if (message.id)
            UUID.internalBinaryWrite(message.id, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string name = 2; */
        if (message.name !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.name);
        /* model.User owner = 3; */
        if (message.owner)
            User.internalBinaryWrite(message.owner, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        /* repeated model.User board_users = 4; */
        for (let i = 0; i < message.boardUsers.length; i++)
            User.internalBinaryWrite(message.boardUsers[i], writer.tag(4, WireType.LengthDelimited).fork(), options).join();
        /* google.protobuf.Timestamp created_at = 5; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(5, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.Board
 */
export const Board = new Board$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateBoardRequest$Type extends MessageType<CreateBoardRequest> {
    constructor() {
        super("model.CreateBoardRequest", [
            { no: 1, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<CreateBoardRequest>): CreateBoardRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<CreateBoardRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateBoardRequest): CreateBoardRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string name */ 1:
                    message.name = reader.string();
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
    internalBinaryWrite(message: CreateBoardRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string name = 1; */
        if (message.name !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.name);
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.CreateBoardRequest
 */
export const CreateBoardRequest = new CreateBoardRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetBoardRequest$Type extends MessageType<GetBoardRequest> {
    constructor() {
        super("model.GetBoardRequest", [
            { no: 1, name: "id", kind: "message", T: () => UUID }
        ]);
    }
    create(value?: PartialMessage<GetBoardRequest>): GetBoardRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<GetBoardRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetBoardRequest): GetBoardRequest {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* model.UUID id */ 1:
                    message.id = UUID.internalBinaryRead(reader, reader.uint32(), options, message.id);
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
    internalBinaryWrite(message: GetBoardRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.UUID id = 1; */
        if (message.id)
            UUID.internalBinaryWrite(message.id, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.GetBoardRequest
 */
export const GetBoardRequest = new GetBoardRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetBoardResponse$Type extends MessageType<GetBoardResponse> {
    constructor() {
        super("model.GetBoardResponse", [
            { no: 1, name: "board", kind: "message", T: () => Board }
        ]);
    }
    create(value?: PartialMessage<GetBoardResponse>): GetBoardResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<GetBoardResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetBoardResponse): GetBoardResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* model.Board board */ 1:
                    message.board = Board.internalBinaryRead(reader, reader.uint32(), options, message.board);
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
    internalBinaryWrite(message: GetBoardResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.Board board = 1; */
        if (message.board)
            Board.internalBinaryWrite(message.board, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.GetBoardResponse
 */
export const GetBoardResponse = new GetBoardResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EventBoardUpdate$Type extends MessageType<EventBoardUpdate> {
    constructor() {
        super("model.EventBoardUpdate", [
            { no: 1, name: "metadata", kind: "message", T: () => EventMetadata },
            { no: 2, name: "board", kind: "message", T: () => Board }
        ]);
    }
    create(value?: PartialMessage<EventBoardUpdate>): EventBoardUpdate {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<EventBoardUpdate>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: EventBoardUpdate): EventBoardUpdate {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* model.EventMetadata metadata */ 1:
                    message.metadata = EventMetadata.internalBinaryRead(reader, reader.uint32(), options, message.metadata);
                    break;
                case /* model.Board board */ 2:
                    message.board = Board.internalBinaryRead(reader, reader.uint32(), options, message.board);
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
    internalBinaryWrite(message: EventBoardUpdate, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.EventMetadata metadata = 1; */
        if (message.metadata)
            EventMetadata.internalBinaryWrite(message.metadata, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* model.Board board = 2; */
        if (message.board)
            Board.internalBinaryWrite(message.board, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.EventBoardUpdate
 */
export const EventBoardUpdate = new EventBoardUpdate$Type();
