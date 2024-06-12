export interface RpcError extends Error {
  code?: string;
  meta?: Record<string, string>;
  methodName: string;
  name: string;
  serviceName: string;
}

export const isRpcError = (
  error: RpcError | unknown
): error is RpcError => {
  const { name } = error as Error;
  return name === "RpcError";
};

export enum AuthEnum {
  // sessionStorage key for storing JWT token
  TOKEN_KEY = "u-token",
}
