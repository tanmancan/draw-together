import { MessageType } from "@protobuf-ts/runtime";

export const protoFromBlob = async <T extends object>(
  data: Blob,
  model: MessageType<T>
): Promise<T> => {
  const buffer = await data.arrayBuffer();
  const msgBin = new Uint8Array(buffer);
  return model.fromBinary(msgBin);
};
