// @generated by protobuf-ts 2.9.4 with parameter ts_nocheck
// @generated from protobuf file "proto/service.proto" (package "service", syntax proto3)
// tslint:disable
// @ts-nocheck
import { PointerService } from "./service";
import type { UpdatePointerResponse } from "./model/pointer";
import type { UpdatePointerRequest } from "./model/pointer";
import { ChatService } from "./service";
import type { GetBoardMessagesResponse } from "./model/chat";
import type { GetBoardMessagesRequest } from "./model/chat";
import type { ChatMessageResponse } from "./model/chat";
import type { ChatMessageRequest } from "./model/chat";
import { BoardService } from "./service";
import type { DrawingDetectResponse } from "./model/drawing";
import type { GetBoardDrawingsResponse } from "./model/drawing";
import type { GetBoardDrawingsRequest } from "./model/drawing";
import type { UpdateDrawingResponse } from "./model/drawing";
import type { UpdateDrawingRequest } from "./model/drawing";
import type { GetBoardRequest } from "./model/board";
import type { GetBoardResponse } from "./model/board";
import type { CreateBoardRequest } from "./model/board";
import type { RpcTransport } from "@protobuf-ts/runtime-rpc";
import type { ServiceInfo } from "@protobuf-ts/runtime-rpc";
import { UserService } from "./service";
import type { GetUserResponse } from "./model/user";
import type { Empty } from "../google/protobuf/empty";
import { stackIntercept } from "@protobuf-ts/runtime-rpc";
import type { CreateUserResponse } from "./model/user";
import type { CreateUserRequest } from "./model/user";
import type { UnaryCall } from "@protobuf-ts/runtime-rpc";
import type { RpcOptions } from "@protobuf-ts/runtime-rpc";
/**
 * @generated from protobuf service service.UserService
 */
export interface IUserServiceClient {
    /**
     * @generated from protobuf rpc: CreateUser(model.CreateUserRequest) returns (model.CreateUserResponse);
     */
    createUser(input: CreateUserRequest, options?: RpcOptions): UnaryCall<CreateUserRequest, CreateUserResponse>;
    /**
     * @generated from protobuf rpc: GetUser(google.protobuf.Empty) returns (model.GetUserResponse);
     */
    getUser(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetUserResponse>;
    /**
     * @generated from protobuf rpc: DeleteUser(google.protobuf.Empty) returns (model.GetUserResponse);
     */
    deleteUser(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetUserResponse>;
}
/**
 * @generated from protobuf service service.UserService
 */
export class UserServiceClient implements IUserServiceClient, ServiceInfo {
    typeName = UserService.typeName;
    methods = UserService.methods;
    options = UserService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: CreateUser(model.CreateUserRequest) returns (model.CreateUserResponse);
     */
    createUser(input: CreateUserRequest, options?: RpcOptions): UnaryCall<CreateUserRequest, CreateUserResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateUserRequest, CreateUserResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetUser(google.protobuf.Empty) returns (model.GetUserResponse);
     */
    getUser(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetUserResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, GetUserResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DeleteUser(google.protobuf.Empty) returns (model.GetUserResponse);
     */
    deleteUser(input: Empty, options?: RpcOptions): UnaryCall<Empty, GetUserResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<Empty, GetUserResponse>("unary", this._transport, method, opt, input);
    }
}
/**
 * @generated from protobuf service service.BoardService
 */
export interface IBoardServiceClient {
    /**
     * @generated from protobuf rpc: CreateBoard(model.CreateBoardRequest) returns (model.GetBoardResponse);
     */
    createBoard(input: CreateBoardRequest, options?: RpcOptions): UnaryCall<CreateBoardRequest, GetBoardResponse>;
    /**
     * @generated from protobuf rpc: GetBoard(model.GetBoardRequest) returns (model.GetBoardResponse);
     */
    getBoard(input: GetBoardRequest, options?: RpcOptions): UnaryCall<GetBoardRequest, GetBoardResponse>;
    /**
     * @generated from protobuf rpc: DeleteBoard(model.GetBoardRequest) returns (model.GetBoardResponse);
     */
    deleteBoard(input: GetBoardRequest, options?: RpcOptions): UnaryCall<GetBoardRequest, GetBoardResponse>;
    /**
     * @generated from protobuf rpc: UpdateDrawing(model.UpdateDrawingRequest) returns (model.UpdateDrawingResponse);
     */
    updateDrawing(input: UpdateDrawingRequest, options?: RpcOptions): UnaryCall<UpdateDrawingRequest, UpdateDrawingResponse>;
    /**
     * @generated from protobuf rpc: GetBoardDrawings(model.GetBoardDrawingsRequest) returns (model.GetBoardDrawingsResponse);
     */
    getBoardDrawings(input: GetBoardDrawingsRequest, options?: RpcOptions): UnaryCall<GetBoardDrawingsRequest, GetBoardDrawingsResponse>;
    /**
     * @generated from protobuf rpc: DrawingDetect(model.GetBoardDrawingsRequest) returns (model.DrawingDetectResponse);
     */
    drawingDetect(input: GetBoardDrawingsRequest, options?: RpcOptions): UnaryCall<GetBoardDrawingsRequest, DrawingDetectResponse>;
}
/**
 * @generated from protobuf service service.BoardService
 */
export class BoardServiceClient implements IBoardServiceClient, ServiceInfo {
    typeName = BoardService.typeName;
    methods = BoardService.methods;
    options = BoardService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: CreateBoard(model.CreateBoardRequest) returns (model.GetBoardResponse);
     */
    createBoard(input: CreateBoardRequest, options?: RpcOptions): UnaryCall<CreateBoardRequest, GetBoardResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<CreateBoardRequest, GetBoardResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetBoard(model.GetBoardRequest) returns (model.GetBoardResponse);
     */
    getBoard(input: GetBoardRequest, options?: RpcOptions): UnaryCall<GetBoardRequest, GetBoardResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetBoardRequest, GetBoardResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DeleteBoard(model.GetBoardRequest) returns (model.GetBoardResponse);
     */
    deleteBoard(input: GetBoardRequest, options?: RpcOptions): UnaryCall<GetBoardRequest, GetBoardResponse> {
        const method = this.methods[2], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetBoardRequest, GetBoardResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: UpdateDrawing(model.UpdateDrawingRequest) returns (model.UpdateDrawingResponse);
     */
    updateDrawing(input: UpdateDrawingRequest, options?: RpcOptions): UnaryCall<UpdateDrawingRequest, UpdateDrawingResponse> {
        const method = this.methods[3], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdateDrawingRequest, UpdateDrawingResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetBoardDrawings(model.GetBoardDrawingsRequest) returns (model.GetBoardDrawingsResponse);
     */
    getBoardDrawings(input: GetBoardDrawingsRequest, options?: RpcOptions): UnaryCall<GetBoardDrawingsRequest, GetBoardDrawingsResponse> {
        const method = this.methods[4], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetBoardDrawingsRequest, GetBoardDrawingsResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: DrawingDetect(model.GetBoardDrawingsRequest) returns (model.DrawingDetectResponse);
     */
    drawingDetect(input: GetBoardDrawingsRequest, options?: RpcOptions): UnaryCall<GetBoardDrawingsRequest, DrawingDetectResponse> {
        const method = this.methods[5], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetBoardDrawingsRequest, DrawingDetectResponse>("unary", this._transport, method, opt, input);
    }
}
/**
 * @generated from protobuf service service.ChatService
 */
export interface IChatServiceClient {
    /**
     * @generated from protobuf rpc: SendMessage(model.ChatMessageRequest) returns (model.ChatMessageResponse);
     */
    sendMessage(input: ChatMessageRequest, options?: RpcOptions): UnaryCall<ChatMessageRequest, ChatMessageResponse>;
    /**
     * @generated from protobuf rpc: GetBoardMessages(model.GetBoardMessagesRequest) returns (model.GetBoardMessagesResponse);
     */
    getBoardMessages(input: GetBoardMessagesRequest, options?: RpcOptions): UnaryCall<GetBoardMessagesRequest, GetBoardMessagesResponse>;
}
/**
 * @generated from protobuf service service.ChatService
 */
export class ChatServiceClient implements IChatServiceClient, ServiceInfo {
    typeName = ChatService.typeName;
    methods = ChatService.methods;
    options = ChatService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: SendMessage(model.ChatMessageRequest) returns (model.ChatMessageResponse);
     */
    sendMessage(input: ChatMessageRequest, options?: RpcOptions): UnaryCall<ChatMessageRequest, ChatMessageResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<ChatMessageRequest, ChatMessageResponse>("unary", this._transport, method, opt, input);
    }
    /**
     * @generated from protobuf rpc: GetBoardMessages(model.GetBoardMessagesRequest) returns (model.GetBoardMessagesResponse);
     */
    getBoardMessages(input: GetBoardMessagesRequest, options?: RpcOptions): UnaryCall<GetBoardMessagesRequest, GetBoardMessagesResponse> {
        const method = this.methods[1], opt = this._transport.mergeOptions(options);
        return stackIntercept<GetBoardMessagesRequest, GetBoardMessagesResponse>("unary", this._transport, method, opt, input);
    }
}
/**
 * @generated from protobuf service service.PointerService
 */
export interface IPointerServiceClient {
    /**
     * @generated from protobuf rpc: UpdatePointer(model.UpdatePointerRequest) returns (model.UpdatePointerResponse);
     */
    updatePointer(input: UpdatePointerRequest, options?: RpcOptions): UnaryCall<UpdatePointerRequest, UpdatePointerResponse>;
}
/**
 * @generated from protobuf service service.PointerService
 */
export class PointerServiceClient implements IPointerServiceClient, ServiceInfo {
    typeName = PointerService.typeName;
    methods = PointerService.methods;
    options = PointerService.options;
    constructor(private readonly _transport: RpcTransport) {
    }
    /**
     * @generated from protobuf rpc: UpdatePointer(model.UpdatePointerRequest) returns (model.UpdatePointerResponse);
     */
    updatePointer(input: UpdatePointerRequest, options?: RpcOptions): UnaryCall<UpdatePointerRequest, UpdatePointerResponse> {
        const method = this.methods[0], opt = this._transport.mergeOptions(options);
        return stackIntercept<UpdatePointerRequest, UpdatePointerResponse>("unary", this._transport, method, opt, input);
    }
}
