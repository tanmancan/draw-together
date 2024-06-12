import {
  NotificationSeverityEnum,
  showNotification,
} from "../../components/ToastNotification/events";
import { User } from "../../proto-ts/proto/model/user";
import { getToken, protoClients } from "../rpc/clients";
import { AuthEnum, isRpcError } from "../rpc/types";

export const handleCreateUser = async (name: string): Promise<User | null> => {
  const res = await protoClients.getUserClient().createUser({
    name,
  });
  const body = res.response;

  if (body?.token && body?.user?.id?.value) {
    window.sessionStorage.setItem(AuthEnum.TOKEN_KEY, body?.token);
    protoClients.reset();
    return body.user;
  }

  return null;
};

// Fetches user from API using a JWT token
export const handleGetUser = async (): Promise<User | null> => {
  try {
    const res = await protoClients.getUserClient().getUser({});
    const body = res.response;

    return body?.user?.id?.value ? body.user : null;
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

    if (getToken()) {
      showNotification(title, body, severity);
    }
    return null;
  }
};
