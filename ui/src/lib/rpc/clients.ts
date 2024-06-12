import { AuthEnum } from "./types";
import { GrpcWebFetchTransport } from "@protobuf-ts/grpcweb-transport";
import {
  BoardServiceClient,
  ChatServiceClient,
  PointerServiceClient,
  UserServiceClient,
} from "../../proto-ts/proto/service.client";

export const getToken = (): string => {
  return window.sessionStorage.getItem(AuthEnum.TOKEN_KEY) ?? "";
};

let RPC_BASE_URL = "/";

if (import.meta.env.DEV) {
  RPC_BASE_URL = import.meta.env.VITE_RPC_BASE_URL ?? RPC_BASE_URL;
}

const X_HEADER_TOKEN = `X-Request-Token`;

// in milliseconds
const RPC_TIMEOUT = 30000;

class ProtoClients {
  constructor(
    private transport = new GrpcWebFetchTransport({
      baseUrl: RPC_BASE_URL,
      meta: {
        [X_HEADER_TOKEN]: getToken(),
      },
      timeout: RPC_TIMEOUT,
      format: "binary",
    }),
    private boardClient = new BoardServiceClient(transport),
    private chatClient = new ChatServiceClient(transport),
    private userClient = new UserServiceClient(transport),
    private pointerClient = new PointerServiceClient(transport)
  ) {}

  reset() {
    this.transport = new GrpcWebFetchTransport({
      baseUrl: RPC_BASE_URL,
      meta: {
        [X_HEADER_TOKEN]: getToken(),
      },
      timeout: RPC_TIMEOUT,
      format: "binary",
    });

    this.boardClient = new BoardServiceClient(this.transport);
    this.chatClient = new ChatServiceClient(this.transport);
    this.userClient = new UserServiceClient(this.transport);
    this.pointerClient = new PointerServiceClient(this.transport);
  }

  getBoardClient() {
    return this.boardClient;
  }

  getChatClient() {
    return this.chatClient;
  }

  getUserClient() {
    return this.userClient;
  }

  getPointerClient() {
    return this.pointerClient;
  }
}

export const protoClients = new ProtoClients();
