import {
  NotificationSeverityEnum,
  showNotification,
} from "../../components/ToastNotification/events";
import { BASE_HOST, BASE_PORT, BASE_PROTOCOL } from "../common";
import { getToken } from "../rpc/clients";
import { onMessageHandler } from "./handlers";

const PREFIX = "ws";
const PROTOCOL = BASE_PROTOCOL === "https:" ? `wss:` : `ws:`;
let PORT = BASE_PORT;
let HOST = `${PROTOCOL}//${BASE_HOST}${PORT ? `:${PORT}` : ""}`;

if (import.meta.env.DEV) {
  PORT = import.meta.env.VITE_WS_PORT ?? PORT;
  HOST = `${import.meta.env.VITE_WS_HOST}:${PORT}` ?? HOST;
}

function noOpHandler() {}

class WsBaseClient {
  protected ws?: WebSocket;
  protected params: Record<string, string> = {};
  protected readyState?: number;

  protected retryInterval: number = 3000;
  protected maxRetryCount: number = 5;
  protected currentRetryCount: number = 0;

  protected debugLog: boolean = false;

  protected onOpenHandler?: (e: Event) => void;
  protected onMessageHandler?: (e: MessageEvent) => void;
  protected onCloseHandler?: (e: CloseEvent) => void;
  protected onErrorHandler?: (e: Event) => void;

  constructor() {
    this.openHandler = this.openHandler.bind(this);
    this.messageHandler = this.messageHandler.bind(this);
    this.closeHandler = this.closeHandler.bind(this);
    this.errorHandler = this.errorHandler.bind(this);
    this.init = this.init.bind(this);
    this.connect = this.connect.bind(this);
  }

  init(
    params: Record<string, string> = {},
    onOpenHandler?: (e: Event) => void,
    onMessageHandler?: (e: MessageEvent) => void,
    onCloseHandler?: (e: CloseEvent) => void,
    onErrorHandler?: (e: Event) => void,
    debugLog: boolean = false
  ) {
    if (this.ws) {
      return;
    }

    this.params = params;
    this.debugLog = debugLog;

    this.onOpenHandler = onOpenHandler;
    this.onMessageHandler = onMessageHandler;
    this.onCloseHandler = onCloseHandler;
    this.onErrorHandler = onErrorHandler;

    this.connect();
  }

  protected connect(retry: boolean = false) {
    if (retry) {
      this.currentRetryCount++;
      this.teardown();

      if (this.currentRetryCount > this.maxRetryCount) {
        showNotification(
          "connection to server lost",
          "max retry attempted. giving up",
          NotificationSeverityEnum.ERROR
        );
        console.error(`max retry attempted giving up`);
        return;
      }

      showNotification(
        "connection to server lost",
        `retrying connection failed. attempt: ${this.currentRetryCount} of ${this.maxRetryCount}`,
        NotificationSeverityEnum.WARNING
      );
      console.warn(
        `retrying connection failed. attempt: ${this.currentRetryCount} of ${this.maxRetryCount}`
      );
    }

    this.ws = new WebSocket(this.getUrl());
    this.ws.addEventListener("open", this.openHandler);
    this.ws.addEventListener("message", this.messageHandler);
    this.ws.addEventListener("close", this.closeHandler);
    this.ws.addEventListener("error", this.errorHandler);
  }

  teardown() {
    const ws = this.getWs();

    if (!ws) {
      return;
    }

    if (this.getReadyState() === WebSocket.OPEN) {
      ws?.close(1000);
    }

    ws.removeEventListener("open", this.openHandler);
    ws.removeEventListener("message", this.messageHandler);
    ws.removeEventListener("close", this.closeHandler);
    ws.removeEventListener("error", this.errorHandler);
  }

  enableDebugLog() {
    this.debugLog = true;
  }

  disableDebugLog() {
    this.debugLog = false;
  }

  getWs(): WebSocket {
    if (!this.ws) {
      throw Error("init() must be called first.");
    }
    return this.ws;
  }

  protected getReadyState(): number | undefined {
    return this?.ws?.readyState;
  }

  protected getUrl(): string {
    const url = new URL(`${HOST}/${PREFIX}`);

    Object.entries(this.params).forEach(([key, val]) => {
      url.searchParams.set(key, val);
      url.searchParams.set(key, val);
    });

    url.searchParams.set("token", getToken() ?? "");

    return url.toString();
  }

  protected openHandler(e: Event) {
    if (this.debugLog) {
      console.log("onOpenEvent: ", e);
    }
    if (this.currentRetryCount > 0) {
      showNotification(
        "re-connected to server",
        "successfully re-connected to server",
        NotificationSeverityEnum.SUCCESS
      );
    }
    this.currentRetryCount = 0;
    this.onOpenHandler ? this.onOpenHandler(e) : noOpHandler();
  }

  protected messageHandler(e: MessageEvent) {
    if (this.debugLog) {
      console.log("onMessageEvent: ", e);
    }
    onMessageHandler(e);
    this.onMessageHandler ? this.onMessageHandler(e) : noOpHandler();
  }

  protected closeHandler(e: CloseEvent) {
    if (this.debugLog) {
      console.log("onCloseEvent: ", e);
    }

    if (!e.wasClean) {
      if (this.currentRetryCount === 0) {
        showNotification(
          "connection to server lost",
          "will retry connection",
          NotificationSeverityEnum.ERROR
        );
      }
      window.setTimeout(this.retryConnection.bind(this), this.retryInterval);
    }
    this.onCloseHandler ? this.onCloseHandler(e) : noOpHandler();
  }

  protected errorHandler(e: Event) {
    if (this.debugLog) {
      console.log("onErrorEvent: ", e);
    }
    this.onErrorHandler ? this.onErrorHandler(e) : noOpHandler();
  }

  protected retryConnection() {
    this.connect(true);
  }
}

export default WsBaseClient;
