// @generated by protobuf-ts 2.9.4 with parameter ts_nocheck
// @generated from protobuf file "proto/model/user.proto" (package "model", syntax proto3)
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
import { UUID } from "./common";
/**
 * @generated from protobuf message model.User
 */
export interface User {
    /**
     * @generated from protobuf field: model.UUID id = 1;
     */
    id?: UUID;
    /**
     * @generated from protobuf field: string name = 2;
     */
    name: string;
    /**
     * @generated from protobuf field: google.protobuf.Timestamp created_at = 3;
     */
    createdAt?: Timestamp;
}
/**
 * Create a new user
 *
 * @generated from protobuf message model.CreateUserRequest
 */
export interface CreateUserRequest {
    /**
     * @generated from protobuf field: string name = 1;
     */
    name: string;
}
/**
 * Create user response
 * Returns a JWT token used to validate the user
 *
 * @generated from protobuf message model.CreateUserResponse
 */
export interface CreateUserResponse {
    /**
     * @generated from protobuf field: string token = 1;
     */
    token: string;
    /**
     * @generated from protobuf field: model.User user = 2;
     */
    user?: User;
}
/**
 * Returns the current user
 *
 * @generated from protobuf message model.GetUserResponse
 */
export interface GetUserResponse {
    /**
     * @generated from protobuf field: model.User user = 1;
     */
    user?: User;
}
/**
 * @generated from protobuf message model.EventUserUpdate
 */
export interface EventUserUpdate {
    /**
     * @generated from protobuf field: model.EventMetadata metadata = 1;
     */
    metadata?: EventMetadata;
    /**
     * @generated from protobuf field: model.User user = 2;
     */
    user?: User;
}
// @generated message type with reflection information, may provide speed optimized methods
class User$Type extends MessageType<User> {
    constructor() {
        super("model.User", [
            { no: 1, name: "id", kind: "message", T: () => UUID },
            { no: 2, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 3, name: "created_at", kind: "message", T: () => Timestamp }
        ]);
    }
    create(value?: PartialMessage<User>): User {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<User>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: User): User {
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
                case /* google.protobuf.Timestamp created_at */ 3:
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
    internalBinaryWrite(message: User, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.UUID id = 1; */
        if (message.id)
            UUID.internalBinaryWrite(message.id, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* string name = 2; */
        if (message.name !== "")
            writer.tag(2, WireType.LengthDelimited).string(message.name);
        /* google.protobuf.Timestamp created_at = 3; */
        if (message.createdAt)
            Timestamp.internalBinaryWrite(message.createdAt, writer.tag(3, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.User
 */
export const User = new User$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateUserRequest$Type extends MessageType<CreateUserRequest> {
    constructor() {
        super("model.CreateUserRequest", [
            { no: 1, name: "name", kind: "scalar", T: 9 /*ScalarType.STRING*/ }
        ]);
    }
    create(value?: PartialMessage<CreateUserRequest>): CreateUserRequest {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.name = "";
        if (value !== undefined)
            reflectionMergePartial<CreateUserRequest>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateUserRequest): CreateUserRequest {
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
    internalBinaryWrite(message: CreateUserRequest, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
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
 * @generated MessageType for protobuf message model.CreateUserRequest
 */
export const CreateUserRequest = new CreateUserRequest$Type();
// @generated message type with reflection information, may provide speed optimized methods
class CreateUserResponse$Type extends MessageType<CreateUserResponse> {
    constructor() {
        super("model.CreateUserResponse", [
            { no: 1, name: "token", kind: "scalar", T: 9 /*ScalarType.STRING*/ },
            { no: 2, name: "user", kind: "message", T: () => User }
        ]);
    }
    create(value?: PartialMessage<CreateUserResponse>): CreateUserResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        message.token = "";
        if (value !== undefined)
            reflectionMergePartial<CreateUserResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: CreateUserResponse): CreateUserResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* string token */ 1:
                    message.token = reader.string();
                    break;
                case /* model.User user */ 2:
                    message.user = User.internalBinaryRead(reader, reader.uint32(), options, message.user);
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
    internalBinaryWrite(message: CreateUserResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* string token = 1; */
        if (message.token !== "")
            writer.tag(1, WireType.LengthDelimited).string(message.token);
        /* model.User user = 2; */
        if (message.user)
            User.internalBinaryWrite(message.user, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.CreateUserResponse
 */
export const CreateUserResponse = new CreateUserResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class GetUserResponse$Type extends MessageType<GetUserResponse> {
    constructor() {
        super("model.GetUserResponse", [
            { no: 1, name: "user", kind: "message", T: () => User }
        ]);
    }
    create(value?: PartialMessage<GetUserResponse>): GetUserResponse {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<GetUserResponse>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: GetUserResponse): GetUserResponse {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* model.User user */ 1:
                    message.user = User.internalBinaryRead(reader, reader.uint32(), options, message.user);
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
    internalBinaryWrite(message: GetUserResponse, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.User user = 1; */
        if (message.user)
            User.internalBinaryWrite(message.user, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.GetUserResponse
 */
export const GetUserResponse = new GetUserResponse$Type();
// @generated message type with reflection information, may provide speed optimized methods
class EventUserUpdate$Type extends MessageType<EventUserUpdate> {
    constructor() {
        super("model.EventUserUpdate", [
            { no: 1, name: "metadata", kind: "message", T: () => EventMetadata },
            { no: 2, name: "user", kind: "message", T: () => User }
        ]);
    }
    create(value?: PartialMessage<EventUserUpdate>): EventUserUpdate {
        const message = globalThis.Object.create((this.messagePrototype!));
        if (value !== undefined)
            reflectionMergePartial<EventUserUpdate>(this, message, value);
        return message;
    }
    internalBinaryRead(reader: IBinaryReader, length: number, options: BinaryReadOptions, target?: EventUserUpdate): EventUserUpdate {
        let message = target ?? this.create(), end = reader.pos + length;
        while (reader.pos < end) {
            let [fieldNo, wireType] = reader.tag();
            switch (fieldNo) {
                case /* model.EventMetadata metadata */ 1:
                    message.metadata = EventMetadata.internalBinaryRead(reader, reader.uint32(), options, message.metadata);
                    break;
                case /* model.User user */ 2:
                    message.user = User.internalBinaryRead(reader, reader.uint32(), options, message.user);
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
    internalBinaryWrite(message: EventUserUpdate, writer: IBinaryWriter, options: BinaryWriteOptions): IBinaryWriter {
        /* model.EventMetadata metadata = 1; */
        if (message.metadata)
            EventMetadata.internalBinaryWrite(message.metadata, writer.tag(1, WireType.LengthDelimited).fork(), options).join();
        /* model.User user = 2; */
        if (message.user)
            User.internalBinaryWrite(message.user, writer.tag(2, WireType.LengthDelimited).fork(), options).join();
        let u = options.writeUnknownFields;
        if (u !== false)
            (u == true ? UnknownFieldHandler.onWrite : u)(this.typeName, message, writer);
        return writer;
    }
}
/**
 * @generated MessageType for protobuf message model.EventUserUpdate
 */
export const EventUserUpdate = new EventUserUpdate$Type();
