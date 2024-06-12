import { protoClients } from "../rpc/clients";
import { Board } from "../../proto-ts/proto/model/board";
import {
  DrawingDetectResponse,
  GetBoardDrawingsResponse,
  ImageData,
} from "../../proto-ts/proto/model/drawing";
import { CircuitBreaker } from "../circuitbreaker";
import {
  NotificationSeverityEnum,
  showNotification,
} from "../../components/ToastNotification/events";
import { isRpcError } from "../rpc/types";

const cb = new CircuitBreaker("boardClient");

export const handleCreateBoard = async (
  name: string
): Promise<Board | null> => {
  const protectedCreateBoard = cb.protect(
    async () =>
      await protoClients.getBoardClient().createBoard({
        name,
      })
  );
  const res = await protectedCreateBoard();
  const body = res?.response;

  return body?.board?.id?.value ? body.board : null;
};

export const handleGetBoard = async (
  boardID: string
): Promise<Board | null> => {
  const protectedGetBoard = cb.protect(
    async () =>
      await protoClients.getBoardClient().getBoard({
        id: {
          value: boardID,
        },
      })
  );
  const res = await protectedGetBoard();
  const body = res?.response;
  return body?.board?.id?.value ? body.board : null;
};

export const handleDrawingUpdate = async (
  boardID: string,
  imageData: ImageData
) => {
  const protectedDrawingUpdate = cb.protect(
    async () =>
      await protoClients.getBoardClient().updateDrawing({
        boardId: {
          value: boardID,
        },
        imageData,
      })
  );
  const res = await protectedDrawingUpdate();
  return res?.response ?? false;
};

export const handleGetBoardDrawings = async (
  boardID: string
): Promise<GetBoardDrawingsResponse | null> => {
  try {
    const protectedDrawingUpdate = cb.protect(
      async () =>
        await protoClients.getBoardClient().getBoardDrawings({
          boardId: {
            value: boardID,
          },
        })
    );
    const res = await protectedDrawingUpdate();
    return res?.response ?? null;
  } catch (error) {
    console.error(error);
    let severity = NotificationSeverityEnum.ERROR,
      title = "",
      body = "";

    if (isRpcError(error)) {
      const { code, methodName, serviceName, message } = error;
      severity = NotificationSeverityEnum.ERROR;
      title = `${code}: ${methodName} (${serviceName})`;
      body = message;
    }

    showNotification(title, body, severity);
    return null;
  }
};

export const handleDrawingDetect = async (
  boardID: string
): Promise<DrawingDetectResponse | null> => {
  try {
    const protectedDrawingDetect = cb.protect(
      async () =>
        await protoClients.getBoardClient().drawingDetect({
          boardId: {
            value: boardID,
          },
        })
    );
    const res = await protectedDrawingDetect();
    return res?.response ?? null;
  } catch (error) {
    console.error(error);
    let severity = NotificationSeverityEnum.ERROR,
      title = "",
      body = "";

    if (isRpcError(error)) {
      const { code, methodName, serviceName, message } = error;
      severity = NotificationSeverityEnum.ERROR;
      title = `${code}: ${methodName} (${serviceName})`;
      body = message;
    }

    showNotification(title, body, severity);
    return null;
  }
};
