export const NOTIFICATION_EVENT = "notification";

export enum NotificationSeverityEnum {
  INFO = "blue",
  WARNING = "orange",
  ERROR = "red",
  SUCCESS = "green",
}

export interface INotificationPayload {
  severity: NotificationSeverityEnum;
  title: string;
  body: string;
  duration?: number;
}

export class NotificationEvent extends CustomEvent<INotificationPayload> {
  constructor(init: CustomEventInit<INotificationPayload>) {
    super(NOTIFICATION_EVENT, init);
  }
}

export const showNotification = (
  title: string,
  body: string,
  severity: NotificationSeverityEnum = NotificationSeverityEnum.INFO
) => {
  const detail: INotificationPayload = {
    title,
    body,
    severity,
  };
  const ev = new NotificationEvent({
    detail,
  });
  window.dispatchEvent(ev);
};
