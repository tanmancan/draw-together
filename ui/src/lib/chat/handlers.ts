import { EventChatMessage } from "../../proto-ts/proto/model/chat";
import { CircuitBreaker } from "../circuitbreaker";
import { protoClients } from "../rpc/clients";

const cb = new CircuitBreaker("chatClient");

export const handleSendChatMessage = async (
  boardId: string,
  message: string[]
): Promise<boolean> => {
  const protectedSendMessage = cb.protect(
    async () =>
      await protoClients.getChatClient().sendMessage({
        boardId: {
          value: boardId,
        },
        message,
      })
  );
  const res = await protectedSendMessage();
  const body = res?.response;

  return body?.success ?? false;
};

export const handleGetBoardMessages = async (
  boardId: string
): Promise<EventChatMessage[]> => {
  const res = await protoClients.getChatClient().getBoardMessages({
    boardId: {
      value: boardId,
    },
  });
  const body = res.response;

  return body.messages ?? [];
};
