import { CircuitBreaker } from "../circuitbreaker";
import { protoClients } from "../rpc/clients";
import { PointerStack } from "../viewport/types";

const cb = new CircuitBreaker("pointerClient");

export const handlePointerUpdate = async (
  pointerStack: PointerStack,
  boardId: string
): Promise<boolean> => {
  const protectedPointerUpdate = cb.protect(
    async () =>
      await protoClients.getPointerClient().updatePointer({
        boardId: {
          value: boardId,
        },
        pointerPositions: pointerStack,
      })
  );
  const res = await protectedPointerUpdate();
  const body = res?.response;

  return body?.success ?? false;
};
